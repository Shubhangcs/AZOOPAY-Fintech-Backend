package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/levionstudio/fintech/internal/models"
	"github.com/levionstudio/fintech/internal/utils"
)

// FinZenG serves Vendor Payment (api_provider "DVSID", routes /payout/create/dev
// and /payout/status-check/dev).

var finzengHTTPClient = &http.Client{Timeout: 70 * time.Second}

var istLocation = time.FixedZone("IST", 5*60*60+30*60)

var finzengDisallowedChars = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// finzengClean drops special characters and collapses repeated spaces, as the
// FinZenG docs ask for beneficiary details (spaces between words are kept).
func finzengClean(s string) string {
	return strings.Join(strings.Fields(finzengDisallowedChars.ReplaceAllString(s, " ")), " ")
}

func finzengConfigured() bool {
	return utils.FinzengAPI != "" && utils.FinzengClientID != "" && utils.FinzengAPIKey != "" &&
		utils.FinzengClientSecret != "" && utils.FinzengServiceSecret != "" && utils.FinzengFromAccount != ""
}

func setFinzengHeaders(req *http.Request) {
	req.Header.Set("client-id", utils.FinzengClientID)
	req.Header.Set("fps3-api-key", utils.FinzengAPIKey)
	req.Header.Set("client-secret", utils.FinzengClientSecret)
	req.Header.Set("service-secret", utils.FinzengServiceSecret)
	req.Header.Set("accept", "application/json")
}

// finzengCode trims the responseCode; FinZenG has sent values like " S00000P".
func finzengCode(res *models.FinzengPayoutAPIResponseModel) string {
	return strings.ToUpper(strings.TrimSpace(res.ResponseCode))
}

// finzengBankRef returns the bank reference only when it is a real value.
func finzengBankRef(res *models.FinzengPayoutAPIResponseModel) string {
	ref := strings.TrimSpace(res.BankReferenceNumber)
	if strings.Contains(strings.ToLower(ref), "not available") {
		return ""
	}
	return ref
}

// callFinzengPayoutAPI initiates a payout. Outcomes:
//   - SUCCESS                 -> SUCCESS
//   - PENDING                 -> PENDING
//   - FAILURE S00000F         -> FAILED (payout rejected)
//   - FAILURE E001..E009      -> FAILED (auth/IP/tariff/KYC - not processed)
//   - FAILURE SN00001         -> FAILED (request validation - not processed)
//   - FAILURE E010 (duplicate reference_id), UNKNOWN, unreadable reply or a
//     network error after sending -> PENDING, since money may have moved;
//     resolve via status check before any refund.
func callFinzengPayoutAPI(logger *slog.Logger, pt *models.PayoutTransactionModel) (resp *models.FinzengPayoutAPIResponseModel, finalStatus, orderID, operatorTxnID string) {
	if !finzengConfigured() {
		logger.Error("finzeng payout api not configured", "payout_transaction_id", pt.PayoutTransactionID)
		return nil, "FAILED", "", ""
	}

	body := map[string]any{
		"reference_id":       pt.PartnerRequestID,
		"purpose_message":    finzengClean("Payout to " + pt.BeneficiaryName),
		"from_account":       utils.FinzengFromAccount,
		"to_account":         strings.TrimSpace(pt.AccountNumber),
		"transfer_type":      strings.ToUpper(strings.TrimSpace(pt.TransferType)),
		"sender_mobile":      strings.TrimSpace(pt.RetailerPhone),
		"transfer_amount":    fmt.Sprintf("%.2f", pt.Amount),
		"beneficiary_name":   finzengClean(pt.BeneficiaryName),
		"beneficiary_mobile": strings.TrimSpace(pt.MobileNumber),
		"ifsc":               strings.ToUpper(strings.TrimSpace(pt.IFSCCode)),
		"beneficiary_bank":   finzengClean(pt.BankName),
		"currency_code":      "INR",
	}

	b, err := json.Marshal(body)
	if err != nil {
		logger.Error("finzeng payout marshal failed", "error", err, "payout_transaction_id", pt.PayoutTransactionID)
		return nil, "FAILED", "", ""
	}

	req, err := http.NewRequest(http.MethodPost, utils.FinzengAPI+utils.FinzengPayoutInitiate, bytes.NewReader(b))
	if err != nil {
		logger.Error("finzeng payout request build failed", "error", err, "payout_transaction_id", pt.PayoutTransactionID)
		return nil, "FAILED", "", ""
	}
	setFinzengHeaders(req)
	req.Header.Set("content-type", "application/json")

	httpRes, err := finzengHTTPClient.Do(req)
	if err != nil {
		logger.Error("finzeng payout api call failed", "error", err, "payout_transaction_id", pt.PayoutTransactionID)
		return nil, "PENDING", "", ""
	}
	defer httpRes.Body.Close()

	var apiResp models.FinzengPayoutAPIResponseModel
	if err := json.NewDecoder(httpRes.Body).Decode(&apiResp); err != nil {
		logger.Error("finzeng payout decode failed", "error", err, "http_status", httpRes.StatusCode, "payout_transaction_id", pt.PayoutTransactionID)
		return nil, "PENDING", "", ""
	}

	resp = &apiResp
	orderID = string(apiResp.SanTxnID)
	operatorTxnID = finzengBankRef(&apiResp)
	code := finzengCode(&apiResp)

	switch strings.ToUpper(strings.TrimSpace(apiResp.Status)) {
	case "SUCCESS":
		finalStatus = "SUCCESS"
	case "PENDING":
		finalStatus = "PENDING"
	case "FAILURE", "FAILED":
		switch {
		case code == "E010":
			finalStatus = "PENDING"
		case code == "S00000F" || code == "SN00001" || strings.HasPrefix(code, "E"):
			finalStatus = "FAILED"
		default:
			finalStatus = "PENDING"
		}
	default:
		finalStatus = "PENDING"
	}

	if finalStatus != "SUCCESS" {
		logger.Warn("finzeng payout not successful", "status", apiResp.Status, "response_code", code, "msg", apiResp.Message, "errors", apiResp.Errors, "payout_transaction_id", pt.PayoutTransactionID)
	}
	return
}

