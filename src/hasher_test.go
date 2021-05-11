package src

import (
	"fmt"
	"testing"
)

func TestMD5Hasher_Hash(t *testing.T) {
	type testCase struct {
		input  []byte
		opt    HasherOption
		exp    string
		expErr error
	}
	testCases := []testCase{{
		input: []byte("input"),
		opt:   WithMD5(),
		exp:   "a43c1b0aa53a0c908810c06ab1ff3967",
	}, {
		input: []byte("long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"long long long long long long long long long long long long long long long long long long long long long long long " +
			"input"),
		opt: WithMD5(),
		exp: "5d72973efa07bdfb34bc562dbd52c792",
	}, {
		input: []byte("ø˜˜ß©πø˜£ø§∂ƒ©"),
		opt:   WithMD5(),
		exp:   "dda3f39521dab212e5648d3604eb0fdf",
	}, {
		input: []byte{},
		opt:   WithMD5(),
		exp:   "d41d8cd98f00b204e9800998ecf8427e",
	}, {
		input:  []byte("input"),
		opt:    nil,
		expErr: fmt.Errorf("hasher func is not configured, use HasherOption for choose one of them"),
	}}

	for _, tc := range testCases {
		h := NewHasher(tc.opt)
		res, err := h.Hash(tc.input)
		if tc.expErr != nil {
			if err == nil {
				t.Errorf("expected error: %v, got nil", tc.expErr)
				continue
			}
			if err.Error() != tc.expErr.Error() {
				t.Errorf("expected error: %v, got %v", tc.expErr, err)
				continue
			}
			continue
		}
		if tc.exp != res {
			t.Errorf("expected: %#v, got %v", tc.exp, res)
		}
	}
}
