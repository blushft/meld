package options

import (
	"errors"

	"github.com/blushft/meld/utility"
)

type Option struct {
	Name  string
	Value interface{}
}

type Options []Option

type Setter func(Options)

func (opts Options) Set(name string, val interface{}) {
	opts = append(opts, Option{name, val})
}

func (opts Options) Value(name string) interface{} {
	for _, opt := range opts {
		if opt.Name == name {
			return opt.Value
		}
	}

	return nil
}

func (opts Options) Bind(name string, v interface{}) error {
	o := opts.Value(name)
	if v == nil {
		return errors.New("no option with name " + name + " could be found")
	}
	return utility.BindInterface(o, v)
}

func Name(n string) Setter {
	return func(o Options) {
		if val := o.Value("Name"); val == nil {
			val = n
			return
		}

		o = append(o, Option{"Name", n})
	}
}
