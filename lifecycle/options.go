package lifecycle

import "context"

type LifecycleOptions struct {
	Context context.Context

	BeforeInit  []func(c context.Context) error
	BeforeStart []func(c context.Context) error
	BeforeStop  []func(c context.Context) error
	AfterInit   []func(c context.Context) error
	AfterStart  []func(c context.Context) error
	AfterStop   []func(c context.Context) error

	start func(context.Context) error
	stop  func(context.Context) error
}

type LifecycleOption func(*LifecycleOptions)

func newOptions(opts ...LifecycleOption) *LifecycleOptions {
	options := &LifecycleOptions{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(options)
	}

	return options
}

func Context(ctx context.Context) LifecycleOption {
	return func(o *LifecycleOptions) {
		o.Context = ctx
	}
}

func BeforeInit(f ...func(c context.Context) error) LifecycleOption {
	return func(o *LifecycleOptions) {
		o.BeforeInit = append(o.BeforeInit, f...)
	}
}

func BeforeStart(f ...func(c context.Context) error) LifecycleOption {
	return func(o *LifecycleOptions) {
		o.BeforeStart = append(o.BeforeStart, f...)
	}
}

func BeforeStop(f ...func(c context.Context) error) LifecycleOption {
	return func(o *LifecycleOptions) {
		o.BeforeStop = append(o.BeforeStop, f...)
	}

}

func AfterInit(f ...func(c context.Context) error) LifecycleOption {
	return func(o *LifecycleOptions) {
		o.AfterInit = append(o.AfterInit, f...)
	}
}

func AfterStart(f ...func(c context.Context) error) LifecycleOption {
	return func(o *LifecycleOptions) {
		o.AfterStart = append(o.AfterStart, f...)
	}
}

func AfterStopCE(f ...func(c context.Context) error) LifecycleOption {
	return func(o *LifecycleOptions) {
		o.AfterStop = append(o.AfterStop, f...)
	}
}
