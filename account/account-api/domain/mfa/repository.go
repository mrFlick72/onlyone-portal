package mfa

import "context"

type MfaRepository interface {
	FindAll(ctx context.Context) ([]MfaDevice, error)
}
