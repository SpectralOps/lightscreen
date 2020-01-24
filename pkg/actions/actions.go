package actions

import (
	"github.com/jondot/lightscreen/pkg/core"
)

type Actions struct {
	MutationMap   map[string]core.ActionBuilder
	ValidationMap map[string]core.ActionBuilder
}

func Default() *Actions {
	MutationMap := map[string]core.ActionBuilder{
		"resolve_sha": NewResolveSHAMutation,
	}

	ValidationMap := map[string]core.ActionBuilder{
		"admit_sha": NewAdmitSHAValidation,
	}

	return &Actions{
		MutationMap,
		ValidationMap,
	}
}

func (a *Actions) Clear() {
	a.MutationMap = map[string]core.ActionBuilder{}
	a.ValidationMap = map[string]core.ActionBuilder{}
}

func (a *Actions) AddMutation(key string, mb core.ActionBuilder) {
	a.MutationMap[key] = mb
}

func (a *Actions) GetMutation(key string, config map[string]interface{}, opts *core.ActionContext) (core.Action, error) {
	return a.MutationMap[key](config, opts)
}

func (a *Actions) AddValidation(key string, vb core.ActionBuilder) {
	a.ValidationMap[key] = vb
}

func (a *Actions) GetValidation(key string, config map[string]interface{}, opts *core.ActionContext) (core.Action, error) {
	return a.ValidationMap[key](config, opts)
}
