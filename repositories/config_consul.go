package repositories

import "ars_projekat/model"

type ConfigurationConsulRepository struct {
	// TODO add connection
}

func (c ConfigurationConsulRepository) Get() {
	//TODO implement me
	panic("implement me")
}
func NewConfigConsulRepository() model.ConfigurationRepository {
	return ConfigurationConsulRepository{}
}
