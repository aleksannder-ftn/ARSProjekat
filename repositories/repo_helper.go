package repositories

import (
	"fmt"
)

const (
	configurations = "configs/%s/%s/"
	all            = "configs"
	groups         = "configs/groups/%s/%s"
)

func ConstructConfigKey(name string, version string) string {
	return fmt.Sprintf(configurations, name, version)
}

func ConstructConfigGroupKey(name string, version string, labels []string) string {
	key := fmt.Sprintf(groups, name, version)
	for i, v := range labels {
		if i == len(labels) {
			key += v
		} else {
			key += v + ";"
		}
	}

	return fmt.Sprintf(key)
}
