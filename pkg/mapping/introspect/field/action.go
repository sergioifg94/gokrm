package field

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/sergioifg94/gokrm/pkg/mapping"
	"github.com/sergioifg94/gokrm/pkg/mapping/introspect"
)

type FieldMappingAction struct {
	ResultKey      string
	TargetSelector string
	SourceSelector func(reflect.Value) reflect.Value
}

var _ mapping.MappingAction = &FieldMappingAction{}

func (a *FieldMappingAction) Apply(source interface{}, resultSet mapping.ResultSet) error {
	result, ok := resultSet.GetResult(a.ResultKey)
	if !ok {
		return fmt.Errorf("result %s not found in result set", a.ResultKey)
	}

	sourceField := a.SourceSelector(reflect.ValueOf(source))
	targetField := ParseTargetSelector(a.TargetSelector)(reflect.ValueOf(result))

	targetField.Set(sourceField)
	return nil
}

func ParseTargetSelector(targetSelector string) func(reflect.Value) reflect.Value {
	fields := strings.Split(targetSelector, ".")

	selector := indexTargetSelector(func(source reflect.Value) reflect.Value {
		return source
	}, fields[0])

	if len(fields) == 1 {
		return selector
	}

	for i := 1; i < len(fields); i++ {
		field := fields[i]

		selector = indexTargetSelector(selector, field)
	}

	return selector
}

func indexTargetSelector(current introspect.FieldSelector, field string) introspect.FieldSelector {
	fieldName, index, ok := indexedField(field)
	if !ok {
		return fieldSelectorFor(current, field)
	}

	return func(source reflect.Value) reflect.Value {
		fieldValue := fieldSelectorFor(current, fieldName)(source)
		if fieldValue.IsNil() {
			slice := reflect.MakeSlice(fieldValue.Type(), index+1, index+1)
			fieldValue.Set(slice)
		} else if fieldValue.Len() < index+1 {
			fieldValue.SetLen(index + 1)
		}

		return fieldValue.Index(index)
	}
}

func fieldSelectorFor(current introspect.FieldSelector, fieldName string) introspect.FieldSelector {
	return func(source reflect.Value) reflect.Value {
		currentValue := current(source)

		if currentValue.Kind() == reflect.Ptr {
			currentValue = currentValue.Elem()
		}

		field, _ := currentValue.Type().FieldByName(fieldName)
		fieldValue := currentValue.FieldByName(fieldName)

		if field.Type.Kind() == reflect.Ptr {
			fieldType := field.Type.Elem()

			if fieldValue == (reflect.Value{}) {
				fieldValue.Set(reflect.New(fieldType))
			}
		}

		return fieldValue
	}
}

func indexedField(input string) (field string, index int, ok bool) {
	var re = regexp.MustCompile(`(?m)(?P<Field>[a-zA-Z0-9_]+)\[(?P<Index>[0-9]+)\]`)

	if !re.MatchString(input) {
		return
	}

	matches := re.FindStringSubmatch(input)
	field = matches[1]
	index, _ = strconv.Atoi(matches[2])
	ok = true
	return
}
