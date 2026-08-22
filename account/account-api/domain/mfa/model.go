package mfa

type MfaDevice struct {
	UserName    string `json:"userName"`
	MfaMethod   string `json:"mfaMethod"`
	MfaChannel  string `json:"mfaChannel"`
	MfaDeviceId string `json:"mfaDeviceId"`
	Default     bool   `json:"default"`
}
