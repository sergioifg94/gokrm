package introspect

import (
	"reflect"

	"github.com/sergioifg94/gokrm/pkg/mapping"
)

// IdentificationActionResolver is an ActionResolver that identifies new objects
// to add to the ResultSet, it resolves actions that instantiate the new result
// object
type IdentificationActionResolver struct {
	// ResultKey is a function that, given a StructTag, returns the key
	// that identifies the result object to instantiate
	ResultKey ResultKeyFunc

	// Constructor is a function that, given a result key, instantiates
	// a new resulting object
	Constructor func(resultKey string) interface{}
}

type ResultKeyFunc func(tag reflect.StructTag) (string, bool)

var _ mapping.MappingActionResolver = &IdentificationActionResolver{}

func (r *IdentificationActionResolver) ResolveMappingActions(source interface{}) ([]mapping.MappingAction, error) {
	value := reflect.ValueOf(source)
	sourceType := value.Type()

	identifiedResults := map[string]bool{}
	r.identifyResultsIn(identifiedResults, sourceType)

	return r.createActions(identifiedResults), nil
}

func (r *IdentificationActionResolver) identifyResultsIn(identifiedResults map[string]bool, sourceType reflect.Type) {
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)

		resultKey, ok := r.ResultKey(field.Tag)
		if !ok {
			if field.Type.Kind() == reflect.Struct {
				r.identifyResultsIn(identifiedResults, field.Type)
			}
			continue
		}

		identifiedResults[resultKey] = true
	}
}

func (r *IdentificationActionResolver) createActions(identifiedResults map[string]bool) []mapping.MappingAction {
	actions := []mapping.MappingAction{}

	for resultKey := range identifiedResults {
		k := resultKey
		action := &ConstructorAction{
			ResultKey: resultKey,
			Constructor: func() interface{} {
				return r.Constructor(k)
			},
		}

		actions = append(actions, action)
	}

	return actions
}

// ConstructorAction is a mapping action that adds a new result object to the
// result set
type ConstructorAction struct {
	// ResultKey is the key that identifies the resulting object to be
	// instantiated
	ResultKey string

	// Constructor is the function that instantiates the new object
	Constructor func() interface{}
}

var _ mapping.MappingAction = &ConstructorAction{}

func (c *ConstructorAction) Apply(_ interface{}, resultSet mapping.ResultSet) error {
	resultSet.AddResult(c.ResultKey, c.Constructor())
	return nil
}

func FromLookup(key string) ResultKeyFunc {
	return func(tag reflect.StructTag) (string, bool) {
		return tag.Lookup(key)
	}
}
