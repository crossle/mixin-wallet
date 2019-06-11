package mixin

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
)

type MixinKey struct {
	Address  string
	ViewKey  string
	SpendKey string
}

func LocalGenerateKey() (MixinKey, error) {
	seed := make([]byte, 64)
	_, err := rand.Read(seed)
	if err != nil {
		return MixinKey{}, err
	}
	addr := common.NewAddressFromSeed(seed)
	return MixinKey{Address: addr.String(), ViewKey: addr.PrivateViewKey.String(), SpendKey: addr.PrivateSpendKey.String()}, nil
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

func LocalGenerateAddress(viewKey, spendKey string) (*common.Address, error) {
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
	return &addr, nil
}
