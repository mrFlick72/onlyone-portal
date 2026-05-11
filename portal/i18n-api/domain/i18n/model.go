package i18n-api

import "strings"
type Locale string

func (l Locale) Normalized() Locale {
	return Locale(strings.ToLower(strings.ReplaceAll(string(l), "-", "_")))
}

func (l Locale) IsEmpty() bool {
	return strings.TrimSpace(string(l)) == ""
}

type BundleKey struct {
	Application string
	Page        string
	Language    Locale
}

type MessageBundle struct {
	Application string
	Page        string
	Language    Locale
	Messages    map[string]any
}
