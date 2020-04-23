package fetchers

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"errors"
)

func TestInt64_Fetch(t *testing.T) {
	envKey := "TEST_VALUE"

	tcs := []struct {
		Name     string
		Default  int64
		Expected int64
		Before   func()
		ErrCheck func(t *testing.T, err error)
	}{
		{
			Name:     "happy path",
			Default:  5,
			Expected: 10,
			Before: func() {
				os.Setenv(envKey, "10")
			},
		},
		{
			Name:     "not set",
			Default:  5,
			Expected: 5,
			ErrCheck: func(t *testing.T, err error) {
				require.True(t, errors.Is(err, ErrNotSet))
			},
		},
		{
			Name:     "parse error",
			Default:  5,
			Expected: 5,
			Before: func() {
				os.Setenv(envKey, "")
			},
			ErrCheck: func(t *testing.T, err error) {
				require.NotNil(t, err)
				require.False(t, errors.Is(err, ErrNotSet))
			},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%v %v", i, tc.Name), func(t *testing.T) {
			defer os.Unsetenv(envKey)

			if tc.Before != nil {
				tc.Before()
			}

			var v int64
			err := NewInt64(&v, tc.Default).Fetch(envKey)

			if tc.ErrCheck != nil {
				tc.ErrCheck(t, err)
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, tc.Expected, v)
		})
	}
}

func TestInt64_Value(t *testing.T) {
	v := int64(55)
	require.Nil(t, Int64{}.Value())
	require.Equal(t, v, Int64{Pointer: &v}.Value())
}
