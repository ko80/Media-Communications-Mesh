package registry

import (
	"context"
	"control-plane-agent/internal/model"
	"fmt"
)

type multipointGroupRegistry struct {
	Registry
}

var MultipointGroupRegistry multipointGroupRegistry

func (r *multipointGroupRegistry) Init() {
	r.handler = r
}

func (r *multipointGroupRegistry) HandleReadMany(req Request) {
	out := []model.MultipointGroup{}
	for id, item := range r.items {
		group, ok := item.(model.MultipointGroup)
		if !ok {
			req.Reply <- Reply{Err: ErrTypeCastFailed}
			return
		}
		group.ID = id
		if group.Status != nil {
			if req.Flags&FlagAddStatus != 0 {
				v := *group.Status
				group.Status = &v
			} else {
				group.Status = nil
			}
		}
		if group.Config != nil {
			if req.Flags&FlagAddConfig != 0 {
				v := *group.Config
				group.Config = &v
			} else {
				group.Config = nil
			}
		}
		out = append(out, group)
	}
	req.Reply <- Reply{Data: out}
}

func (r *multipointGroupRegistry) ListAll(ctx context.Context, addStatus, addConfig bool) ([]model.MultipointGroup, error) {
	req := NewRequest(OpReadMany)
	if addStatus {
		req.Flags |= FlagAddStatus
	}
	if addConfig {
		req.Flags |= FlagAddConfig
	}
	reply, err := r.ExecRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("multipoint group registry list all exec req err: %w", err)
	}
	ret, ok := reply.Data.([]model.MultipointGroup)
	if !ok {
		return nil, fmt.Errorf("multipoint group registry list all err: %w", ErrTypeCastFailed)
	}
	return ret, nil
}

func (r *multipointGroupRegistry) HandleReadOne(req Request, item interface{}) {
	group, ok := item.(model.MultipointGroup)
	if !ok {
		req.Reply <- Reply{Err: ErrTypeCastFailed}
		return
	}
	group.ID = req.ID
	if group.Status != nil {
		v := *group.Status
		group.Status = &v
	}
	if group.Config != nil {
		if req.Flags&FlagAddConfig != 0 {
			v := *group.Config
			group.Config = &v
		} else {
			group.Config = nil
		}
	}
	req.Reply <- Reply{Data: group}
}

func (r *multipointGroupRegistry) Get(ctx context.Context, id string, addConfig bool) (model.MultipointGroup, error) {
	req := NewRequest(OpReadOne)
	req.ID = id
	if addConfig {
		req.Flags |= FlagAddConfig
	}
	reply, err := r.ExecRequest(ctx, req)
	if err != nil {
		return model.MultipointGroup{}, fmt.Errorf("multipoint group registry get exec req err: %w", err)
	}
	ret, ok := reply.Data.(model.MultipointGroup)
	if !ok {
		return model.MultipointGroup{}, fmt.Errorf("multipoint group registry get err: %w", ErrTypeCastFailed)
	}
	return ret, nil
}

func (r *multipointGroupRegistry) Add(ctx context.Context, conn model.MultipointGroup) (string, error) {
	req := NewRequest(OpAddOne)
	req.Data = conn
	reply, err := r.ExecRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("multipoint group registry add exec req err: %w", err)
	}
	return reply.ID, nil
}

func (r *multipointGroupRegistry) Delete(ctx context.Context, id string) error {
	req := NewRequest(OpDeleteOne)
	req.ID = id
	_, err := r.ExecRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("multipoint group registry del exec req err: %w", err)
	}
	return nil
}
