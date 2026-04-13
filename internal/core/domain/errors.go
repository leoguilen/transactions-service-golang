package domain

import "errors"

var (
	ErrInvalidAccount                  = errors.New("invalid account")
	ErrInvalidAccountID                = errors.New("invalid account ID")
	ErrAccountNotFound                 = errors.New("account not found")
	ErrAccountAlreadyExists            = errors.New("account already exists")
	ErrInvalidAccountDocumentNumber    = errors.New("invalid account document number")
	ErrInvalidTransaction              = errors.New("invalid transaction")
	ErrTransactionAccountNotExists     = errors.New("transaction account does not exist")
	ErrTransactionAccountInvalid       = errors.New("transaction account is invalid")
	ErrTransactionOperationTypeInvalid = errors.New("transaction operation type is invalid")
	ErrTransactionAmountInvalid        = errors.New("transaction amount is invalid")
)
