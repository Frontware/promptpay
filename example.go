package PromptPay

import "fmt"

// Example code for godoc
func Example() {

	myPayment := NewPayment()
	myPayment.Amount = 45.10 // THB
	myPayment.Account = "0105540087061"
	qrcode := myPayment.String()
	fmt.Println("QRCode string ", qrcode)

	// Output:
	// QRCode string
}
