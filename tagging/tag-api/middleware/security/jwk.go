package security

import (
	"context"

	"github.com/lestrrat-go/jwx/v3/jwk"
)

type Jwk struct {
	Url string
}

func (receiver *Jwk) JwkSets() (jwk.Set, error) {
	set, err := jwk.Fetch(context.Background(), receiver.Url)

	if err != nil {
		return nil, err
	}
	return set, err
}
