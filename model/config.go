package model

// TODO add version as struct, add labels as field (used for filtering)
type Configuration struct {
	Name       string            `json:"name"`
	Id         int32             `json:"id"`
	Version    string            `json:"version"`
	Parameters map[string]string `json:"parameters"`
}

// TODO add struct methods

type ConfigurationRepository interface {
	//TODO add CRUD methods
}
