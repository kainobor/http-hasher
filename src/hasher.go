package src

import (
	"crypto/md5"
	"fmt"
	"hash"
)

type (
	Hasher interface {
		Hash([]byte) (string, error)
	}

	hasher struct {
		newHash NewHash
	}

	HasherOption func(*hasher)

	NewHash func() hash.Hash
)

var _ Hasher = &hasher{}

func NewHasher(opts ...HasherOption) *hasher {
	h := &hasher{}
	for _, opt := range opts {
		if opt != nil {
			opt(h)
		}
	}

	return h
}

func WithMD5() HasherOption {
	return func(h *hasher) {
		h.newHash = md5.New
	}
}

func (h *hasher) Hash(input []byte) (string, error) {
	if h.newHash == nil {
		return "", fmt.Errorf("hasher func is not configured, use HasherOption for choose one of them")
	}

	internalHasher := h.newHash()
	if _, err := internalHasher.Write(input); err != nil {
		return "", fmt.Errorf("write input: %w", err)
	}

	hashSum := internalHasher.Sum(nil)

	return fmt.Sprintf("%x", hashSum), nil
}
