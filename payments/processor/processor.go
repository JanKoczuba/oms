package processor

import pb "github.com/JanKoczuba/commons/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
