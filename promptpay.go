package promptpay

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strings"
)

// PromptPay basic structure accomodating basic information
// for creating QRCode.
//
// EMVco
//
// Version number
// field 00
// length 02
// data 01
//
// Type of QR Code
// field 01
// length 02
// data 11 -> many use | 12 -> one use
//
// Merchant information
// field 29
// length 37
// Application ID (mean of payment)
//  -sub-field 00
//  -length 16
//  data A000000677010111 -> PromptPay
// Phone number
//  -sub-field 01
//  -length 13
//  -data cut the prefix, phone number is preceded by 00 and country code (66)
// ID Card
//  -sub-field 02
//  -length 13
//
// Country Code
// field 58
// length 02
// data "TH" -> "Thailand"
//
// Transaction currency
// field 53
// length 03
// data 764 "THB" ISO_4217
//
// Transaction amount
// field 54
// length up to 13
//
// CRC
// field 63
// length 04
// checksum using CRC16 includes all fields including CRC field and length "6304"
type PromptPay struct {
	merchant string
	phone    string
	amount   float64
	oneTime  bool
}

// Gen returns a string to create a PromptPay QRCode.
// It uses PromptPay structure to generate a string.
// Two arguments must be provided in order to be able to generate a string.
// Either merchant ID and amount or phone number and amount must be provided.
// If both merchant ID and phone are provided, merchant ID will be used as default
func (p *PromptPay) Gen() (string, error) {

	var buffer bytes.Buffer

	if strings.TrimSpace(p.merchant) == "" && strings.TrimSpace(p.phone) == "" {
		return "", errors.New("both merchant and phone are empty")
	}

	if p.amount < 0.0 {
		return "", errors.New("amount can't be negative")
	}

	// QRCode specifications
	if p.oneTime {
		if _, err := buffer.WriteString("00020101021229370016A000000677010111"); err != nil {
			return "", err
		}
	} else {
		if _, err := buffer.WriteString("00020101021129370016A000000677010111"); err != nil {
			return "", err
		}
	}

	// merchant information
	if len(p.merchant) == 13 {
		if _, err := buffer.WriteString("0213"); err != nil {
			return "", err
		}
		if _, err := buffer.WriteString(p.merchant); err != nil {
			return "", err
		}
	} else if (len(p.phone) == 9 && p.phone[0] != '0') ||
		(len(p.phone) == 10 && p.phone[0] == '0' && p.phone[1] != '0') ||
		(len(p.phone) == 11 && p.phone[0:2] == "66" && p.phone[2] != '0') {
		if _, err := buffer.WriteString("01130066"); err != nil {
			return "", err
		}
		if len(p.phone) == 9 {
			if _, err := buffer.WriteString(p.phone); err != nil {
				return "", err
			}
		} else if len(p.phone) == 10 {
			if _, err := buffer.WriteString(p.phone[1:]); err != nil {
				return "", err
			}
		} else {
			if _, err := buffer.WriteString(p.phone[2:]); err != nil {
				return "", err
			}
		}
	} else {
		return "", errors.New("invalid merchant information")
	}

	// transaction type
	if _, err := buffer.WriteString("5802TH5303764"); err != nil {
		return "", err
	}

	// transaction amount
	if p.amount > 0.0 {
		amountStr := fmt.Sprintf("%.1f", p.amount)
		if math.Mod(p.amount*100.0, 10.0) != 0.0 {
			amountStr = fmt.Sprintf("%.2f", p.amount)
		}
		if _, err := buffer.WriteString("54"); err != nil {
			return "", err
		}
		if _, err := buffer.WriteString(fmt.Sprintf("%02d", len(amountStr))); err != nil {
			return "", err
		}
		if _, err := buffer.WriteString(amountStr); err != nil {
			return "", err
		}
	}

	// Checksum CRC16
	if _, err := buffer.WriteString("6304"); err != nil {
		return "", err
	}
	// Generate hash
	hash := Checksum(XModemRev, buffer.Bytes())
	if _, err := buffer.WriteString(fmt.Sprintf("%X", hash)); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
