package card

import (
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/paymentintent"
	//"golang.org/x/text/currency"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {

	stripe.Key = c.Secret
	//create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}
	//params.AddMetadata("key","value")

	//here is the reason for this function
	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}
	return pi, "", err
}

//take stripe error if any, and convert to string
func cardErrorMessage(code stripe.ErrorCode) string {
	var msg string = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card was deeclined"
	case stripe.ErrorCodeExpiredCard:
		msg = "Your Card is Expire"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Incorrect CVC"
	case stripe.ErrorCodeIncorrectZip:
		msg = "Incorrect ZIP/POSTAL Code"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "The Amount is Too Large to charge your card"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "The Amount is Too Small to charge your card"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "Insufficient Balance"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Your Postal Code is invalid"

	default:
		msg = "Your card was deeclined"

	}
	return msg
}
