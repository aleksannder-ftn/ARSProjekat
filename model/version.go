package model

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

func (v *Version) SetMajor(major int) {
	v.Major = major
}

func (v *Version) SetMinor(minor int) {
	v.Minor = minor
}

func (v *Version) SetPatch(patch int) {
	v.Patch = patch
}

func IncrementVersion(current Version) Version {
	newVersion := current

	newVersion.Patch++

	if newVersion.Patch > 9 {
		newVersion.Patch = 0
		newVersion.Minor++

		if newVersion.Minor > 9 {
			newVersion.Minor = 0
			newVersion.Major++
		}
	}

	return newVersion
}

type VersionRepository interface {
	Delete()
	Update()
	Create()
	Find()
}
