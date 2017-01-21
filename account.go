package main

import (
	"encoding/json"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"time"
)

const PlanMonthly = "padlock-cloud-monthly"

type Account struct {
	Email    string
	Created  time.Time
	Customer *stripe.Customer
}

func (acc *Account) Subscription() *stripe.Sub {
	if acc.Customer == nil {
		return nil
	}
	subs := acc.Customer.Subs.Values

	if len(subs) == 0 {
		return nil
	}

	return subs[0]
}

// Implements the `Key` method of the `Storable` interface
func (acc *Account) Key() []byte {
	return []byte(acc.Email)
}

// Implementation of the `Storable.Deserialize` method
func (acc *Account) Deserialize(data []byte) error {
	return json.Unmarshal(data, acc)
}

// Implementation of the `Storable.Serialize` method
func (acc *Account) Serialize() ([]byte, error) {
	return json.Marshal(acc)
}

func (acc *Account) CreateCustomer() error {
	params := &stripe.CustomerParams{
		Email: acc.Email,
		Plan:  PlanMonthly,
	}

	var err error
	acc.Customer, err = customer.New(params)

	return err
}

func (acc *Account) UpdateCustomer() error {
	if acc.Customer == nil {
		return nil
	}

	if c, err := customer.Get(acc.Customer.ID, nil); err != nil {
		return err
	} else {
		acc.Customer = c
		return nil
	}
}

func (acc *Account) SetPaymentSource(token string) error {
	params := &stripe.CustomerParams{}
	params.SetSource(token)

	var err error
	acc.Customer, err = customer.Update(acc.Customer.ID, params)
	return err
}

func (acc *Account) HasActiveSubscription() bool {
	sub := acc.Subscription()
	return sub != nil && sub.Status == "active"
}

func (acc *Account) RemainingTrialPeriod() time.Duration {
	sub := acc.Subscription()

	if sub == nil {
		return 0
	}

	trialEnd := time.Unix(sub.TrialEnd, 0)
	remaining := trialEnd.Sub(time.Now())
	if remaining < 0 {
		return 0
	} else {
		return remaining
	}
}

func (acc *Account) RemainingTrialDays() int {
	return int(acc.RemainingTrialPeriod().Hours()/24) + 1
}

func NewAccount(email string) (*Account, error) {
	acc := &Account{
		Email:   email,
		Created: time.Now(),
	}

	return acc, nil
}
