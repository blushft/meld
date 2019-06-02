package metadata

type Metadata map[string]map[string]string

func NewMetadata() Metadata {
	return map[string]map[string]string{}
}
