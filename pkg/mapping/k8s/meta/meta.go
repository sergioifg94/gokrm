package meta

import (
	"github.com/sergioifg94/gokrm/pkg/mapping"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MetaMappingAction struct {
	MetaMapping MetaResourceMapping
}

type MetaResourceMapping interface {
	MapMeta(resourceKey string, sourceMeta, targetMeta metav1.Object) error
}

var _ mapping.MappingAction = &MetaMappingAction{}

func (a *MetaMappingAction) Apply(source interface{}, resourcePoolManager mapping.ResultSet) error {
	sourceMeta, err := meta.Accessor(source)
	if err != nil {
		return nil
	}

	for resourceKey, resource := range resourcePoolManager.GetAllResults() {
		resourceMeta, err := meta.Accessor(resource)
		if err != nil {
			continue
		}

		if err := a.MetaMapping.MapMeta(resourceKey, sourceMeta, resourceMeta); err != nil {
			return err
		}
	}

	return nil
}

type IdentityMetaResourceMapping struct{}

var _ MetaResourceMapping = &IdentityMetaResourceMapping{}

func (m *IdentityMetaResourceMapping) MapMeta(resourceKey string, sourceMeta, targetMeta metav1.Object) error {
	targetMeta.SetName(sourceMeta.GetName())
	targetMeta.SetNamespace(sourceMeta.GetNamespace())
	targetMeta.SetLabels(sourceMeta.GetLabels())
	targetMeta.SetAnnotations(sourceMeta.GetAnnotations())

	return nil
}

type MetaActionResolver struct {
	MetaMapping MetaResourceMapping
}

var _ mapping.MappingActionResolver = &MetaActionResolver{}

func (r *MetaActionResolver) ResolveMappingActions(source interface{}) ([]mapping.MappingAction, error) {
	_, err := meta.Accessor(source)
	if err != nil {
		return make([]mapping.MappingAction, 0), nil
	}

	return []mapping.MappingAction{
		&MetaMappingAction{
			MetaMapping: r.MetaMapping,
		},
	}, nil
}

type FieldMetaResourceMapping struct {
	MapName        func(string) string
	MapNamespace   func(string) string
	MapLabels      func(sourceLabels, targetLabels map[string]string)
	MapAnnotations func(sourceAnnotations, targetAnnotations map[string]string)
}

func IdentityString(v string) string     { return v }
func IdentityMap(_, _ map[string]string) {}

var _ MetaResourceMapping = &FieldMetaResourceMapping{}

func (m *FieldMetaResourceMapping) MapMeta(resourceKey string, sourceMeta, targetMeta metav1.Object) error {
	targetMeta.SetName(m.MapName(sourceMeta.GetName()))
	targetMeta.SetNamespace(m.MapNamespace(sourceMeta.GetNamespace()))

	if targetMeta.GetLabels() == nil {
		targetMeta.SetLabels(map[string]string{})
	}
	m.MapLabels(sourceMeta.GetLabels(), targetMeta.GetLabels())

	if targetMeta.GetAnnotations() == nil {
		targetMeta.SetAnnotations(map[string]string{})
	}
	m.MapAnnotations(sourceMeta.GetAnnotations(), targetMeta.GetAnnotations())

	return nil
}

func NewFieldMetaResourceMapping() *FieldMetaResourceMapping {
	return &FieldMetaResourceMapping{
		MapName:        IdentityString,
		MapNamespace:   IdentityString,
		MapLabels:      IdentityMap,
		MapAnnotations: IdentityMap,
	}
}
