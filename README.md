![logo](https://vgy.me/i0c6tm.jpg)

# PromptPay

Generate QRCode for Thai PromptPay

Golang API to generate QRCode

# Build on your computer

## Install dependencies

> make init

## Build

> make build

## Example

```golang

import pp "github.com/Frontware/promptpay"
// Example code for godoc
func main() {

	myPayment := pp.NewPayment()
	myPayment.Amount = 45.10 // THB
	myPayment.Account = "0105540087061"
	qrcode := myPayment.String()
	fmt.Println("QRCode string ", qrcode)

	// Output:
	// QRCode string
}


```

## Documentation

### API

https://godoc.org/github.com/Frontware/promptpay

### Specifications

EMV QR Code specification: https://www.emvco.com/wp-content/plugins/pmpro-customizations/oy-getfile.php?u=/wp-content/uploads/documents/EMVCo-Consumer-Presented-QR-Specification-v1.pdf

-----------------------------------------------
<sup>© 2018 Frontware International. All Rights Reserved.</sup>