package mapping

// MappingAction is the building block of the mapping process. Represents
// an individual action that updates the target resources
type MappingAction interface {
	// Apply performs the action to the result set using source
	Apply(source interface{}, resultSet ResultSet) error
}

// ResultSet manages the results being mapped
type ResultSet interface {
	AddResult(key string, result interface{})
	GetResult(key string) (interface{}, bool)
	GetAllResults() map[string]interface{}
}

// MappingActionResolver discovers MappingActions to be applied for a given
// source object
type MappingActionResolver interface {
	ResolveMappingActions(source interface{}) ([]MappingAction, error)
}
