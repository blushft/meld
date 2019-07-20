package message

type Value interface {
	Interface() interface{}
	String() string
	Encode() []byte
}

type value struct {
	v interface{}
}

func NewValue(v interface{}) Value {
	return &value{
		v: v,
	}
}

func (v *value) Interface() interface{} {
	return v.v
}

func (v *value) String() string {
	if s, ok := v.v.(string); ok {
		return s
	}

	return ""
}

func (v *value) Encode() []byte {
	return []byte{}
}
