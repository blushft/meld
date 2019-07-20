package encoding

type Encoder interface{}

type Register func(encoder Encoder) error
