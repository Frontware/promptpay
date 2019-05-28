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
package main

import (
	"fmt"

	pp "github.com/Frontware/promptpay"
)

func main() {
	payment := pp.PromptPay{
		PromptPayID: "0105540087061", // Tax-ID/ID Card/E-Wallet
		Amount:      100.55,          // Positive amount
	}

	qrcode, _ := payment.Gen() // Generate string to be use in QRCode
	fmt.Println(qrcode)        // Print string
}
```

## Documentation

### API

[![GoDoc](https://godoc.org/github.com/Frontware/promptpay?status.svg)](https://godoc.org/github.com/Frontware/promptpay)

### Specifications

EMV QR Code specification: https://www.emvco.com/wp-content/plugins/pmpro-customizations/oy-getfile.php?u=/wp-content/uploads/documents/EMVCo-Consumer-Presented-QR-Specification-v1.pdf

-----------------------------------------------
<sup>© 2018 Frontware International. All Rights Reserved.</sup>