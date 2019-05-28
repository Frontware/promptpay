package promptpay

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Type type of PromptPay
type Type int

const (
	// UNKNOWN unknown PromptPay type
	UNKNOWN Type = 0
	// ID tax/card ID
	ID Type = 1
	// PHONE mobile/phone
	PHONE Type = 2
	// EWALLET E-Wallet
	EWALLET Type = 3

	idPayLoadFormat       = "00"
	idPOIMethod           = "01"
	idMerchantInfo        = "29"
	idCountryCode         = "58"
	idTransactionCurrency = "53"
	idTransactionAmount   = "54"
	idCRC                 = "63"

	payloadFormat = "01"
	poiStatic     = "11"
	poiDynamic    = "12"
	merchantInfo  = "00"
	botIDPhone    = "01"
	botIDTax      = "02"
	botIDEwallet  = "03"
	guidPromptPay = "A000000677010111"
	currencyCode  = "764"
	phoneCode     = "66"
	countryCode   = "TH"

	sizeCRC = 4
)

// PromptPay basic structure accomodating basic information
type PromptPay struct {
	PromptPayID string
	Amount      float64
	OneTime     bool
}

func f(id, value string) string {
	return fmt.Sprintf("%s%02d%s", id, len(value), value)
}

func formatPhone(value string) string {
	return fmt.Sprintf("00%s%s", phoneCode, value)
}

func formatAmount(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func formatCRC(value uint16) string {
	return fmt.Sprintf("%04X", value)
}

// Gen returns a string to create a PromptPay QRCode.
// It uses PromptPay structure to generate a string.
// Two arguments must be provided in order to be able to generate a string.
// PromptPayID and amount must be provided.
func (p *PromptPay) Gen() (string, error) {

	var buffer bytes.Buffer
	var merchant bytes.Buffer

	if strings.TrimSpace(p.PromptPayID) == "" {
		return "", errors.New("empty PromptPayID")
	}

	if p.Amount < 0.0 {
		return "", errors.New("amount can't be negative")
	}

	// QRCode specification
	if _, err := buffer.WriteString(f(idPayLoadFormat, payloadFormat)); err != nil {
		return "", err
	}

	if !p.OneTime {
		if _, err := buffer.WriteString(f(idPOIMethod, poiStatic)); err != nil {
			return "", err
		}
	} else {
		if _, err := buffer.WriteString(f(idPOIMethod, poiDynamic)); err != nil {
			return "", err
		}
	}

	// merchant information
	if _, err := merchant.WriteString(f(merchantInfo, guidPromptPay)); err != nil {
		return "", err
	}

	switch p.GetPromptPayType() {
	case ID:
		if _, err := merchant.WriteString(f(botIDTax, p.PromptPayID)); err != nil {
			return "", err
		}
	case PHONE:
		if _, err := merchant.WriteString(f(botIDPhone, formatPhone(p.PromptPayID))); err != nil {
			return "", err
		}
	case EWALLET:
		if _, err := merchant.WriteString(f(botIDEwallet, p.PromptPayID)); err != nil {
			return "", err
		}
	default:
		return "", errors.New("invalid merchant information")
	}

	// write merchant information to buffer
	if _, err := buffer.WriteString(f(idMerchantInfo, merchant.String())); err != nil {
		return "", err
	}

	// transaction type
	if _, err := buffer.WriteString(f(idCountryCode, countryCode)); err != nil {
		return "", err
	}
	if _, err := buffer.WriteString(f(idTransactionCurrency, currencyCode)); err != nil {
		return "", err
	}

	// transaction amount
	if p.Amount > 0.0 {
		amountStr := formatAmount(p.Amount)
		if _, err := buffer.WriteString(f(idTransactionAmount, amountStr)); err != nil {
			return "", err
		}
	}

	// CRC field
	if _, err := buffer.WriteString(fmt.Sprintf("%s%02d", idCRC, sizeCRC)); err != nil {
		return "", err
	}

	// Generate CRC16 checksum
	crc := Checksum(XModemRev, buffer.Bytes())
	if _, err := buffer.WriteString(formatCRC(crc)); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// GetPromptPayType returns PromptPayType.
// Check whether the PromptPayID is and ID or a phone.
func (p *PromptPay) GetPromptPayType() Type {
	// Tax/Card ID
	if len(p.PromptPayID) == 13 {
		return ID
	}
	// Phone or mobile
	if (len(p.PromptPayID) == 9) && p.PromptPayID[0] != '0' ||
		(len(p.PromptPayID) == 10 && p.PromptPayID[0] == '0' && p.PromptPayID[1] != '0') ||
		(len(p.PromptPayID) == 11 && p.PromptPayID[0:2] == phoneCode && p.PromptPayID[2] != '0') {
		bound := len(p.PromptPayID) - 9
		if bound > 0 {
			p.PromptPayID = p.PromptPayID[bound:]
		}
		return PHONE
	}
	// TODO: implement for E-wallet

	return UNKNOWN
}
