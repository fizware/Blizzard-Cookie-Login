package authenticator

import (
	_ "embed"
	"github.com/SpiderOak/wincertstore"
)

//go:embed battlenet.crt
var battleCert []byte

//go:embed battlenet.key
var battleKey []byte

func installCert() error {
	_ = removeCert()
	store, err := wincertstore.OpenSystemStore(wincertstore.Root)
	if err != nil {
		return err
	}
	defer func(store *wincertstore.Store) {
		_ = store.Close()
	}(store)
	err = store.AppendCertsFromPEM(battleCert)
	return err
}

func removeCert() error {
	store, err := wincertstore.OpenSystemStore(wincertstore.Root)
	if err != nil {
		return err
	}
	defer func(store *wincertstore.Store) {
		_ = store.Close()
	}(store)
	err = store.RemoveCertsFromPEM(battleCert)
	return err
}
