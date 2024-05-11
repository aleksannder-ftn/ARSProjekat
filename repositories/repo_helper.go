package repositories

import (
	"fmt"
)

const (
	configurations = "configs/%s/%s"
	all            = "configurations"
)

func constructKey(id string) string {
	return fmt.Sprintf(configurations, id)
}
