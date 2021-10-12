package field

import (
	"reflect"
	"testing"

	"github.com/sergioifg94/gokrm/pkg/mapping"
)

func TestFieldMappingAction(t *testing.T) {
	type sourceType struct {
		Foo string
		Bar string
	}

	type level2 struct {
		Field2 string
	}

	type level1 struct {
		Field1 string
		Level2 []level2
	}

	type targetType struct {
		Field  string
		Level1 level1
	}

	action := &FieldMappingAction{
		ResultKey:      "test",
		TargetSelector: "Level1.Level2[1].Field2",
		SourceSelector: func(v reflect.Value) reflect.Value {
			return v.FieldByName("Foo")
		},
	}

	source := sourceType{
		Foo: "testing",
		Bar: "hello",
	}

	target := targetType{}
	resultSet := mapping.NewMappedResultSet()
	resultSet.AddResult("test", &target)

	err := action.Apply(source, resultSet)
	if err != nil {
		t.Fatal(err)
	}
}
