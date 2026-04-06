package packages

import "fmt"

type UnknownPackageError struct {
	Name string
}

func (e UnknownPackageError) Error() string {
	return fmt.Sprintf("unknown package %q", e.Name)
}
