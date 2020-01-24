package core

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ActionBuilder func(config map[string]interface{}, opts *ActionContext) (Action, error)

type Action interface {
	Name() string
	Run(context.Context, *unstructured.Unstructured) error
}

type ActionContext struct {
	Logger          *zap.SugaredLogger
	MetricsRegistry *prometheus.Registry
}
