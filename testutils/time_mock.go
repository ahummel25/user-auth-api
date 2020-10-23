//go:build !prod && !dev
// +build !prod,!dev

package testutils

import "time"

type TimeProvider interface {
	Now() time.Time
}

type MockTimeProvider struct {
	FrozenTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.FrozenTime
}

func NewMockTimeProvider(t time.Time) TimeProvider {
	return &MockTimeProvider{FrozenTime: t}
}

var CurrentTime TimeProvider = NewMockTimeProvider(time.Now().UTC())

func SetFixedTime(t time.Time) {
	CurrentTime = &MockTimeProvider{FrozenTime: t}
}

func ResetTime() {
	CurrentTime = &MockTimeProvider{FrozenTime: time.Now().UTC()}
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
