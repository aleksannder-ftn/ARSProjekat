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

type VersionRepository interface {
	Delete()
	Update()
	Create()
	Find()
}
