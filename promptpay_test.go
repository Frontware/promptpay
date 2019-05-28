package promptpay

import (
	"testing"
)

func TestPromptPay(t *testing.T) {
	paycheck := PromptPay{}
	promptPayStr, err := paycheck.Gen()

	if err == nil {
		t.Error("Expect error: invalid PromptPayID")
	}

	paycheck.PromptPayID = "0105540087061"
	paycheck.Amount = 100.75

	expect := "00020101021129370016A000000677010111021301055400870615802TH53037645406100.756304C86C"

	// Use ID
	promptPayStr, err = paycheck.Gen()
	if err != nil {
		t.Fatal(err.Error())
	}
	if promptPayStr != expect {
		t.Errorf("PromptPay doesn't match.\nGet\t%s\nExpect\t%s\n", promptPayStr, expect)
	}
	t.Log(promptPayStr)

	for i := 0; i < 128; i++ {
		paycheck.Amount = float64(i)
		promptPayStr, err = paycheck.Gen()
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(promptPayStr)
	}

	// check for invalid id
	paycheck.PromptPayID = "012345678912"
	promptPayStr, err = paycheck.Gen()
	if err == nil {
		t.Error("Expect error: invalid ID")
	}

	expect = "00020101021129370016A000000677010111011300668111111115802TH53037645406100.25630415A6"

	// Use phone
	paycheck.PromptPayID = "0811111111"
	paycheck.Amount = 100.25

	promptPayStr, err = paycheck.Gen()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(promptPayStr)
	if promptPayStr != expect {
		t.Errorf("PromptPay doesn't match\nGet\t%s\nExpect\t%s\n", promptPayStr, expect)
	}
	t.Log(promptPayStr)

	for i := 0; i < 128; i++ {
		paycheck.Amount = float64(i)
		promptPayStr, err = paycheck.Gen()
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(promptPayStr)
	}

	// check for invalid phone
	paycheck.PromptPayID = "668123456745"
	promptPayStr, err = paycheck.Gen()
	if err == nil {
		t.Error("Expect error: invalid phone")
	}
}
