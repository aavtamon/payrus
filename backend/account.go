// PizTec Corporation, 2024. All Rights Reserved.

package main

import (
	"strings"
	"sync/atomic"
	"time"
)

type AccountCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DOB      string `json:"dob"`
}

type Payment struct {
	ConversionRate float64 `json:"conversion_rate"`
}

type Card struct {
	Number     string  `json:"number"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Expiration string  `json:"expiration"`
	Balance    float64 `json:"balance"`
}

type Account struct {
	Id          int64               `json:"id"`
	Credentials *AccountCredentials `json:"credentials"`
	Payment     *Payment            `json:"payment"`
	Card        *Card               `json:"card"`
}

var accounts []*Account
var accountId int64 = time.Now().UnixMilli()

func FindAccountByCredentials(credentials *AccountCredentials) *Account {
	account := FindAccountByEmail(credentials.Email)
	if account == nil {
		return nil
	}

	if !verifyPassword(account, credentials.Password) {
		return nil
	}

	return account
}

func FindAccountById(accountId int64) *Account {
	for _, account := range accounts {
		if account.Id == accountId {
			return account
		}
	}

	return nil

}

func FindAccountByEmail(email string) *Account {
	for _, account := range accounts {
		if account.Credentials.Email == email {
			return account
		}
	}

	return nil

}

func CreateAccount(credentials *AccountCredentials) (*Account, error) {
	account := &Account{
		Id:          createUniqueAccountId(),
		Credentials: credentials,
		Payment: &Payment{
			ConversionRate: 95,
		},
	}

	accounts = append(accounts, account)

	return account, nil
}

func (account *Account) CreateCard(cardRequest *Card) error {
	account.Card = &Card{
		Number:     "1111 2222 3333 4444",
		FirstName:  strings.ToUpper(cardRequest.FirstName),
		LastName:   strings.ToUpper(cardRequest.LastName),
		Balance:    0,
		Expiration: "01/27",
	}

	return nil
}

func createUniqueAccountId() int64 {
	return atomic.AddInt64(&accountId, 1)
}

func verifyPassword(account *Account, password string) bool {
	return account.Credentials.Password == password
}