// callFinzengPayoutStatusAPI checks a payout (multipart form: sanTxnId, cTxnId,
// txnDate). Only SUCCESS and FAILURE/S00000F are final; SN00001 (validation),
// SN00002 (not found), UNKNOWN and errors leave the payout PENDING. The stored
// order / bank reference are kept when the reply does not carry them.
func callFinzengPayoutStatusAPI(logger *slog.Logger, pt *models.PayoutTransactionModel) (resp *models.FinzengPayoutAPIResponseModel, finalStatus, orderID, operatorTxnID string) {
	finalStatus = "PENDING"
	orderID = pt.OrderID
	operatorTxnID = pt.OperatorTransactionID

	if !finzengConfigured() {
		logger.Error("finzeng payout api not configured", "payout_transaction_id", pt.PayoutTransactionID)
		return
	}

	var form bytes.Buffer
	mw := multipart.NewWriter(&form)
	_ = mw.WriteField("sanTxnId", pt.OrderID)
	_ = mw.WriteField("cTxnId", pt.PartnerRequestID)
	_ = mw.WriteField("txnDate", pt.CreatedAT.In(istLocation).Format("2006-01-02"))
	if err := mw.Close(); err != nil {
		logger.Error("finzeng status form build failed", "error", err, "payout_transaction_id", pt.PayoutTransactionID)
		return
	}

	req, err := http.NewRequest(http.MethodPost, utils.FinzengAPI+utils.FinzengPayoutStatus, &form)
	if err != nil {
		logger.Error("finzeng status request build failed", "error", err, "payout_transaction_id", pt.PayoutTransactionID)
		return
	}
	setFinzengHeaders(req)
	req.Header.Set("content-type", mw.FormDataContentType())

	httpRes, err := finzengHTTPClient.Do(req)
	if err != nil {
		logger.Error("finzeng status api call failed", "error", err, "payout_transaction_id", pt.PayoutTransactionID)
		return
	}
	defer httpRes.Body.Close()

	var apiResp models.FinzengPayoutAPIResponseModel
	if err := json.NewDecoder(httpRes.Body).Decode(&apiResp); err != nil {
		logger.Error("finzeng status decode failed", "error", err, "http_status", httpRes.StatusCode, "payout_transaction_id", pt.PayoutTransactionID)
		return
	}

	resp = &apiResp
	if id := string(apiResp.SanTxnID); id != "" {
		orderID = id
	}
	if ref := finzengBankRef(&apiResp); ref != "" {
		operatorTxnID = ref
	}

	code := finzengCode(&apiResp)
	switch strings.ToUpper(strings.TrimSpace(apiResp.Status)) {
	case "SUCCESS":
		finalStatus = "SUCCESS"
	case "FAILURE", "FAILED":
		if code == "S00000F" {
			finalStatus = "FAILED"
		} else {
			logger.Warn("finzeng status not final", "response_code", code, "msg", apiResp.Message, "errors", apiResp.Errors, "payout_transaction_id", pt.PayoutTransactionID)
		}
	}
	return
}

func callFinzengBalanceAPI() (*models.FinzengBalanceAPIResponseModel, error) {
	if !finzengConfigured() {
		return nil, fmt.Errorf("finzeng payout api not configured")
	}

	req, err := http.NewRequest(http.MethodPost, utils.FinzengAPI+utils.FinzengBalance, nil)
	if err != nil {
		return nil, err
	}
	setFinzengHeaders(req)
	req.Header.Set("content-type", "application/json")

	httpRes, err := finzengHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	var res models.FinzengBalanceAPIResponseModel
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("finzeng balance decode (status %d): %w", httpRes.StatusCode, err)
	}
	return &res, nil
}
