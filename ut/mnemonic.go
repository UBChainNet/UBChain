package ut

import (
	"encoding/hex"
	"github.com/UBChainNet/UBChain/crypto/bip32"
	"github.com/UBChainNet/UBChain/crypto/bip39"
	"github.com/UBChainNet/UBChain/crypto/ecc/secp256k1"
	"github.com/UBChainNet/UBChain/crypto/seed"
)

func Entropy() (string, error) {
	s, err := seed.GenerateSeed(seed.MinSeedBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(s), nil
}

func Mnemonic(entropyStr string) (string, error) {
	entropy, err := hex.DecodeString(entropyStr)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

func MnemonicToEc(mnemonic string) (*secp256k1.PrivateKey, error) {
	bytes, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	masterKey, err := bip32.NewMasterKey(bytes)
	if err != nil {
		return nil, err
	}
	return secp256k1.ParseStringToPrivate(hex.EncodeToString(masterKey.Key))
}
