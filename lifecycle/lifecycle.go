package lifecycle

import (
	"context"
	"sync"
)

type Starter interface {
	Start(context.Context) error
}

type Checker interface {
	Check(context.Context) error
}

type Stopper interface {
	Stop(context.Context) error
}

type Initializer interface {
	Init(context.Context) error
}

type Lifecycle interface {
	Starter
	Stopper
	Checker
	Initializer

	Options() *LifecycleOptions
}

func NewLifecycle(opts ...LifecycleOption) Lifecycle {
	options := newOptions(opts...)

	return &lifecycle{
		opts: options,
	}
}

type LifecycleWrapper func(lc Lifecycle) Lifecycle

type lifecycle struct {
	sync.RWMutex
	started bool

	opts *LifecycleOptions
}

func (lc *lifecycle) Start(ctx context.Context) error {
	lc.Lock()
	defer lc.Unlock()
	for _, bs := range lc.opts.BeforeStart {
		if err := bs(lc.opts.Context); err != nil {
			return err
		}
	}

	if err := lc.opts.start(lc.opts.Context); err != nil {
		return err
	}

	for _, as := range lc.opts.AfterStart {
		if err := as(lc.opts.Context); err != nil {
			return err
		}
	}

	lc.started = true
	return nil
}

func (lc *lifecycle) Stop(ctx context.Context) error {
	lc.Lock()
	defer lc.Unlock()

	return lc.opts.stop(ctx)
}

func (lc *lifecycle) Check(ctx context.Context) error {
	panic("not implmented")
}

func (lc *lifecycle) Init(ctx context.Context) error {
	panic("not implemented")
}

func (lc *lifecycle) Options() *LifecycleOptions {
	return lc.opts
}
