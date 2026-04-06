package packages

import "fmt"

type Registry interface {
	List() []Definition
	Lookup(name string) (Definition, error)
}

type StaticRegistry struct {
	ordered []Definition
	byName  map[string]Definition
}

func NewStaticRegistry(definitions []Definition) (*StaticRegistry, error) {
	registry := &StaticRegistry{
		ordered: make([]Definition, 0, len(definitions)),
		byName:  make(map[string]Definition, len(definitions)),
	}

	for _, definition := range definitions {
		if err := definition.Validate(); err != nil {
			return nil, err
		}
		if _, exists := registry.byName[definition.Name]; exists {
			return nil, fmt.Errorf("duplicate package %q", definition.Name)
		}
		registry.ordered = append(registry.ordered, definition)
		registry.byName[definition.Name] = definition
	}

	return registry, nil
}

func (r *StaticRegistry) List() []Definition {
	definitions := make([]Definition, len(r.ordered))
	copy(definitions, r.ordered)
	return definitions
}

func (r *StaticRegistry) Lookup(name string) (Definition, error) {
	definition, ok := r.byName[name]
	if !ok {
		return Definition{}, UnknownPackageError{Name: name}
	}
	return definition, nil
}
