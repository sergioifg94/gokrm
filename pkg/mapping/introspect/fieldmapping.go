package introspect

import (
	"reflect"

	"github.com/sergioifg94/gokrm/pkg/mapping"
)

type FieldMappingActionResolver struct {
	ActionBuilder FieldMappingActionBuilder
}

type FieldMappingActionBuilder interface {
	ActionFor(field reflect.StructField, fieldSelector FieldSelector) (mapping.MappingAction, bool)
}

type FieldSelector func(source reflect.Value) reflect.Value

func FieldSelectorID(source reflect.Value) reflect.Value {
	return source
}

var _ mapping.MappingActionResolver = &FieldMappingActionResolver{}

func (r *FieldMappingActionResolver) ResolveMappingActions(source interface{}) ([]mapping.MappingAction, error) {
	// Use reflection to get the fields of the source type
	sourceValue := reflect.ValueOf(source)
	sourceType := sourceValue.Type()

	return r.resolveActionsFor(sourceType, FieldSelectorID), nil
}

func (r *FieldMappingActionResolver) resolveActionsFor(sourceType reflect.Type, currentSelector func(reflect.Value) reflect.Value) []mapping.MappingAction {
	// Initialize resulting actions
	actions := []mapping.MappingAction{}

	// Iterate through the source object fields and check the tag information
	// to generate new mapping actions from them
	for i := 0; i < sourceType.NumField(); i++ {
		fieldIndex := i
		field := sourceType.Field(i)
		fieldSelector := func(s reflect.Value) reflect.Value {
			currentValue := currentSelector(s)
			return currentValue.Field(fieldIndex)
		}

		action, ok := r.ActionBuilder.ActionFor(field, fieldSelector)
		if ok {
			actions = append(actions, action)
		}

		if field.Type.Kind() == reflect.Struct {
			actions = append(actions, r.resolveActionsFor(field.Type, fieldSelector)...)
		}
	}

	return actions
}
