package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestQfq(t *testing.T) {
	tests := map[string]struct {
		val  Qfq
		err1 error
		err2 error
	}{
		"simple": {val: Qfq{Weight: uint32Ptr(2), Lmax: uint32Ptr(4)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalQfq(&testcase.val)
			if err1 != nil {
				if errors.Is(err1, testcase.err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Qfq{}
			err2 := unmarshalQfq(data, &val)
			if err2 != nil {
				if errors.Is(err2, testcase.err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Qfq missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalQfq(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
