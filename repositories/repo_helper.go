package repositories

import (
	"fmt"
)

const (
	configurations = "configs/%s/%s/"
	allConfigs     = "configs"
)

var (
	groups    = "groups/%s/%s/%v/%s"
	allGroups = "groups"
)

func ConstructConfigKey(name string, version string) string {
	return fmt.Sprintf(configurations, name, version)
}

func ConstructConfigGroupKey(name string, version string, labels string, configName string) string {
	if labels == "" {
		return fmt.Sprintf("groups/%s/%s/%s", name, version, configName)
	}
	return fmt.Sprintf(groups, name, version, labels, configName)
}
