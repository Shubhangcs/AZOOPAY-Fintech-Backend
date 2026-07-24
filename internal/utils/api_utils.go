package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	RechargeKitAPI1      = os.Getenv("RECHARGE_KIT_API_1")
	RechargeKitAPI2      = os.Getenv("RECHARGE_KIT_API_2")
	RechargeKitVerifyAPI = os.Getenv("RECHARGE_KIT_VERIFY_API")
	RechargeKitAPIToken  = os.Getenv("RECHARGE_KIT_API_TOKEN")
	PayntricAPI          = os.Getenv("PAYNTRIC_API")
	PayntricAPIToken     = os.Getenv("PAYNTRIC_API_TOKEN")
	PayntricUsername     = os.Getenv("PAYNTRIC_API_USERNAME")
)

const (
	PennyDrop                   = "/validation/penny-drop"
	Payout                      = "/rkitpayout/payoutTransfer"
	PayntricPayout              = "/payout/pay"
	PayoutStatus                = "/recharge/statusCheck"
	PayntricPayoutStatus        = "/payout/status"
	MobileRecharge              = "/recharge/prepaid"
	PostpaidMobileRecharge      = "/recharge/postpaid"
	PrepaidPlanFetch            = "/recharge/prepaidPlanFetch"
	PostpaidBillFetch           = "/recharge/postPaidBillFetch"
	DTHRecharge                 = "/recharge/dth"
	ElectricityBill             = "/recharge/billpayment"
	ElectricityBillFetch        = "/recharge/electricityBillFetch"
	BalanceCheck                = "/recharge/balanceCheck"
	PayntricBalanceCheck        = "/merchant/balance"
	AEPSSubMerchantSignup       = "/submerchant/aeps/signup-min-kyc"
	AEPSMerchantEKYCStatusCheck = "/submerchant/aeps/ekyc-status"
	AEPSBiometricKYC            = "/submerchant/aeps/biometric-kyc"
	AEPSOutletLogin             = "/aeps/outlet-login"
	AEPSOutletLoginStatus       = "/aeps/outlet-login-status"
	AEPSMiniStatement           = "/aeps/mini-statement"
	AEPSCashWithdrawal          = "/aeps/cash-withdrawal"
	AEPSCashDeposit             = "/aeps/cash-deposit"
	AEPSBalanceEnquiry          = "/aeps/balance-enquiry"
	AEPSBankList                = "/aeps/banks"
	AEPSGenerateOTP             = "/aeps/transaction-otp"
)

var apiHTTPClient = &http.Client{Timeout: 70 * time.Second}

func PostRequest(url, authKey, authValue string, body map[string]any, res any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("PostRequest marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("PostRequest build: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(authKey, authValue)

	resp, err := apiHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("PostRequest do: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return fmt.Errorf("PostRequest decode (status %d): %w", resp.StatusCode, err)
	}
	return nil
}
func PostRequest2(url, authKey1, authValue1, authKey2, authValue2 string, body map[string]any, res any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("PostRequest marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("PostRequest build: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(authKey1, authValue1)
	req.Header.Set(authKey2, authValue2)

	resp, err := apiHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("PostRequest do: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return fmt.Errorf("PostRequest decode (status %d): %w", resp.StatusCode, err)
	}
	return nil
}

func GetRequest(url, authKey, authValue string, res any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("GetRequest build: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authKey, authValue)

	resp, err := apiHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("GetRequest do: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return fmt.Errorf("GetRequest decode (status %d): %w", resp.StatusCode, err)
	}
	return nil
}

func GetRequest2(url, authKey1, authValue1, authKey2, authValue2 string, res any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("GetRequest build: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authKey1, authValue1)
	req.Header.Set(authKey2, authValue2)

	resp, err := apiHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("GetRequest do: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return fmt.Errorf("GetRequest decode (status %d): %w", resp.StatusCode, err)
	}
	return nil
}
