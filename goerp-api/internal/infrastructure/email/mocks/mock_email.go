package mocks

type MockEmailService struct {
	SendCodeFunc func(to, code string) error
}

func (m *MockEmailService) SendCode(to, code string) error {
	return m.SendCodeFunc(to, code)
}
