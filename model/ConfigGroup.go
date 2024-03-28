package model

// TODO implement version as struct
type ConfigurationGroup struct {
	Name           string          `json:"name"`
	Id             int32           `json:"id"`
	Version        string          `json:"version"`
	Configurations []Configuration `json:"configurations"`
}

// TODO add methods for struct

type ConfigurationGroupRepository interface {
	// TODO add crud for config group repository
}
