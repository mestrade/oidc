// +build go1.12
// +build dontbuildme

package oidc

import (
	"encoding/base64"
	//"encoding/json"
	//"fmt"
	//"strings"
)

// alg
// RS256 == RSASSA-PKCS1-v1_5 + SHA256
// RS384 == RSASSA-PKCS1-v1_5 + SHA384
// RS512 == RSASSA-PKCS1-v1_5 + SHA512
// ES256 == ECDSA P-256 + SHA256
// ES384 == ECDSA P-384 + SHA384
// ES512 == ECDSA P-521 + SHA512
// PS256
// PS384
// PS512
// none => EXIT |
//
// unsecured JWS == len sig == 0
type IdTokenHeader struct {
	// raw: base64
	Raw []byte

	// header
	Kid string `json:"kid"`
	Alg string `json:"alg"`
}

func (h IdTokenHeader) String() string {
	return fmt.Sprintf("HEADER:\n\talg: '%s'\n\tkid: '%s'\n", h.Alg, h.Kid)
}

// Claims
type IdTokenClaims struct {
	// raw: base64
	Raw []byte

	// claims
	Sub           string `json:"sub"`
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	Exp           int    `json:"exp"`
	Iat           int    `json:"iat"`
	Email         string `json:"email"`
	Nonce         string `json:"nonce"`
	EmailVerified bool   `json:"email_verified"` // Addition to provide some additionnal "security" and avoid abuse of oauth for login (optional)
	Azp           string `json:"azp"`            // Addition TBD
}

func (h IdTokenClaims) String() string {
	return fmt.Sprintf("CLAIMS:\n\tsub: '%s'\n\tiss: '%s'\n\taud: '%s'\n\texp: '%d'\n\tiat: '%d'\n\temail: '%s'\n\tnonce: '%s'\n",
		h.Sub,
		h.Iss,
		h.Aud,
		h.Exp,
		h.Iat,
		h.Email,
		h.Nonce)
}

// signature
type IdTokenSignature struct {
	// raw: base64
	Raw []byte

	// debase64 signature
	Blob []byte
}

func (s IdTokenSignature) String() string {
	return fmt.Sprintf("SIG: '%s'\n", s.Raw)
}

type IdToken struct {
	Hdr    IdTokenHeader    // Idtoken Header
	Claims IdTokenClaims    // Idtoken Claims
	Sig    IdTokenSignature // Idtoken signature
	Raw    string           // the raw token..
}

func (idt *IdToken) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n",
		idt.Hdr,
		idt.Claims,
		idt.Sig)
}

// FieldFunc() or Split()
// XXX TODO should be renamed to parseSafeIdToken
func newIdToken(idtoken string) (*IdToken, error) {
	var hdr IdTokenHeader
	var claims IdTokenClaims
	var sig IdTokenSignature

	//fmt.Printf("NEW ID TOKEN!!\n")

	//tok := strings.SplitN(idtoken, ".", 3)
	tok := strings.Split(idtoken, ".")

	if len(tok) != 3 {
		//return nil, errors.New("invalid token for us")
		return nil, ErrParse
	}

	// no signature, NOPE.. invalid.
	if len(tok[0]) == 0 || len(tok[1]) == 0 || len(tok[2]) == 0 {
		//return nil, errors.New("invalid token for us")
		return nil, ErrParse
	}

	//
	// header
	//
	hdrJson, err := base64.RawURLEncoding.DecodeString(tok[0])
	if err != nil {
		return nil, err
	}
	// unmarshal header
	err = json.Unmarshal(hdrJson, &hdr)
	if err != nil {
		return nil, err
	}
	hdr.Raw = []byte(tok[0])
	//fmt.Printf("HEADER: %v\n", h.String())

	//
	// claims
	//
	claimsJson, err := base64.RawURLEncoding.DecodeString(tok[1])
	if err != nil {
		return nil, err
	}

	// unmarshal claims
	err = json.Unmarshal(claimsJson, &claims)
	if err != nil {
		return nil, err
	}
	claims.Raw = []byte(tok[1])

	//
	// signature
	//
	sigBin, err := base64.RawURLEncoding.DecodeString(tok[2])
	if err != nil {
		return nil, err
	}
	sig.Blob = sigBin
	sig.Raw = []byte(tok[2])

	it := IdToken{
		Hdr:    hdr,
		Claims: claims,
		Sig:    sig,
	}
	return &it, nil
}
