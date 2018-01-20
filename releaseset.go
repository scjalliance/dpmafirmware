package dpmafirmware

const startingBranchCap = 32

// ReleaseSet is a set of releases.
type ReleaseSet []Release

// Branches returns the release set with all releases grouped by branch.
func (rs ReleaseSet) Branches() (branches []Branch) {
	seen := make(map[string]int) // Maps branch names to indices within branches
	for i := range rs {
		branch := rs[i].Branch()
		if b, exists := seen[branch]; exists {
			branches[b].Releases = append(branches[b].Releases, rs[i])
		} else {
			b := len(branches)
			seen[branch] = b
			branches = append(branches, Branch{Name: branch})
			branches[b].Releases = make(ReleaseSet, 0, startingBranchCap)
			branches[b].Releases = append(branches[b].Releases, rs[i])
		}
	}
	return branches
}

// Branch returns the group of releases for the given branch.
func (rs ReleaseSet) Branch(branch string) (group Branch) {
	group.Name = branch
	group.Releases = make(ReleaseSet, 0, startingBranchCap)
	for i := range rs {
		if rs[i].Branch() == branch {
			group.Releases = append(group.Releases, rs[i])
		}
	}
	return
}

// Latest returns the most recent release within the set.
func (rs ReleaseSet) Latest() (r Release) {
	if len(rs) == 0 {
		return
	}
	r = rs[0]
	for i := 0; i < len(rs); i++ {
		if rs[i].Version > r.Version {
			r = rs[i]
		}
	}
	return
}

func (rs ReleaseSet) Len() int           { return len(rs) }
func (rs ReleaseSet) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ReleaseSet) Less(i, j int) bool { return rs[i].Version > rs[j].Version }
