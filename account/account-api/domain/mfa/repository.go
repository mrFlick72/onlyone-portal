package mfa

import "context"

type MfaRepository interface {
	FindAll(ctx context.Context) ([]MfaDevice, error)
	StartEnrollment(ctx context.Context, mfaMethod string, mfaChannel string) (string, error)
	Associate(ctx context.Context, ticket string, code string) error
	SetDefault(ctx context.Context, mfaDeviceId string) error
}
