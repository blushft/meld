package codec

type Codec interface{}

type Register func(codec Codec) error
