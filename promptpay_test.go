package promptpay

import (
	"testing"
)

func TestPromptPay(t *testing.T) {
	paycheck := PromptPay{
		merchant: "0105540087061",
		amount:   100.75,
	}

	// Use ID
	promptPayStr, err := paycheck.Gen()
	if err != nil {
		t.Fatal(err.Error())
	}
	if promptPayStr != "00020101021129370016A000000677010111021301055400870615802TH53037645406100.756304C86C" {
		t.Fatal("PromptPay doesn't match")
	}
	t.Log(promptPayStr)

	for i := 0; i < 128; i++ {
		paycheck.amount = float64(i)
		promptPayStr, err = paycheck.Gen()
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(promptPayStr)
	}

	// Use phone
	paycheck.merchant = ""
	paycheck.phone = "0811111111"
	paycheck.amount = 100.25

	promptPayStr, err = paycheck.Gen()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(promptPayStr)
	if promptPayStr != "00020101021129370016A000000677010111011300668111111115802TH53037645406100.25630415A6" {
		t.Fatal("PromptPay doesn't match")
	}
	t.Log(promptPayStr)

	for i := 0; i < 128; i++ {
		paycheck.amount = float64(i)
		promptPayStr, err = paycheck.Gen()
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(promptPayStr)
	}
}
