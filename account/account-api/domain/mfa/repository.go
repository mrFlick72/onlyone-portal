package mfa

import "context"

type MfaRepository interface {
	FindAll(ctx context.Context) ([]MfaDevice, error)
	StartEnrollment(ctx context.Context, mfaMethod string, mfaChannel string) (string, error)
	Associate(ctx context.Context, ticket string, code string) error
}
