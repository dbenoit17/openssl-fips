package openssl_test

import (
	"bytes"
	"github.com/golang-fips/openssl-fips/openssl"
	"github.com/golang-fips/openssl-fips/openssl/bbig"
	"math/big"
	"testing"
)

func TestSharedKeyECDH(t *testing.T) {
	if !openssl.Enabled() {
		t.Skip("boringcrypto: skipping test, FIPS not enabled")
	}
	// Test vector from CAVS
	x, ok := new(big.Int).SetString(
		"ead218590119e8876b29146ff89ca61770c4edbbf97d38ce385ed281d8a6b230",
		16,
	)
	if !ok {
		panic("bad hex")
	}
	y, ok := new(big.Int).SetString(
		"28af61281fd35e2fa7002523acc85a429cb06ee6648325389f59edfce1405141",
		16,
	)
	if !ok {
		panic("bad hex")
	}
	k, ok := new(big.Int).SetString(
		"7d7dc5f71eb29ddaf80d6214632eeae03d9058af1fb6d22ed80badb62bc1a534",
		16,
	)
	if !ok {
		panic("bad hex")
	}
	peerPublicKey := []byte{
		// uncompressed
		0x04,
		// X
		0x70, 0x0c, 0x48, 0xf7, 0x7f, 0x56, 0x58, 0x4c,
		0x5c, 0xc6, 0x32, 0xca, 0x65, 0x64, 0x0d, 0xb9,
		0x1b, 0x6b, 0xac, 0xce, 0x3a, 0x4d, 0xf6, 0xb4,
		0x2c, 0xe7, 0xcc, 0x83, 0x88, 0x33, 0xd2, 0x87,
		// Y
		0xdb, 0x71, 0xe5, 0x09, 0xe3, 0xfd, 0x9b, 0x06,
		0x0d, 0xdb, 0x20, 0xba, 0x5c, 0x51, 0xdc, 0xc5,
		0x94, 0x8d, 0x46, 0xfb, 0xf6, 0x40, 0xdf, 0xe0,
		0x44, 0x17, 0x82, 0xca, 0xb8, 0x5f, 0xa4, 0xac,
	}
	expected := []byte{
		0x46, 0xfc, 0x62, 0x10, 0x64, 0x20, 0xff, 0x01,
		0x2e, 0x54, 0xa4, 0x34, 0xfb, 0xdd, 0x2d, 0x25,
		0xcc, 0xc5, 0x85, 0x20, 0x60, 0x56, 0x1e, 0x68,
		0x04, 0x0d, 0xd7, 0x77, 0x89, 0x97, 0xbd, 0x7b,
	}

	bx, by, bk := bbig.Enc(x), bbig.Enc(y), bbig.Enc(k)

	priv, err := openssl.NewPrivateKeyECDH("P-256", bx, by, bk)
	if err != nil {
		t.Log("NewPrivateKeyECDH failed", err)
		t.Fail()
	}
	derived, err := openssl.SharedKeyECDH(priv, peerPublicKey)
	if err != nil {
		t.Log("SharedKeyECDH failed", err)
		t.Fail()
	}
	if !bytes.Equal(derived, expected) {
		t.Log("derived shared secret doesn't match")
		t.Fail()
	}
}
