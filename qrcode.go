package PromptPay

import (
	"fmt"
	"regexp"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/snksoft/crc"
)

const (
	idPayloadFormat                                  = "00"
	idPOIMethod                                      = "01"
	idMerchantInformationBOT                         = "29"
	ID_TRANSACTION_CURRENCY                          = "53"
	ID_TRANSACTION_AMOUNT                            = "54"
	ID_COUNTRY_CODE                                  = "58"
	ID_CRC                                           = "63"
	PAYLOAD_FORMAT_EMV_QRCPS_MERCHANT_PRESENTED_MODE = "01"
	POI_METHOD_STATIC                                = "11"
	POI_METHOD_DYNAMIC                               = "12"
	MERCHANT_INFORMATION_TEMPLATE_ID_GUID            = "00"
	BOT_ID_MERCHANT_PHONE_NUMBER                     = "01"
	BOT_ID_MERCHANT_TAX_ID                           = "02"
	BOT_ID_MERCHANT_EWALLET_ID                       = "03"
	GUID_PROMPTPAY                                   = "A000000677010111"
	TRANSACTION_CURRENCY_THB                         = "764"
	COUNTRY_CODE_TH                                  = "TH"
)

// Payment is the payment definition
type Payment struct {
	Amount              float32
	Country             string
	Currency            string
	transactionCurrency string // ISO 4267
	OneTime             bool
	Account             string // Can be tax id, phone number or personal id card
	Version             string
}

// NewPayment initialize new payment struct with default values for THB payment in Thailand
func NewPayment() (payment Payment) {
	payment = Payment{
		Currency:            "THB",
		Country:             COUNTRY_CODE_TH,
		transactionCurrency: TRANSACTION_CURRENCY_THB,
		Version:             "000201",
	}
	return
}

var iso4217 map[string]string // https://en.wikipedia.org/wiki/ISO_4217

func init() {
	iso4217 = make(map[string]string)
	iso4217 = map[string]string{
		"THB": "764",
		"EUR": "978",
	}
}

// SetCurrency set currency iso code
func (p *Payment) SetCurrency(currency string) {
	currency = strings.ToUpper(currency)
	p.transactionCurrency = iso4217[currency]
}

func f(id string, value string) string {
	return fmt.Sprintf("%s%02d%s", id, len(value), value)
}

func serialize(xs []string) string {
	return strings.Join(xs, "")
}

// sanitizeTarget cleans the target string
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
	data = append(data, f(idPayloadFormat, PAYLOAD_FORMAT_EMV_QRCPS_MERCHANT_PRESENTED_MODE))
	if p.Amount != 0 {
		data = append(data, f(idPOIMethod, POI_METHOD_DYNAMIC))
	} else {
		data = append(data, f(idPOIMethod, POI_METHOD_STATIC))
	}
	merchantInfo := serialize([]string{f(MERCHANT_INFORMATION_TEMPLATE_ID_GUID, GUID_PROMPTPAY), f(targetType, formatTarget(target))})
	data = append(data, f(idMerchantInformationBOT, merchantInfo))
	data = append(data, f(ID_COUNTRY_CODE, COUNTRY_CODE_TH))
	data = append(data, f(ID_TRANSACTION_CURRENCY, p.transactionCurrency))
	data = append(data, f(idPayloadFormat, PAYLOAD_FORMAT_EMV_QRCPS_MERCHANT_PRESENTED_MODE))
	if p.Amount != 0 {
		data = append(data, f(ID_TRANSACTION_AMOUNT, formatAmount(p.Amount)))
	}

	dataToCrc := fmt.Sprintf("%s%s%s", serialize(data), ID_CRC, "04")
	crcValue := crc.CalculateCRC(crc.CCITT, []byte(dataToCrc))
	data = append(data, f(ID_CRC, formatCrc(crcValue)))
	return serialize(data)
}

// QRCode returns png as []byte
func (p *Payment) QRCode() (png []byte, err error) {
	png, err = qrcode.Encode(p.String(), qrcode.High, 512)
	return
}
