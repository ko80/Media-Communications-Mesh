package registry

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Operation int

const (
	OpReadMany  Operation = iota // input: id range, output: array of struct
	OpReadOne                    // input: id,       output: struct
	OpAddOne                     // input: struct,   output: id
	OpDeleteOne                  // input: id,       output: none
)

type Request struct {
	Operation Operation
	ID        string
	Data      interface{}
	Flags     uint32
	Reply     chan Reply
}

type Reply struct {
	ID   string
	Data interface{}
	Err  error
}

func NewRequest(op Operation) Request {
	return Request{
		Operation: op,
		Reply:     make(chan Reply, 1),
	}
}

type RequestHandler interface {
	HandleReadMany(req Request)
	HandleReadOne(req Request, item interface{})
}

type Registry struct {
	queue   chan Request
	items   map[string]interface{}
	handler RequestHandler
}

func (r *Registry) Run(ctx context.Context) error {
	if r.handler == nil {
		return errors.New("no request handler assigned")
	}

	r.queue = make(chan Request, 100) // TODO: capacity to be configured
	r.items = make(map[string]interface{})

	for {
		select {
		case <-ctx.Done():
			return nil

		case req := <-r.queue:
			switch req.Operation {

			case OpReadMany:
				if req.ID != "" {
					req.Reply <- Reply{Err: errors.New("registry doesn't support read many over id range")}
					break
				}
				r.handler.HandleReadMany(req)

			case OpReadOne:
				item := r.items[req.ID]
				if item == nil {
					req.Reply <- Reply{Err: ErrResourceNotFound}
					break
				}
				r.handler.HandleReadOne(req, item)

			case OpAddOne:
				id := uuid.NewString()
				if _, ok := r.items[id]; ok {
					req.Reply <- Reply{Err: fmt.Errorf("new uuid already exists in registry (%v)", id)}
					break
				}
				r.items[id] = req.Data
				req.Reply <- Reply{ID: id}

			case OpDeleteOne:
				if _, ok := r.items[req.ID]; !ok {
					req.Reply <- Reply{Err: ErrResourceNotFound}
					break
				}
				delete(r.items, req.ID)
				req.Reply <- Reply{}
			}
		}
	}
}

func (r *Registry) ExecRequest(ctx context.Context, req Request) (Reply, error) {
	select {
	case <-ctx.Done():
		return Reply{}, ctx.Err()
	case r.queue <- req:
	}

	select {
	case <-ctx.Done():
		return Reply{}, ctx.Err()
	case reply := <-req.Reply:
		if reply.Err != nil {
			return Reply{}, reply.Err
		}
		return reply, nil
	}
}
