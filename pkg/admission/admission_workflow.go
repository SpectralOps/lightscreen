package admission

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	sigadmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/jondot/lightscreen/pkg/actions"
	"github.com/jondot/lightscreen/pkg/core"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

type AdmissionConfig struct {
	Validations []map[string]interface{} `yaml:"validations"`
	Mutations   []map[string]interface{} `yaml:"mutations"`
}

type AdmissionMetrics struct {
	AdmissionStatus *prometheus.CounterVec
	Duration        *prometheus.HistogramVec
}
type AdmissionWorkflow struct {
	Mutations   []core.Action
	Validations []core.Action
	Logger      *zap.SugaredLogger
	Metrics     AdmissionMetrics
}

func NewAdmissionWorkflow(actionBuilders *actions.Actions, config io.Reader, logger *zap.SugaredLogger) (*AdmissionWorkflow, error) {
	m := AdmissionMetrics{
		AdmissionStatus: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "lightscreen_admission_status",
				Help: "Invalid admission request",
			},
			[]string{"status", "admitted"},
		),
		Duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "lightscreen_admission_duration",
				Help: "Duration of admission",
			},
			[]string{},
		),
	}
	a := &AdmissionWorkflow{Logger: logger, Metrics: m}
	opts := &core.ActionContext{Logger: logger, MetricsRegistry: metrics.Registry}

	data, err := ioutil.ReadAll(config)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	t := AdmissionConfig{}

	err = yaml.Unmarshal([]byte(data), &t)

	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	for _, m := range t.Mutations {
		mut, err := actionBuilders.GetMutation(m["type"].(string), m, opts)
		if err != nil {
			return nil, err
		}

		a.Mutations = append(a.Mutations, mut)
	}

	for _, m := range t.Validations {
		val, err := actionBuilders.GetValidation(m["type"].(string), m, opts)
		if err != nil {
			return nil, err
		}

		a.Validations = append(a.Validations, val)
	}

	return a, nil
}

func (a *AdmissionWorkflow) Execute(ctx context.Context, body []byte) sigadmission.Response {
	a.Logger.Debugf("Request %v", string(body))

	start := time.Now()

	defer func() { a.Metrics.Duration.WithLabelValues().Observe(time.Now().Sub(start).Seconds()) }()

	obj := &unstructured.Unstructured{}
	// metric defer timing for execute
	if _, _, err := deserializer.Decode(body, nil, obj); err != nil {
		a.Logger.Errorw("Bad request", "error", err)
		a.Metrics.AdmissionStatus.With(prometheus.Labels{"status": "bad_request", "admitted": "false"}).Inc()

		return sigadmission.Errored(http.StatusBadRequest, err)
	}

	originalObj := obj.DeepCopy()

	for _, m := range a.Mutations {
		err := m.Run(ctx, obj)
		if err != nil {
			a.Logger.Errorw("Mutation error", "error", err)
			a.Metrics.AdmissionStatus.With(prometheus.Labels{"status": "mutation_error", "admitted": "false"}).Inc()
			return sigadmission.Errored(http.StatusBadRequest, err)
		}
	}
	for _, v := range a.Validations {
		err := v.Run(ctx, obj)
		if err != nil {
			a.Logger.Errorw("Failed validation", "error", err)
			a.Metrics.AdmissionStatus.With(prometheus.Labels{"status": "validation_error", "admitted": "false"}).Inc()
			return sigadmission.ValidationResponse(false, err.Error())
		}
	}

	if len(a.Mutations) > 0 {
		originalBytes, _ := originalObj.MarshalJSON()
		objBytes, _ := obj.MarshalJSON()
		res := sigadmission.PatchResponseFromRaw(originalBytes, objBytes)
		a.Logger.Debugw("response", "response", res)
		a.Metrics.AdmissionStatus.With(prometheus.Labels{"status": "ok_with_patch", "admitted": "true"}).Inc()
		return res
	}
	res := sigadmission.ValidationResponse(true, "")
	a.Logger.Debugw("response", "response", res)
	a.Metrics.AdmissionStatus.With(prometheus.Labels{"status": "ok", "admitted": "true"}).Inc()
	return res
}
