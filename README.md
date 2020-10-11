<img src="https://www.frontware.com/images/img/fw-logo.png" alt="Frontware" width="120"/>

# PromptPay

Generate QRCode for Thai PromptPay

![logo](https://vgy.me/i0c6tm.jpg)

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

[Try it on Go playground](https://play.golang.org/p/9EneQrObJ6T)


## Documentation

### API

[![PkgGoDev](https://pkg.go.dev/badge/github.com/frontware/frontware)](https://pkg.go.dev/github.com/frontware/promptpay)

### Specifications

EMV QR Code specification: https://www.emvco.com/wp-content/plugins/pmpro-customizations/oy-getfile.php?u=/wp-content/uploads/documents/EMVCo-Consumer-Presented-QR-Specification-v1.pdf

-----------------------------------------------
<sup>© 2020 Frontware International. All Rights Reserved.</sup>
