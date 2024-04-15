package security_test

import "github.com/stretchr/testify/mock"

type MockSignature struct {
	mock.Mock
}

func (m *MockSignature) VerifySignature(sig [64]byte, msg [32]byte, pubKey [32]byte) (bool, error) {
	returnValues := m.Called(sig, msg, pubKey)

	return returnValues.Get(0).(bool), returnValues.Error(1)
}
