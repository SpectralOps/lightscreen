package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/jondot/lightscreen/pkg/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type AdmitSHAValidation struct {
	refs map[string]string
	opts *core.ActionContext
}

func NewAdmitSHAValidation(config map[string]interface{}, opts *core.ActionContext) (core.Action, error) {
	mapping, ok := config["admit"]
	if !ok {
		return nil, errors.New("Missing 'admit' mapping in configuration")
	}
	refs := makeConfigMap(mapping.(map[interface{}]interface{}))
	return &AdmitSHAValidation{
		refs: refs,
		opts: opts,
	}, nil
}

func (a *AdmitSHAValidation) Name() string {
	return "admit_sha"
}
func (a *AdmitSHAValidation) Run(ctx context.Context, p *unstructured.Unstructured) error {
	o, err := gabs.Consume(p.Object)
	if err != nil {
		a.opts.Logger.Infow("Cannot understand this request", "object", p.Object)
		return errors.New("Bad request")
	}
	containers := o.Path("spec.containers").Data().([]interface{})

	for _, c := range containers {
		container, _ := gabs.Consume(c)
		img := container.Path("image").Data()
		if img == nil {
			a.opts.Logger.Infow("Theres no image for this container", "container", c)
			return errors.New(fmt.Sprintf("There is no image for container %v", c))
		}
		_, ok := a.refs[img.(string)]

		if !ok {
			a.opts.Logger.Infow("Image not allowed", "image", img)
			return errors.New(fmt.Sprintf("Image %v is not allowed to be admitted", img))
		}
	}
	return nil
}
