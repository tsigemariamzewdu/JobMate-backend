package interfaces

// IOTPSender defines SMS sending (service interface)
type IOTPSender interface {
	SendOTP(phone string, code string) error
}
