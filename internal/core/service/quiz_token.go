package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/S-nudhana/stray2stay-quiz-backend/internal/core/domain"
)

// tokenPayload is the data embedded in a signed question token. Keeping the
// breed answer here (rather than in server-side memory) lets a stateless
// instance - or any instance behind a load balancer - verify an answer
// without a shared session store.
type tokenPayload struct {
	PetType   string `json:"petType"`
	BreedName string `json:"breedName"`
	ImageURL  string `json:"imageUrl"`
	IssuedAt  int64  `json:"issuedAt"`
}

type tokenSigner struct {
	secret []byte
}

func newTokenSigner(secret string) *tokenSigner {
	return &tokenSigner{secret: []byte(secret)}
}

func (s *tokenSigner) sign(p tokenPayload) (string, error) {
	raw, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	payloadB64 := base64.RawURLEncoding.EncodeToString(raw)
	sig := s.hmac(payloadB64)
	return payloadB64 + "." + sig, nil
}

func (s *tokenSigner) verify(token string) (tokenPayload, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return tokenPayload{}, domain.ErrInvalidToken
	}
	payloadB64, sig := parts[0], parts[1]

	expectedSig := s.hmac(payloadB64)
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return tokenPayload{}, domain.ErrInvalidToken
	}

	raw, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return tokenPayload{}, domain.ErrInvalidToken
	}

	var p tokenPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return tokenPayload{}, domain.ErrInvalidToken
	}

	return p, nil
}

func (s *tokenSigner) hmac(payloadB64 string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(payloadB64))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
