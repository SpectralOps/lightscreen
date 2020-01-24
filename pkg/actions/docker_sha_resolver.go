package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/jondot/lightscreen/pkg/core"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ResolveSHAAction struct {
	resolver *DockerShaResolver
	opts     *core.ActionContext
}

func NewResolveSHAMutation(config map[string]interface{}, opts *core.ActionContext) (core.Action, error) {
	val, ok := config["finder"]
	var finder DigestFinder = &CraneDigestFinder{}
	if ok && val.(string) == "dict" {
		conf, ok := config["finder_config"]
		if ok {
			finder = &DictDigestFinder{entries: makeConfigMap(conf.(map[interface{}]interface{}))}
		}
	}
	resolver := NewDockerShaResolver(finder)
	return &ResolveSHAAction{
		resolver: resolver,
		opts:     opts,
	}, nil
}

func (r *ResolveSHAAction) Name() string {
	return "resolve_sha"
}

func (r *ResolveSHAAction) Run(ctx context.Context, p *unstructured.Unstructured) error {
	o, err := gabs.Consume(p.Object)
	if err != nil {
		r.opts.Logger.Infow("Cannot understand this request", "object", p.Object)
		return errors.New("Bad request")
	}
	containers := o.Path("spec.containers").Data().([]interface{})

	for _, c := range containers {
		container, _ := gabs.Consume(c)
		img := container.Path("image").Data()
		if img == nil {
			r.opts.Logger.Infow("Theres no image for this container", "container", c)
			return errors.New(fmt.Sprintf("There is no image for container %v", c))
		}
		resolved, err := r.resolver.Resolve(img.(string))
		if err != nil {
			return err
		}
		container.Set(resolved, "image")
	}
	return nil
}

func NewDockerShaResolver(finder DigestFinder) *DockerShaResolver {
	return &DockerShaResolver{finder}
}

func (r *DockerShaResolver) Resolve(tag string) (string, error) {
	if !strings.Contains(tag, "/") {
		tag = "library/" + tag
	}
	m, err := r.finder.Digest(tag)
	if err != nil {
		return "", err
	}

	return tag + "@" + m, nil
}

type DigestFinder interface {
	Digest(string) (string, error)
}

type DictDigestFinder struct {
	entries map[string]string
}

func (d *DictDigestFinder) Digest(ref string) (string, error) {
	val, ok := d.entries[ref]
	if !ok {
		return "", errors.Errorf("No digest found for %v", ref)
	}
	return val, nil
}

type CraneDigestFinder struct {
}

func (_ *CraneDigestFinder) Digest(ref string) (string, error) {
	return crane.Digest(ref)
}

type DockerShaResolver struct {
	finder DigestFinder
}
