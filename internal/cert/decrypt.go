package cert

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/pkcs12"
)

// decryptString decrypts one of the base64-encoded strings produced by
// libnativelib.so's custom stream cipher (FUN_00027098).
//
// Decoded layout: [8 bytes IV][4 bytes FNV-1a checksum big-endian][ciphertext]
//
// The cipher's 1000-round S-box mixing overwrites key-derived values with
// IV-derived ones, so the key parameter is effectively irrelevant — any
// non-empty string works.
func decryptString(encoded, key string) (string, error) {
	if len(key) == 0 {
		return "", fmt.Errorf("empty key")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64: %w", err)
	}
	if len(data) < 12 {
		return "", fmt.Errorf("too short: %d bytes", len(data))
	}

	iv := data[0:8]
	wantChecksum := binary.BigEndian.Uint32(data[8:12])
	ciphertext := data[12:]

	h := uint32(0x811c9dc5)
	for _, b := range ciphertext {
		h = (h ^ uint32(b)) * 0x1000193
	}
	if h != wantChecksum {
		return "", fmt.Errorf("FNV checksum mismatch (got %08x want %08x)", h, wantChecksum)
	}

	fnvLow := byte(h)

	kb := []byte(key)
	kl := len(kb)
	var s [256]byte
	for i := range 256 {
		s[i] = iv[i&7] ^ kb[i%kl] ^ fnvLow ^ 0x5a
	}

	prev := uint32(s[0])
	for round := uint32(1); round <= 1000; round++ {
		prev = uint32(iv[(round-1)&7]) ^ round ^ (prev<<1|(prev&0xFF)>>7) ^ prev
		carry := prev
		for j := 1; j < 256; j++ {
			b := uint32(s[j])
			carry = carry ^ round ^ uint32(iv[(round-1+uint32(j))&7]) ^ b ^ (b<<1|(b>>7))
			s[j] = byte(carry)
		}
	}
	s[0] = byte(prev)

	if len(ciphertext) == 0 {
		return "", nil
	}
	plain := make([]byte, len(ciphertext))

	c0 := uint32(ciphertext[0])
	iv0 := uint32(iv[0])
	plain[0] = byte((c0^iv0^1)<<5|(c0^iv0)>>3) ^ s[0]

	for i := 1; i < len(ciphertext); i++ {
		x := uint32(ciphertext[i]^ciphertext[i-1]) ^ uint32(i+1) ^ uint32(iv[i&7])
		plain[i] = byte((x<<5)|(x>>3)) ^ s[i&0xFF]
	}

	for len(plain) > 0 && plain[len(plain)-1] == 0 {
		plain = plain[:len(plain)-1]
	}
	return string(plain), nil
}

// encryptedStrings are all base64 strings from libnativelib.so's FUN_000274f4.
// One of them decrypts to the alice.p12 password.
var encryptedStrings = []string{
	"eifddmdX43fNE0EVvAB0rGqlJwEpaNNc5CewrEhfsVgeS0CfbIxNSYDIpvo9wqSt41W2oA==",
	"f3mviQpI8d+9AFzoSe5yWZZhV+1zS+2OYmXpXW6R1ATwvGNoZvyoMD9JJYkoTOoqbxPtDw==",
	"EBGjspzB2dbinvUVNvoOnt1vGOqANoq/D1mlGP4hV2R/Z1Bhau27xUQb5lqMq4cOR9X6gQ==",
	"kSpTSBKiPGHYGVSsZMdeFZwvV4qSkvaHAezjgM1mVVc+DoYC+wMUBFtHUTME/lctvSZWCA==",
	"pIWH2NT68cThX5Pi3E1I71nSLxk9be54H4ENogTNEU3A8ORj1yVUOb1SPom979CcoKCV+w==",
	"QfR/UpgYuqCOtsyHL9YkX9TLEig39cvoqGDdf3mI0OKJFTxfoZvWlA00kDNP7y1n9cITqg==",
	"Hl2IGsuYmC3YanEURDbEzTKpjTLY9M1cdYaJNXGZCaKvxuqVr8OaGyKDjEgM5Kwn",
	"EzEWmVPxIClxmrw910OqKL49KihCFNBKM/oyIP+JoOj89Ui6sl7qUBm6yNIF6Rhg",
	"KLs6XVjpDn61H6cs8e9dss2R4Kr7kTcrk4A1OJBe",
	"AW9n9UJgevqQKrOtfUV5uxdTmq/xXS3TTs5dxIgRtOXx/+YbLEO5p/g5AkqFnOiz/iJb/g==",
	"SOnqE/D+v/JHltS+I49tZ2hzNEssVEUt+xq5wmD9pAUb+l3qLlZkTUN54mJDxhw08HOtAw==",
	"kW7WMyMTVL1U4MlzesEiP+8bcyaLp1PQx7EUskLZavXn4SOkOQVb6MexpjA=",
	"V5DnPprXwgQk6fP14EVmUECoiPVjzWfPkthfgqAd9MJwuIjdQ0/1RRjQ7XA=",
	"a1IptAW7x+081RPxToDQn7c28aB7CxY5RoY88nao6B3n7AvoO8l+09QMvlBmuMinHbzRiFuJBuO8XGxwzOhaPeduqUI7e+x5nHx8EA==",
	"Ls/CAqRwvStPokd8o0FTZiKm7e2P72g37ZvIiEOzGUQIhM8/hb0S5VUxEjojmWukUh73dyelQhVNefBiY8kLxPCk/1edN3IlbQ==",
	"3UaW8cbQliNYzy4LZsdO6VocJJELnIYbrkGSMG59zGvAllTZ6fT2KxTjgiBO9/RLqFxM2SsuPNn2STC6xMfsexpsphOxpCTB5usY2g==",
	"PqU/ShbwBjLAeiiguKtKE29cWzJJVzI9dV4aUw==",
	"H1JG9Jm8ZTurGOBYPm038CMYesVBRWVUPDyzMqXWLclr8NIP4qxlRh42qnil7Y57mHowNBhVFn0=",
	"houQvL2SQOVuiQVlE8bNdWB+wwY=",
	"vWrOPxfXl/mPQ8N7VYvvyMNngG5QUD7kMQu6gP0Y5YKZ7AMeMPGXTQ==",
	"8lHPu45LMHfDBNhHi41XnV9clRVn4g/PHeWdMxlfHM8mty6ntMie1B5dz8tKRId1B3xMFQ==",
	"CKq0vH+aIH890hDIf637g4ac5aXy1Wqkd3Erb3IjdNbqE+LrrbqfacudJehNRKrz71VvxA==",
}

// findP12Password tries every known encrypted string against p12Data and
// returns the first plaintext that successfully decodes the certificate.
func findP12Password(p12Data []byte) (string, error) {
	for _, enc := range encryptedStrings {
		plain, err := decryptString(enc, "x")
		if err != nil {
			continue
		}
		blocks, err := pkcs12.ToPEM(p12Data, plain)
		if err == nil && len(blocks) > 0 {
			return plain, nil
		}
	}
	return "", fmt.Errorf("alice.p12: no known password matched")
}
