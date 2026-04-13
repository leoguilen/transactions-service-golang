package domain

import (
	"time"
)

type Account struct {
	ID             int
	DocumentNumber string
	CreatedAt      time.Time
}

func NewAccount(documentNumber string) (*Account, error) {
	err := validateAccountParams(documentNumber)
	if err != nil {
		return nil, err
	}

	account := Account{
		DocumentNumber: documentNumber,
	}

	return &account, nil
}

func validateAccountParams(documentNumber string) error {
	documentNumberLen := len(documentNumber)
	if documentNumberLen != 11 && documentNumberLen != 14 {
		return ErrInvalidAccountDocumentNumber
	}
	return nil
}
