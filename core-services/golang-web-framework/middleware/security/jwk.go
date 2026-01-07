package security

import (
	"context"
	"log"

	"github.com/lestrrat-go/jwx/v3/jwk"
)

type Jwk struct {
	Url string
}

func (receiver *Jwk) JwkSets() (jwk.Set, error) {
	set, err := jwk.Fetch(context.Background(), receiver.Url)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return set, err
}
