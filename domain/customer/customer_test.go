package customer

import (
	"errors"
	"testing"
)

func TestCustomer_NewCustomer(t *testing.T) {
	type testCase struct {
		test        string
		name        string
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Empty name validation",
			name:        "",
			expectedErr: ErrInvalidPerson,
		}, {
			test:        "Valid name",
			name:        "jshan",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCustomer(tc.name)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
