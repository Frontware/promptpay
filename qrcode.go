package PromptPay

import (
	"fmt"
	"regexp"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/snksoft/crc"
)

type Payment struct {
	Amount   float32
	Country  string
	Currency string
	OneTime  bool
	Account  string
	Version  string
}

const ID_PAYLOAD_FORMAT = "00"
const ID_POI_METHOD = "01"
const ID_MERCHANT_INFORMATION_BOT = "29"
const ID_TRANSACTION_CURRENCY = "53"
const ID_TRANSACTION_AMOUNT = "54"
const ID_COUNTRY_CODE = "58"
const ID_CRC = "63"

const PAYLOAD_FORMAT_EMV_QRCPS_MERCHANT_PRESENTED_MODE = "01"
const POI_METHOD_STATIC = "11"
const POI_METHOD_DYNAMIC = "12"
const MERCHANT_INFORMATION_TEMPLATE_ID_GUID = "00"
const BOT_ID_MERCHANT_PHONE_NUMBER = "01"
const BOT_ID_MERCHANT_TAX_ID = "02"
const BOT_ID_MERCHANT_EWALLET_ID = "03"
const GUID_PROMPTPAY = "A000000677010111"
const TRANSACTION_CURRENCY_THB = "764"
const COUNTRY_CODE_TH = "TH"

// NewPayment initialize new payment struct with default values
func NewPayment() (payment Payment) {
	payment = Payment{
		Currency: "THB",
		Country:  COUNTRY_CODE_TH,
		Version:  "000201",
	}
	return
}

func f(id string, value string) string {
	return fmt.Sprintf("%s%02d%s", id, len(value), value)
}

func serialize(xs []string) string {
	return strings.Join(xs, "")
}

func sanitizeTarget(id string) string {
	regex := regexp.MustCompile(`[^0-9]`)
	return regex.ReplaceAllString(id, "")
}

func formatTarget(id string) string {
	numbers := sanitizeTarget(id)
	if len(numbers) >= 13 {
		return numbers
	}
	regex := regexp.MustCompile(`^0`)
	countryCoded := regex.ReplaceAllString(id, "66")
	return fmt.Sprintf("%013s", countryCoded)
}

//
func formatAmount(amount float32) string {
	return fmt.Sprintf("%.2f", amount)
}

func formatCrc(crcValue uint64) string {
	return fmt.Sprintf("%04X", crcValue)
}

// String returns string of Promptpay QRCode
func (p Payment) String() string {
	target := sanitizeTarget(p.Account)
	var targetType string
	switch {
	case len(target) >= 15:
		targetType = BOT_ID_MERCHANT_EWALLET_ID
	case len(target) >= 13:
		targetType = BOT_ID_MERCHANT_TAX_ID
	default:
		targetType = BOT_ID_MERCHANT_PHONE_NUMBER
	}

	var data []string
	data = append(data, f(ID_PAYLOAD_FORMAT, PAYLOAD_FORMAT_EMV_QRCPS_MERCHANT_PRESENTED_MODE))
	if p.Amount != 0 {
		data = append(data, f(ID_POI_METHOD, POI_METHOD_DYNAMIC))
	} else {
		data = append(data, f(ID_POI_METHOD, POI_METHOD_STATIC))
	}
	merchantInfo := serialize([]string{f(MERCHANT_INFORMATION_TEMPLATE_ID_GUID, GUID_PROMPTPAY), f(targetType, formatTarget(target))})
	data = append(data, f(ID_MERCHANT_INFORMATION_BOT, merchantInfo))
	data = append(data, f(ID_COUNTRY_CODE, COUNTRY_CODE_TH))
	data = append(data, f(ID_TRANSACTION_CURRENCY, TRANSACTION_CURRENCY_THB))
	data = append(data, f(ID_PAYLOAD_FORMAT, PAYLOAD_FORMAT_EMV_QRCPS_MERCHANT_PRESENTED_MODE))
	if p.Amount != 0 {
		data = append(data, f(ID_TRANSACTION_AMOUNT, formatAmount(p.Amount)))
	}

	dataToCrc := fmt.Sprintf("%s%s%s", serialize(data), ID_CRC, "04")
	crcValue := crc.CalculateCRC(crc.CCITT, []byte(dataToCrc))
	data = append(data, f(ID_CRC, formatCrc(crcValue)))
	return serialize(data)
}

func (p *Payment) QRCode() (png []byte, err error) {
	png, err = qrcode.Encode(p.String(), qrcode.High, 512)
	return
}
