package PromptPay

import (
	"fmt"
	"regexp"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/snksoft/crc"
)

const (
	idPayloadFormat                            = "00"
	idPOIMethod                                = "01"
	idMerchantInformationBOT                   = "29"
	idTransactionCurrency                      = "53"
	idTransactionAmount                        = "54"
	idCountryCode                              = "58"
	idCRC                                      = "63"
	payloadFormatEMVQRCPSmerchantPresentedMode = "01"
	poiMethodStatic                            = "11"
	poiMathodDynamic                           = "12"
	merchantInformationTemplateIDGUID          = "00"
	botIDMerchantPhoneNumber                   = "01"
	botIDerchantTaxID                          = "02"
	botIDMerchantEwalletID                     = "03"
	guidPromptpay                              = "A000000677010111"
	transactionCurrencyTHB                     = "764"
	countryCodeTH                              = "TH"
)

// Payment is the payment definition
type Payment struct {
	Amount              float32              // Default is 0
	Country             string               // Default is TH
	Currency            string               // Default is THB
	transactionCurrency string               // ISO 4267
	OneTime             bool                 // One time payment type
	Account             string               // Can be tax id, phone number or personal id card
	Version             string               // Default is 000201
	QRCodeQuality       qrcode.RecoveryLevel // Default is HIGH
}

// NewPayment initialize new payment struct with default values for THB payment in Thailand
func NewPayment() (payment Payment) {
	payment = Payment{
		Currency:            "THB",
		Country:             countryCodeTH,
		transactionCurrency: transactionCurrencyTHB,
		Version:             "000201",
		QRCodeQuality:       qrcode.High,
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
		targetType = botIDMerchantEwalletID
	case len(target) >= 13:
		targetType = botIDerchantTaxID
	default:
		targetType = botIDMerchantPhoneNumber
	}

	var data []string
	data = append(data, f(idPayloadFormat, payloadFormatEMVQRCPSmerchantPresentedMode))
	if p.Amount != 0 {
		data = append(data, f(idPOIMethod, poiMathodDynamic))
	} else {
		data = append(data, f(idPOIMethod, poiMethodStatic))
	}
	merchantInfo := serialize([]string{f(merchantInformationTemplateIDGUID, guidPromptpay), f(targetType, formatTarget(target))})
	data = append(data, f(idMerchantInformationBOT, merchantInfo))
	data = append(data, f(idCountryCode, countryCodeTH))
	data = append(data, f(idTransactionCurrency, p.transactionCurrency))
	data = append(data, f(idPayloadFormat, payloadFormatEMVQRCPSmerchantPresentedMode))
	if p.Amount != 0 {
		data = append(data, f(idTransactionAmount, formatAmount(p.Amount)))
	}

	dataToCrc := fmt.Sprintf("%s%s%s", serialize(data), idCRC, "04")
	crcValue := crc.CalculateCRC(crc.CCITT, []byte(dataToCrc))
	data = append(data, f(idCRC, formatCrc(crcValue)))
	return serialize(data)
}

// QRCode returns png as []byte
func (p *Payment) QRCode() (png []byte, err error) {
	png, err = qrcode.Encode(p.String(), p.QRCodeQuality, 512)
	return
}
