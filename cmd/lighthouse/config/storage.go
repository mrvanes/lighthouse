package config

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/go-oidfed/lighthouse/storage"
)

type storageConf struct {
	DataDir     string      `yaml:"data_dir"`
	BackendType backendType `yaml:"backend"`
}

func (c *storageConf) validate() error {
	if c.DataDir == "" {
		return errors.New("error in storage conf: data_dir must be specified")
	}
	return nil
}

var defaultStorageConf = storageConf{
	BackendType: BackendTypeBadger,
}

type backendType string

const (
	BackendTypeJSON   backendType = "json"
	BackendTypeBadger backendType = "badger"
)

// IsValid checks if the storage type is a known value
func (bt backendType) IsValid() bool {
	switch bt {
	case BackendTypeJSON, BackendTypeBadger:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (bt backendType) String() string {
	return string(bt)
}

// LoadStorageBackends loads and returns the storage backends for the passed Config
func LoadStorageBackends(c storageConf) (
	subordinateStorage storage.SubordinateStorageBackend,
	trustMarkedEntitiesStorage storage.TrustMarkedEntitiesStorageBackend,
	walletInstanceStorage storage.WalletInstanceStorageBackend, err error,
) {
	switch c.BackendType {
	case BackendTypeJSON:
		warehouse := storage.NewFileStorage(c.DataDir)
		subordinateStorage = warehouse.SubordinateStorage()
		trustMarkedEntitiesStorage = warehouse.TrustMarkedEntitiesStorage()
		walletInstanceStorage = warehouse.WalletInstanceStorage()
	case BackendTypeBadger:
		warehouse, err := storage.NewBadgerStorage(c.DataDir)
		if err != nil {
			return nil, nil, nil, err
		}
		subordinateStorage = warehouse.SubordinateStorage()
		trustMarkedEntitiesStorage = warehouse.TrustMarkedEntitiesStorage()
		walletInstanceStorage = nil
	default:
		return nil, nil, nil, errors.Errorf("unknown storage backend type: %s", c.BackendType)
	}
	log.Info("Loaded storage backend")
	return
}
