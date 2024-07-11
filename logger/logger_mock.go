package logger

import (
	"github.com/stretchr/testify/mock"
)

// MockLogger implementation
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Info(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, err error, args ...any) {
	m.Called(msg, err, args)
}

func (m *MockLogger) WithOperation(operationID string) LoggerInterface {
	args := m.Called(operationID)
	return args.Get(0).(LoggerInterface)
}
