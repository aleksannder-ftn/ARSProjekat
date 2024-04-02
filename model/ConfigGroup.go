package model

// ConfigurationGroup TODO implement version as struct
type ConfigurationGroup struct {
	Name           string          `json:"name"`
	Id             int32           `json:"id"`
	Version        Version         `json:"version"`
	Configurations []Configuration `json:"configurations"`
}

func (cg *ConfigurationGroup) SetName(name string) {
	cg.Name = name
}

func (cg *ConfigurationGroup) SetId(id int32) {
	cg.Id = id
}

func (cg *ConfigurationGroup) SetVersion(version Version) {
	cg.Version = version
}

func (cg *ConfigurationGroup) SetConfigurations(configs []Configuration) {
	cg.Configurations = configs
}

// TODO add methods for struct

/*
	Ne znam sta nam jos fali za config group, s obzirom da odvajamo na dva repoa predpostavljam da

ce config group imati neke dodatne nacine za pretragu. To cemo prodiskutovati na discordu
*/
type ConfigurationGroupRepository interface {
	Delete()
	Update()
	Create()
	FindById()
	// TODO add crud for config group repository
}
