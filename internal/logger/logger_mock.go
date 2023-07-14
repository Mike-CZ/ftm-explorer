package logger

// MockLogger is a mock of ILogger interface.
// All methods will return with no action.
type MockLogger struct{}

// NewMockLogger creates a new mock instance.
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Critical mocks base method.
func (m *MockLogger) Critical(_ ...interface{}) {
	return
}

// Criticalf mocks base method.
func (m *MockLogger) Criticalf(_ string, _ ...interface{}) {
	return
}

// Debug mocks base method.
func (m *MockLogger) Debug(_ ...interface{}) {
	return
}

// Debugf mocks base method.
func (m *MockLogger) Debugf(_ string, _ ...interface{}) {
	return
}

// Error mocks base method.
func (m *MockLogger) Error(_ ...interface{}) {
	return
}

// Errorf mocks base method.
func (m *MockLogger) Errorf(_ string, _ ...interface{}) {
	return
}

// Fatal mocks base method.
func (m *MockLogger) Fatal(_ ...interface{}) {
	return
}

// Fatalf mocks base method.
func (m *MockLogger) Fatalf(_ string, _ ...interface{}) {
	return
}

// Info mocks base method.
func (m *MockLogger) Info(_ ...interface{}) {
	return
}

// Infof mocks base method.
func (m *MockLogger) Infof(_ string, _ ...interface{}) {
	return
}

// ModuleLogger mocks base method.
func (m *MockLogger) ModuleLogger(_ string) ILogger {
	return m
}

// Notice mocks base method.
func (m *MockLogger) Notice(_ ...interface{}) {
	return
}

// Noticef mocks base method.
func (m *MockLogger) Noticef(_ string, _ ...interface{}) {
	return
}

// Panic mocks base method.
func (m *MockLogger) Panic(_ ...interface{}) {
	return
}

// Panicf mocks base method.
func (m *MockLogger) Panicf(_ string, _ ...interface{}) {
	return
}

// Printf mocks base method.
func (m *MockLogger) Printf(_ string, _ ...interface{}) {
	return
}

// Warning mocks base method.
func (m *MockLogger) Warning(_ ...interface{}) {
	return
}

// Warningf mocks base method.
func (m *MockLogger) Warningf(_ string, _ ...interface{}) {
	return
}
