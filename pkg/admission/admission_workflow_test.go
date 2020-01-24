package admission_test

import (
	"context"
	"strings"
	"testing"

	"encoding/json"

	"github.com/BTBurke/snapshot"
	"github.com/jondot/lightscreen/pkg/actions"
	. "github.com/jondot/lightscreen/pkg/admission"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestAdmissionWorkflow(t *testing.T) {
	log_, _ := zap.NewDevelopment()
	logger := log_.Sugar()

	t.Run("mapping=resolve-but-NO-admission", func(t *testing.T) {
		conf := `mutations:
  - type: resolve_sha
    finder: dict
    finder_config:
      "library/nginx": sha256:deadbeef
validations:
  - type: admit_sha
    admit:
      "sha256:3f": true`
		wf, err := NewAdmissionWorkflow(actions.Default(), strings.NewReader(conf), logger)
		pod := corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx"},
				},
			},
		}
		body, _ := json.Marshal(&pod)
		ar := admissionv1beta1.AdmissionReview{
			Request: &admissionv1beta1.AdmissionRequest{
				Kind:   metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
				Object: runtime.RawExtension{Raw: body},
			},
		}
		as := assert.New(t)
		req, err := json.Marshal(&ar)
		as.Nil(err)

		resp := wf.Execute(context.TODO(), req)
		as.Nil(err)

		snapshot.Assert(t, map[string]interface{}{
			"patches": resp.Patches,
			"res":     resp.AdmissionResponse,
		})
	})
	t.Run("mapping=resolve-YES-admission", func(t *testing.T) {
		conf := `mutations:
  - type: resolve_sha
    finder: dict
    finder_config:
      "library/nginx": sha256:deadbeef
validations:
  - type: admit_sha
    admit:
      "library/nginx@sha256:deadbeef": true`
		wf, err := NewAdmissionWorkflow(actions.Default(), strings.NewReader(conf), logger)
		pod := corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx"},
				},
			},
		}
		body, _ := json.Marshal(&pod)
		ar := admissionv1beta1.AdmissionReview{
			Request: &admissionv1beta1.AdmissionRequest{
				Kind:   metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
				Object: runtime.RawExtension{Raw: body},
			},
		}
		as := assert.New(t)
		req, err := json.Marshal(&ar)
		as.Nil(err)

		resp := wf.Execute(context.TODO(), req)
		as.Nil(err)

		snapshot.Assert(t, map[string]interface{}{
			"patches": resp.Patch,
			"res":     resp.AdmissionResponse,
		})
	})
}
