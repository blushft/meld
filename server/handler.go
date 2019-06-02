package server

import (
	"github.com/blushft/meld/common/metadata"
)

type HandlerOption func(*HandlerOptions)

type HandlerOptions struct {
	Name string
	Meta metadata.Metadata
}
