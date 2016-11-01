package actions

import (
	"fmt"
)

type errComponentVersionNotFound struct {
	componentName string
}

func (e errComponentVersionNotFound) Error() string {
	return fmt.Sprintf("A ComponentVersion for component named '%s' not found.", e.componentName)
}
