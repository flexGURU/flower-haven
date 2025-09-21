package services

type IPayStack interface {
	InitializePayment(email string, amount int64) (string, string, error)
	VerifyPayment(reference string, amount int64) (string, error)
}
