package mapping

import (
	gokrmapping "github.com/sergioifg94/gokrm/pkg/mapping"
	"github.com/sergioifg94/gokrm/pkg/mapping/introspect"
	"github.com/sergioifg94/gokrm/pkg/mapping/introspect/field"
	"github.com/sergioifg94/gokrm/pkg/mapping/k8s/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	RESOURCE_KEY_TAG = "gokrmTarget"
	TARGET_FIELD_TAG = "gokrmTargetField"
)

type ActionResolverOptions struct {
	Identification introspect.IdentificationActionResolver
	Meta           meta.MetaActionResolver
	FieldMapping   introspect.FieldMappingActionResolver
}

type ActionResolverOptionsFunc func(*ActionResolverOptions)

func DefaultOptions(_ *ActionResolverOptions) {}

type ResultKey string

func MapObject(resourceTemplates map[ResultKey]runtime.Object, options func(*ActionResolverOptions), source interface{}) (map[ResultKey]runtime.Object, error) {
	resultSet := gokrmapping.NewMappedResultSet()

	actionResolvers := &ActionResolverOptions{
		Identification: introspect.IdentificationActionResolver{
			ResultKey: introspect.FromLookup(RESOURCE_KEY_TAG),
			Constructor: func(resultKey string) interface{} {
				return resourceTemplates[ResultKey(resultKey)].DeepCopyObject()
			},
		},
		Meta: meta.MetaActionResolver{
			MetaMapping: meta.NewFieldMetaResourceMapping(),
		},
		FieldMapping: introspect.FieldMappingActionResolver{
			ActionBuilder: &field.TagMappingActionBuilder{
				ResultKeyTag:      RESOURCE_KEY_TAG,
				TargetSelectorTag: TARGET_FIELD_TAG,
			},
		},
	}

	options(actionResolvers)

	if err := gokrmapping.MapObject(
		gokrmapping.ActionResolverComposedOf(
			&actionResolvers.Identification,
			&actionResolvers.Meta,
			&actionResolvers.FieldMapping,
		),
		resultSet,
		source,
	); err != nil {
		return nil, err
	}

	resultObjs := map[ResultKey]runtime.Object{}
	for resourceKey, result := range resultSet.GetAllResults() {
		resultObjs[ResultKey(resourceKey)] = result.(runtime.Object)
	}

	return resultObjs, nil
}
