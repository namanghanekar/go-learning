package client

type PaymentClient struct {
}

func NewPaymentClient() *PaymentClient {
	return &PaymentClient{}
}

func (p *PaymentClient) ProcessPayment(
	amount float64,
) (bool, error) {

	return true, nil
}
