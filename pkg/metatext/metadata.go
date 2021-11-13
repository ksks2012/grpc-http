package metatext

import (
	"log"
	"strings"

	"google.golang.org/grpc/metadata"
)

type MetaDataTextMap struct {
	metadata.MD
}

func (m MetaDataTextMap) ForeachKey(handler func(key, val string) error) error {
	log.Printf("ForeachKey")
	for k, vs := range m.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m MetaDataTextMap) Set(key, val string) {
	log.Printf("Set")
	key = strings.ToLower(key)
	m.MD[key] = append(m.MD[key], val)
}
