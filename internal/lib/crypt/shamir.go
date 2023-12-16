package crypt

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"math/big"
)

const num256bits = "114905077578207896536214286520987971709694817738092178524314541933911456618019"

type keyShare struct {
	X *big.Int
	Y *big.Int
}

func (sh keyShare) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(sh); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (sh *keyShare) FromBytes(data []byte) error {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(sh)
	if err != nil {
		return err
	}

	return nil
}

func bigIntToString(n *big.Int) string {
	return string(n.Bytes())
}

func generateRandomCoefficients(degree int, prime *big.Int) []*big.Int {
	coefficients := make([]*big.Int, degree)
	for i := range coefficients {
		coefficients[i], _ = rand.Int(rand.Reader, prime)
	}
	return coefficients
}

func createShares(secret *big.Int, n, k int, prime *big.Int) []keyShare {
	coefficients := generateRandomCoefficients(k-1, prime)
	shares := make([]keyShare, n)
	for i := 1; i <= n; i++ {
		x := big.NewInt(int64(i))
		y := new(big.Int).Set(secret)
		for j, coeff := range coefficients {
			exp := new(big.Int).Exp(x, big.NewInt(int64(j+1)), prime)
			term := new(big.Int).Mul(coeff, exp)
			y.Add(y, term)
			y.Mod(y, prime)
		}
		shares[i-1] = keyShare{X: x, Y: y}
	}
	return shares
}

func lagrangeInterpolate(shares []keyShare, prime *big.Int) *big.Int {
	secret := new(big.Int)
	for i, s := range shares {
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)
		for j, otherShare := range shares {
			if i != j {
				numerator.Mul(numerator, otherShare.X)
				denominator.Mul(denominator, new(big.Int).Sub(otherShare.X, s.X))
			}
		}
		denominator.ModInverse(denominator, prime)
		fraction := new(big.Int).Mul(numerator, denominator)
		term := new(big.Int).Mul(s.Y, fraction)
		secret.Add(secret, term)
	}
	secret.Mod(secret, prime)
	return secret
}

func CreateKeyShares(key []byte, threshold, totalShares int) ([][]byte, error) {
	prime, ok := new(big.Int).SetString(num256bits, 10)
	if !ok {
		return nil, errors.New("failed set prime value")
	}

	shares := createShares(new(big.Int).SetBytes(key), totalShares, threshold, prime)

	result := make([][]byte, 0, len(shares))

	for _, i := range shares {
		var bs []byte
		bs, err := i.ToBytes()
		if err != nil {
			return nil, err
		}
		result = append(result, bs)
	}

	return result, nil
}

func RecoveryKeyShares(byteShares [][]byte) ([]byte, error) {
	shares := make([]keyShare, 0)
	prime, ok := new(big.Int).SetString(num256bits, 10)
	if !ok {
		return nil, errors.New("failed set prime value")
	}

	for _, i := range byteShares {
		var sh keyShare
		err := sh.FromBytes(i)
		if err != nil {
			return nil, err
		}
		shares = append(shares, sh)
	}

	recoveredSecret := bigIntToString(lagrangeInterpolate(shares, prime))

	return []byte(recoveredSecret), nil
}
