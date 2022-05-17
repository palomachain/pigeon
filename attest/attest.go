package attest

import "context"

type Request interface{}

type Evidence interface {
	Bytes() ([]byte, error)
}

type Registry struct {
	r map[string]Attestor
}

func NewRegistry() *Registry {
	return &Registry{
		r: make(map[string]Attestor),
	}
}

func (r *Registry) Register(queueTypeName string, att Attestor) {
	r.r[queueTypeName] = att
}

func (r *Registry) Execute(ctx context.Context, queueTypeName string, req Request) (Evidence, error) {
	att, ok := r.r[queueTypeName]
	if !ok {
		return nil, nil
	}
	return att.ProvideEvidence(ctx, req)
}

type Attestor interface {
	ProvideEvidence(ctx context.Context, req Request) (Evidence, error)
}
