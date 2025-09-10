package storage

import (
	"encoding/json"

)

type KeyAttestation struct {
	Key string `json:"key"`
}

// WalletInstanceInfo holds information about a walletinstance for storage
type WalletInstanceInfo struct {
    HardwareKeyTag string `json:"hardware_key_tag"`
    // KeyAttestation string `json:"key_attestation"`
    KeyAttestation KeyAttestation
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (info *WalletInstanceInfo) UnmarshalJSON(src []byte) error {
	type walletInstanceInfo WalletInstanceInfo
	ii := walletInstanceInfo(*info)
	if err := json.Unmarshal(src, &ii); err != nil {
		return err
	}
	*info = WalletInstanceInfo(ii)
	return nil
}


// WalletInstanceStorageBackend is an interface to store WalletInstanceInfo
type WalletInstanceStorageBackend interface {
	Write(hardwareKeyTag string, info WalletInstanceInfo) error
	Delete(hardwareKeyTag string) error
	WalletInstance(hardwareKeyTag string) (*WalletInstanceInfo, error)
	Load() error
}
