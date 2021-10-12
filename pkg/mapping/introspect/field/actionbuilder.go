package field

import (
	"reflect"

	"github.com/sergioifg94/gokrm/pkg/mapping"
	"github.com/sergioifg94/gokrm/pkg/mapping/introspect"
)

type TagMappingActionBuilder struct {
	ResultKeyTag,
	TargetSelectorTag string
}

var _ introspect.FieldMappingActionBuilder = &TagMappingActionBuilder{}

func (b *TagMappingActionBuilder) ActionFor(field reflect.StructField, fieldSelector introspect.FieldSelector) (mapping.MappingAction, bool) {
	tag := field.Tag

	resultKey, ok := tag.Lookup(b.ResultKeyTag)
	if !ok {
		return nil, false
	}

	targetSelector, ok := tag.Lookup(b.TargetSelectorTag)
	if !ok {
		targetSelector = field.Name
	}

	return &FieldMappingAction{
		ResultKey:      resultKey,
		TargetSelector: targetSelector,
		SourceSelector: fieldSelector,
	}, true
}
