package endpoint

type Headers map[string]string

func NewHeaders(m ...map[string]string) Headers {
	h := make(map[string]string)

	for _, mm := range m {
		for k, v := range mm {
			h[k] = v
		}
	}

	return h
}
