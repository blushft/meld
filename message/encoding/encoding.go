package encoding

import "io"

type Encoder interface {
	io.ReadWriteCloser
	Name() string
}

type Register func(encoder Encoder) error

func JSON() Encoder {
	return &jsonEnc{}
}

type jsonEnc struct {
}

func (j *jsonEnc) Name() string {
	return "json"
}

func (j *jsonEnc) Close() error {
	return nil
}

func (j *jsonEnc) Read(p []byte) (int, error) {
	return 0, nil
}

func (j *jsonEnc) Write(p []byte) (int, error) {
	return 0, nil
}
