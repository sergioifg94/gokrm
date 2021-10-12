package mapping

import (
	"fmt"
	"testing"

	"github.com/sergioifg94/gokrm/pkg/mapping/k8s/meta"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestMapObject(t *testing.T) {
	assert := assert.New(t)

	var replicas int32 = 10

	results, err := MapObject(
		map[ResultKey]runtime.Object{
			"Deployment": &appsv1.Deployment{},
			"Service":    &corev1.Service{},
		},
		func(aro *ActionResolverOptions) {
			aro.Meta.MetaMapping.(*meta.FieldMetaResourceMapping).MapNamespace = func(s string) string {
				return fmt.Sprintf("gokrm-%s", s)
			}
		},
		sourceType{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Replicas:  &replicas,
			ClusterIP: "1.1.1.1",
			Image:     "image:latest",
			Port:      8080,
		},
	)
	assert.NoError(err, "No error expected")

	resultDeployment := results["Deployment"].(*appsv1.Deployment)
	resultService := results["Service"].(*corev1.Service)

	assert.Equal("test", resultDeployment.Name)
	assert.Equal("gokrm-default", resultDeployment.Namespace)
	assert.Equal(int32(10), *resultDeployment.Spec.Replicas)
	assert.Equal("image:latest", resultDeployment.Spec.Template.Spec.Containers[0].Image)

	assert.Equal("test", resultService.Name)
	assert.Equal("gokrm-default", resultService.Namespace)
	assert.Equal("1.1.1.1", resultService.Spec.ClusterIP)
	assert.Equal(int32(8080), resultService.Spec.Ports[0].Port)
}

type sourceType struct {
	metav1.ObjectMeta
	Image     string `gokrmTarget:"Deployment" gokrmTargetField:"Spec.Template.Spec.Containers[0].Image"`
	Replicas  *int32 `gokrmTarget:"Deployment" gokrmTargetField:"Spec.Replicas"`
	ClusterIP string `gokrmTarget:"Service" gokrmTargetField:"Spec.ClusterIP"`
	Port      int32  `gokrmTarget:"Service" gokrmTargetField:"Spec.Ports[0].Port"`
}

var _ runtime.Object = sourceType{}

func (s sourceType) DeepCopyObject() runtime.Object {
	return s
}

func (s sourceType) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}

var _ metav1.ObjectMetaAccessor = &sourceType{}

func (s sourceType) GetObjectMeta() metav1.Object {
	return &s.ObjectMeta
}
