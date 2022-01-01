package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Responses struct {
		Code []int    `validate:"in:200,404,500"`
		Body []string `validate:"len:5|in:admin,stuff"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "test",
			expectedErr: ErrInterfaceIsNotStruct,
		},
		{
			in: User{
				ID:     "012345678901234567890123456789123456",
				Name:   "Name",
				Age:    20,
				Email:  "user@email.com",
				Role:   "admin",
				Phones: []string{"01234567891", "12345678901"},
				meta:   []byte("extra info"),
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1234",
			},
			expectedErr: ErrLenNotEqual,
		},
		{
			in: User{
				ID:     "1",
				Name:   "Name",
				Age:    20,
				Email:  "user@email.com",
				Role:   "admin",
				Phones: []string{"01234567891", "12345678901"},
				meta:   []byte("extra info"),
			},
			expectedErr: ErrLenNotEqual,
		},
		{
			in: Token{
				Header:    []byte{1, 2, 3},
				Payload:   []byte{4, 5, 6},
				Signature: []byte{7, 8, 9},
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 405,
				Body: "Not found",
			},
			expectedErr: ErrNotInRange,
		},
		{
			in: Responses{
				Code: []int{200, 404, 500},
				Body: []string{"admin", "stuff"},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
				return
			}

			require.ErrorAs(t, err, &tt.expectedErr)
		})
	}
}
