package domain

import (
	"time"
)

type OperationType int

const (
	OperationTypePurchase            OperationType = 1
	OperationTypeWithdrawal          OperationType = 2
	OperationTypePurchaseInstallment OperationType = 3
	OperationTypeCreditVoucher       OperationType = 4
)

type Transaction struct {
	ID              int
	AccountID       int
	OperationTypeID OperationType
	Amount          float64
	EventDate       time.Time
}

func NewTransaction(accountID int, operationTypeID int, amount float64) (*Transaction, error) {
	err := validateTransactionParams(operationTypeID, amount)
	if err != nil {
		return nil, err
	}

	transaction := Transaction{
		AccountID:       accountID,
		OperationTypeID: OperationType(operationTypeID),
		Amount:          getSignedAmount(OperationType(operationTypeID), amount),
		EventDate:       time.Now(),
	}

	return &transaction, nil
}

func validateTransactionParams(operationTypeID int, amount float64) error {
	if operationTypeID < int(OperationTypePurchase) || operationTypeID > int(OperationTypeCreditVoucher) {
		return ErrTransactionOperationTypeInvalid
	}
	if amount <= 0 {
		return ErrTransactionAmountInvalid
	}
	return nil
}

func getSignedAmount(operationTypeID OperationType, amount float64) float64 {
	switch operationTypeID {
	case OperationTypePurchase, OperationTypeWithdrawal, OperationTypePurchaseInstallment:
		return -amount
	case OperationTypeCreditVoucher:
		return amount
	default:
		panic("invalid operation type")
	}
}
