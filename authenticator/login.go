package authenticator

import (
	"errors"
	"time"
	"vAuth/battleinfo"
)

var (
	ErrStartProxy = errors.New("unable to start proxy")
	channel       = make(chan error)
	battleAccount *battleinfo.BattleAccount
)

func Login(cookie string) error {
	account, err := battleinfo.GetBattleInfo(cookie)
	if err != nil {
		return err
	}
	if err := closeAndFindBattleNet(); err != nil {
		return errors.New("unable to find battle.net, try opening it")
	}
	logoutBattleNet()
	if err := installCert(); err != nil {
		return ErrStartProxy
	}
	if err := setSystemProxy("127.0.0.1:9582", true); err != nil {
		return ErrStartProxy
	}
	battleAccount = account
	battleProxy, err := startProxy()
	defer func() {
		_ = setSystemProxy("", false)
		_ = removeCert()
		if battleProxy != nil {
			battleProxy.Close()
		}
	}()
	if err != nil {
		return ErrStartProxy
	}
	if err := startBattleNet(); err != nil {
		return errors.New("unable to start battle.net")
	}
	err = <-channel
	time.Sleep(1 * time.Second) // give it time to finish the request
	return err
}
