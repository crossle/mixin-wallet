package mixin

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
)

func LocalGenerateKey() (string, string, string, error) {
	seed := make([]byte, 64)
	_, err := rand.Read(seed)
	if err != nil {
		return "", "", "", err
	}
	addr := common.NewAddressFromSeed(seed)
	fmt.Printf("address:\t%s\n", addr.String())
	fmt.Printf("view key:\t%s\n", addr.PrivateViewKey.String())
	fmt.Printf("spend key:\t%s\n", addr.PrivateSpendKey.String())
	return addr.String(), addr.PrivateViewKey.String(), addr.PrivateSpendKey.String(), nil
}

func ParseKeyFromHex(src string) (crypto.Key, error) {
	var key crypto.Key
	data, err := hex.DecodeString(src)
	if err != nil {
		return crypto.Key{}, err
	}
	if len(data) != len(key) {
		return crypto.Key{}, fmt.Errorf("invalid key length %d", len(data))
	}
	copy(key[:], data)
	return key, nil
}

func generateAddress(viewKey, spendKey string) (*common.Address, error) {
	seed := make([]byte, 64)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, err
	}
	addr := common.NewAddressFromSeed(seed)
	if viewKey != "" {
		key, err := hex.DecodeString(viewKey)
		if err != nil {
			return nil, err
		}
		copy(addr.PrivateViewKey[:], key)
		addr.PublicViewKey = addr.PrivateViewKey.Public()
	}

	key, err := hex.DecodeString(spendKey)
	if err != nil {
		return nil, err
	}
	copy(addr.PrivateSpendKey[:], key)
	addr.PublicSpendKey = addr.PrivateSpendKey.Public()
	fmt.Printf("address:\t%s\n", addr.String())
	fmt.Printf("view key:\t%s\n", addr.PrivateViewKey.String())
	fmt.Printf("spend key:\t%s\n", addr.PrivateSpendKey.String())
	return &addr, nil
}
