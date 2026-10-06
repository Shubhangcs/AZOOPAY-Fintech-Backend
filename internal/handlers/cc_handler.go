package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/levionstudio/fintech/internal/models"
	"github.com/levionstudio/fintech/internal/store"
	"github.com/levionstudio/fintech/internal/utils"
)

type CCHandler struct {
	logger  *slog.Logger
	ccstore store.CreditCardPaymentStore
}

func NewCCHandler(logger *slog.Logger, ccstore store.CreditCardPaymentStore) *CCHandler {
	return &CCHandler{
		logger,
		ccstore,
	}
}

func (ch *CCHandler) HandleGetCCOperators(w http.ResponseWriter, r *http.Request) {
	res, err := getCCOperators()
	if err != nil {
		utils.BadRequest(w, ch.logger, "get cc operators api error", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc operators fetched successfully", "operators": res})
}

func getCCOperators() (*models.CardProviderDetailsModel, error) {
	var res models.CardProviderDetailsModel
	if err := utils.GetRequest(
		utils.RechargeKitAPI2+utils.CCOperatorsFetch,
		"Authorization",
		"Bearer "+utils.RechargeKitAPIToken,
		&res,
	); err != nil {
		return nil, err
	}

	if res.Status == 3 {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func (ch *CCHandler) HandleCreateCCBeneficiary(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCreditCardBeneficiaryRequestModel

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ch.logger, "create cc beneficiary", err)
		return
	}

	if err := ch.ccstore.CreateCreditCardBeneficiary(&req); err != nil {
		utils.ServerError(w, ch.logger, "create cc beneficiary", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc beneficiary created successfully"})
}

func (ch *CCHandler) HandleUpdateCCBeneficiary(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateCreditCardBeneficiaryRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ch.logger, "update cc beneficiary", err)
		return
	}

	if err := ch.ccstore.UpdateCreditCardBeneficiary(&req); err != nil {
		utils.ServerError(w, ch.logger, "update cc beneficiary", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc beneficiary updated successfully"})
}

func (ch *CCHandler) HandleDeleteCCBeneficiary(w http.ResponseWriter, r *http.Request) {
	beneId, err := utils.ReadParamIDInt(r)
	if err != nil {
		utils.BadRequest(w, ch.logger, "delete cc beneficiary", err)
		return
	}

	if err := ch.ccstore.DeleteCreditCardBeneficiary(beneId); err != nil {
		utils.ServerError(w, ch.logger, "delete cc beneficiary", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc beneficiary deleted successfully"})
}

func (ch *CCHandler) HandleCreditCardPayment(w http.ResponseWriter, r *http.Request) {
	beneId, err := utils.ReadParamIDInt(r)
	if err != nil {
		utils.BadRequest(w, ch.logger, "credit card payment", err)
		return
	}

	var req models.CreateCreditCardPaymentTransactionRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ch.logger, "credit card payment", err)
		return
	}
	req.PartnerRequestID = uuid.NewString()

	beneDetails, err := ch.ccstore.GetBeneficiaryByBeneficiaryID(beneId)
	if err != nil {
		utils.ServerError(w, ch.logger, "credit card payment", err)
		return
	}
	req.BeneDetails = beneDetails

	transactionId, err := ch.ccstore.InitilizeCreateCreditCardPaymentTransaction(&req)
	if err != nil {
		if isCCClientErr(err) {
			utils.BadRequest(w, ch.logger, "credit card payment", err)
			return
		}
		utils.ServerError(w, ch.logger, "credit card payment", err)
		return
	}

	res, err := creditCardPayment(&req)
	if err != nil && res != nil {
		if err := ch.ccstore.FinalizeCreateCreditCardPaymentTransaction(transactionId, res); err != nil {
			utils.ServerError(w, ch.logger, "credit card payment", err)
			return
		}
		utils.BadRequest(w, ch.logger, "credit card payment", err)
		return
	} else if err != nil && res == nil {
		utils.ServerError(w, ch.logger, "credit card payment", err)
		return
	}

	if err := ch.ccstore.FinalizeCreateCreditCardPaymentTransaction(transactionId, res); err != nil {
		utils.ServerError(w, ch.logger, "credit card payment", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc bill payment successfull", "response": res})
}

func creditCardPayment(req *models.CreateCreditCardPaymentTransactionRequestModel) (*models.CreditCardBillPaymentAPIResponse, error) {
	var res models.CreditCardBillPaymentAPIResponse

	if err := utils.PostRequest(
		utils.RechargeKitAPI2+utils.CCBillPayment,
		"Authorization",
		"Bearer "+utils.RechargeKitAPIToken,
		map[string]any{
			"mobile_no":          req.BeneDetails.PhoneNumber,
			"account_no":         req.BeneDetails.AccountNumber,
			"ifsc":               req.BeneDetails.IFSCCode,
			"bank_name":          req.BeneDetails.BankName,
			"beneficiary_name":   req.BeneDetails.BeneficiaryName,
			"amount":             req.Amount,
			"partner_request_id": req.PartnerRequestID,
			"transfer_type":      "5",
			"operator_code":      req.BeneDetails.OperatorCode,
		},
		&res,
	); err != nil {
		return nil, err
	}
	fmt.Println(res)

	if res.Status == 3 {
		return &res, errors.New(res.Message)
	}

	return &res, nil
}

func (ch *CCHandler) HandleGetCreditCardBeneficiariesByRetailerID(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ch.logger, "get cc beneficiaries by retailer id", err)
		return
	}

	bene, err := ch.ccstore.GetBeneficiariesByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ch.logger, "get cc beneficiaries by retailer id", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "beneficiaries fetched successfully", "beneficiaries": bene})
}

func isCCClientErr(err error) bool {
	msg := err.Error()
	return msg == "retailer not found" ||
		msg == "retailer KYC is not verified" ||
		msg == "retailer is blocked" ||
		msg == "insufficient wallet balance" ||
		msg == "insufficient balance" ||
		msg == "cc transaction not found" ||
		msg == "cc transaction not found or already refunded" ||
		msg == "only FAILED cc transactions can be refunded"
}

func (ch *CCHandler) HandleGetCreditCardPaymentTransactionByID(w http.ResponseWriter, r *http.Request) {
	transactionId, err := utils.ReadParamIDInt(r)
	if err != nil {
		utils.BadRequest(w, ch.logger, "get cc transaction", err)
		return
	}

	res, err := ch.ccstore.GetCreditCardPaymentTransactionByID(transactionId)
	if err != nil {
		if isCCClientErr(err) {
			utils.BadRequest(w, ch.logger, "get cc transaction", err)
			return
		}
		utils.ServerError(w, ch.logger, "get cc transaction", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc transaction fetched successfully", "transaction": res})
}

func (ch *CCHandler) HandleGetAllCreditCardPaymentTransactions(w http.ResponseWriter, r *http.Request) {
	p := utils.ReadQueryParams(r)

	res, err := ch.ccstore.GetAllCreditCardPaymentTransactions(p)
	if err != nil {
		utils.ServerError(w, ch.logger, "get all cc transactions", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc transactions fetched successfully", "transactions": res})
}

func (ch *CCHandler) HandleGetCreditCardPaymentTransactionsByRetailerID(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ch.logger, "get cc transactions by retailer id", err)
		return
	}
	p := utils.ReadQueryParams(r)

	res, err := ch.ccstore.GetCreditCardPaymentTransactionsByRetailerID(retailerId, p)
	if err != nil {
		utils.ServerError(w, ch.logger, "get cc transactions by retailer id", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc transactions fetched successfully", "transactions": res})
}

func (ch *CCHandler) HandleGetCreditCardPaymentTransactionsByDistributorID(w http.ResponseWriter, r *http.Request) {
	distributorId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ch.logger, "get cc transactions by distributor id", err)
		return
	}
	p := utils.ReadQueryParams(r)

	res, err := ch.ccstore.GetCreditCardPaymentTransactionsByDistributorID(distributorId, p)
	if err != nil {
		utils.ServerError(w, ch.logger, "get cc transactions by distributor id", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc transactions fetched successfully", "transactions": res})
}

func (ch *CCHandler) HandleGetCreditCardPaymentTransactionsByMasterDistributorID(w http.ResponseWriter, r *http.Request) {
	mdId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ch.logger, "get cc transactions by md id", err)
		return
	}
	p := utils.ReadQueryParams(r)

	res, err := ch.ccstore.GetCreditCardPaymentTransactionsByMasterDistributorID(mdId, p)
	if err != nil {
		utils.ServerError(w, ch.logger, "get cc transactions by md id", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc transactions fetched successfully", "transactions": res})
}

func (ch *CCHandler) HandleRefundCreditCardPaymentTransaction(w http.ResponseWriter, r *http.Request) {
	transactionId, err := utils.ReadParamIDInt(r)
	if err != nil {
		utils.BadRequest(w, ch.logger, "refund cc transaction", err)
		return
	}

	if err := ch.ccstore.RefundCreditCardPaymentTransaction(transactionId); err != nil {
		if isCCClientErr(err) {
			utils.BadRequest(w, ch.logger, "refund cc transaction", err)
			return
		}
		utils.ServerError(w, ch.logger, "refund cc transaction", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "cc transaction refunded successfully"})
}
