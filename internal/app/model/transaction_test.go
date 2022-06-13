package model_test

import (
	"testing"

	"github.com/honyshyota/constanta-rest-api/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		user    func() *model.Transaction
		isValid bool
	}{
		{
			name: "valid",
			user: func() *model.Transaction {
				return model.TestTransaction(t)
			},
			isValid: true,
		},
		{
			name: "empty payment",
			user: func() *model.Transaction {
				u := model.TestTransaction(t)
				u.Pay = 0
				return u
			},
			isValid: false,
		},
		{
			name: "negative payment",
			user: func() *model.Transaction {
				u := model.TestTransaction(t)
				u.Pay = -2
				return u
			},
			isValid: false,
		},
		{
			name: "invalid currency",
			user: func() *model.Transaction {
				u := model.TestTransaction(t)
				u.Currency = "RU"
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.user().Validate())
			} else {
				assert.Error(t, tc.user().Validate())
			}
		})
	}
}
