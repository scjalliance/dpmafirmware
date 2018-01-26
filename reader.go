package dpmafirmware

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"hash"
	"io"
	"path/filepath"
	"regexp"
	"sync"
)

var errReaderClosed = errors.New("the dpma firmware reader was cancelled or closed")

// ModelMatcher is the regular expression used to match models within filenames.
// Matching subgroups will contain the model number.
const ModelMatcher = `(?:_|\b)([dD][0-9]{2,})(?:_|\b)`

// Reader reads firmware file data stored in a dpma firmware package.
type Reader struct {
	matcher *regexp.Regexp     // Matches models within file names
	cancel  context.CancelFunc // Request cancellation function for HTTP

	mutex   sync.RWMutex
	stream  io.Closer // Used by Release.Get to close the underlying response body
	gzip    *gzip.Reader
	tar     *tar.Reader
	md5     hash.Hash // MD5 hasher
	header0 *Header   // First regular file entry in the tar archive
	started bool      // Has the first entry been consumed?
	closed  bool      // Has the reader been closed?
}

// NewReader returns a firmware package reader for the data stream r.
// The data must be in tar.gz format.
//
// It is the caller's responsibility to close the reader when finished with it.
//
// It is safe to close the reader asynchronously; doing so allows the reader to
// be cancelled.
//
// It is possible for some portion of r to be consumed even when an error is
// returned by NewReader.
func NewReader(r io.Reader) (reader *Reader, err error) {
	matcher, err := regexp.Compile(ModelMatcher)
	if err != nil {
		return nil, fmt.Errorf("bad filename model matching expression: %v", err)
	}

	// Tee off a copy of everything we read to an MD5 hash
	md5Hash := md5.New()
	teeReader := io.TeeReader(r, md5Hash)

	// Setup the decompressor
	gzipReader, err := gzip.NewReader(teeReader)
	if err != nil {
		return nil, err
	}

	// Setup the extractor
	tarReader := tar.NewReader(gzipReader)

	// Prepare the reader
	reader = &Reader{
		gzip:    gzipReader,
		tar:     tarReader,
		md5:     md5Hash,
		matcher: matcher,
	}

	// Read the first entry to make sure everything works
	reader.header0, err = reader.next()
	if err != nil {
		gzipReader.Close()
		if err == io.EOF {
			err = errors.New("empty tar archive in dpma package")
		}
		return nil, err
	}

	return
}

// Next advances to the next file entry in the package. It will return io.EOF
// when there is no more data to be read.
func (r *Reader) Next() (header *Header, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return nil, errReaderClosed
	}
	if !r.started {
		r.started = true
		return r.header0, nil
	}

	// The next() function manages its own locking and unlocking
	return r.next()
}

// next advances to the next file entry in the archive. It skips non-file
// entries.
func (r *Reader) next() (header *Header, err error) {
	for {
		header, err := r.tar.Next()
		if err != nil {
			return nil, err
		}

		var models ModelSet
		for _, match := range r.matcher.FindAllStringSubmatch(header.Name, -1) {
			models = append(models, match[1])
		}

		switch header.Typeflag {
		case tar.TypeReg:
			return &Header{
				Name:    filepath.Base(header.Name),
				Path:    filepath.Clean(header.Name),
				Size:    header.Size,
				ModTime: header.ModTime,
				Models:  models,
			}, nil
		}
	}
}

// Read reads data for the current file from the package. It returns io.EOF
// when there is no more data to be read.
func (r *Reader) Read(b []byte) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return 0, errReaderClosed
	}

	return r.tar.Read(b)
}

// MD5Sum returns the MD5 hash of the data read so far. It can be called safely
// before or after r has been closed.
func (r *Reader) MD5Sum() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return fmt.Sprintf("%x", r.md5.Sum(nil))
}

// Close releases any resources consumed by the reader.
func (r *Reader) Close() error {
	// Here we intentionally perform cancellation outside of a lock.
	//
	// If the cancellation function can close the underlying stream, this will
	// allow Read() and Next() calls to stop immediately even when they're
	// holding a lock and waiting for stream data.
	//
	// It is safe to repeatedly call cancellation functions.
	if r.cancel != nil {
		r.cancel()
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return nil
	}

	r.closed = true

	var err1, err2 error
	err1 = r.gzip.Close()
	if r.stream != nil {
		err2 = r.stream.Close()
	}
	if err1 != nil {
		return err1
	}
	return err2
}
