package actions_test

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/jondot/lightscreen/pkg/actions"
	"github.com/jondot/lightscreen/pkg/core"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func unstructure(p corev1.Pod) *unstructured.Unstructured {
	out, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	u := &unstructured.Unstructured{}
	err = json.Unmarshal(out, u)
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("u is nil")
	}
	return u

}

func TestAdmitSHA(t *testing.T) {
	log_, _ := zap.NewDevelopment()
	logger := log_.Sugar()
	actx := &core.ActionContext{Logger: logger}
	t.Run("mapping=no-match", func(t *testing.T) {
		a, _ := NewAdmitSHAValidation(map[string]interface{}{"admit": map[interface{}]interface{}{
			"foo": "bar",
		}}, actx)
		pod := corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind: "pod",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Name: "c1", Image: "c1"},
					{Name: "c2", Image: "c1", SecurityContext: &v1.SecurityContext{}},
				},
				InitContainers: []v1.Container{
					{Name: "i1", Image: "i1"},
					{Name: "i2", Image: "i1", SecurityContext: &v1.SecurityContext{}},
				},
			},
		}
		err := a.Run(context.TODO(), unstructure(pod))

		as := assert.New(t)
		as.Error(err)
	})
	t.Run("mapping=full-match", func(t *testing.T) {
		actx := &core.ActionContext{Logger: logger}
		a, _ := NewAdmitSHAValidation(map[string]interface{}{"admit": map[interface{}]interface{}{
			"sha256:2f2f": "true",
		}}, actx)
		pod := corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind: "pod",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Name: "c1", Image: "sha256:2f2f"},
				},
			},
		}
		err := a.Run(context.TODO(), unstructure(pod))

		as := assert.New(t)
		as.Nil(err)
	})
	t.Run("mapping=missing-one-container", func(t *testing.T) {
		actx := &core.ActionContext{Logger: logger}
		a, _ := NewAdmitSHAValidation(map[string]interface{}{"admit": map[interface{}]interface{}{
			"sha256:2f2f": "true",
		}}, actx)
		pod := corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind: "pod",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Name: "c1", Image: "sha256:2f2f"},
					{Name: "c2", Image: "sha256:2222"},
				},
			},
		}
		err := a.Run(context.TODO(), unstructure(pod))

		as := assert.New(t)
		as.Error(err)
	})
	t.Run("mapping=illegal-map", func(t *testing.T) {
		actx := &core.ActionContext{Logger: logger}
		as := assert.New(t)
		_, err := NewAdmitSHAValidation(map[string]interface{}{}, actx)
		as.Error(err)
	})
}
