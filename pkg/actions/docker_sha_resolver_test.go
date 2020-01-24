package actions_test

import (
	"context"
	"testing"

	"github.com/BTBurke/snapshot"
	. "github.com/jondot/lightscreen/pkg/actions"
	"github.com/jondot/lightscreen/pkg/core"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDockerSHAResolver(t *testing.T) {
	log_, _ := zap.NewDevelopment()
	logger := log_.Sugar()

	t.Run("mapping=missing-ref-c2", func(t *testing.T) {
		actx := &core.ActionContext{Logger: logger}
		b := []byte(`name: resolve_sha
finder: dict
finder_config:
  "library/nginx": "sha256:2f2f2f"`)
		c := map[string]interface{}{}
		err := yaml.Unmarshal(b, &c)
		resolveSHAMutation, _ := NewResolveSHAMutation(c, actx)

		pod := corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind: "pod",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Name: "nginx"},
					{Name: "c2"},
				},
			},
		}
		err = resolveSHAMutation.Run(context.TODO(), unstructure(pod))
		as := assert.New(t)
		as.Error(err)
	})
	t.Run("mapping=resolving-well", func(t *testing.T) {
		b := []byte(`name: resolve_sha
finder: dict
finder_config:
  "library/nginx": "sha256:2f2f2f"`)
		actx := &core.ActionContext{Logger: logger}
		c := map[string]interface{}{}
		err := yaml.Unmarshal(b, &c)
		resolveSHAMutation, _ := NewResolveSHAMutation(c, actx)

		pod := corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind: "pod",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Name: "nginx", Image: "nginx"},
				},
			},
		}

		u := unstructure(pod)
		err = resolveSHAMutation.Run(context.TODO(), u)
		as := assert.New(t)
		as.Nil(err)

		s, _ := u.MarshalJSON()
		snapshot.Assert(t, string(s))
	})
}
