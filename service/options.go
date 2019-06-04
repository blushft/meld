package service

import (
	"github.com/blushft/meld/common/metadata"

	"github.com/blang/semver"
	"github.com/blushft/meld/connector"
)

type Options struct {
	Name      string
	Namespace string
	Version   semver.Version
	Labels    map[string]string
	Tags      []string

	Connectors map[string]connector.Connector
	Meta       metadata.Metadata
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	options := &Options{
		Labels:     make(map[string]string),
		Tags:       make([]string, 0),
		Connectors: make(map[string]connector.Connector),
	}
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func Connector(c connector.Connector) Option {
	return func(o *Options) {
		o.Connectors[c.Name()] = c
	}
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}

func Version(v string) Option {
	return func(o *Options) {
		o.Version = semver.MustParse(v)
	}
}

func WithLabel(k, v string) Option {
	return func(o *Options) {
		o.Labels[k] = v
	}
}

func WithLabels(lab map[string]string) Option {
	return func(o *Options) {
		for k, v := range lab {
			o.Labels[k] = v
		}
	}
}

func WithTag(tag ...string) Option {
	return func(o *Options) {
		o.Tags = append(o.Tags, tag...)
	}
}
