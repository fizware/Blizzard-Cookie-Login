package authenticator

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"net"
	"net/http"
	"strings"
	"syscall"
	"vAuth/proxy"
	"vAuth/proxy/mitm"
	"vAuth/proxy/proxyutil"
)

var (
	wininet               = syscall.MustLoadDLL("Wininet.dll")
	procInternetSetOption = wininet.MustFindProc("InternetSetOptionA")
)

func startProxy() (*proxy.Proxy, error) {
	tlsCert, err := tls.X509KeyPair(battleCert, battleKey)
	if err != nil {
		return nil, ErrStartProxy
	}
	cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, ErrStartProxy
	}
	mitmConfig, err := mitm.NewConfig(cert, tlsCert.PrivateKey.(*rsa.PrivateKey), nil)
	if err != nil {
		return nil, ErrStartProxy
	}
	battleProxy := proxy.NewProxy(proxy.Config{
		ListenAddr: &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 9582,
		},
		MITMConfig: mitmConfig,
		OnRequest:  onRequest,
	})
	if err := battleProxy.Start(); err != nil {
		return nil, ErrStartProxy
	}
	return battleProxy, nil
}

func onRequest(session *proxy.Session) (*http.Request, *http.Response) {
	req := session.Request()
	if !strings.HasPrefix(req.RequestURI, "/login/") || !strings.HasSuffix(req.RequestURI, "/login.app?app=app") {
		return nil, nil
	}
	res := proxyutil.NewResponse(302, strings.NewReader(""), req)
	res.Header.Set("Authentication-State", "DONE")
	res.Header.Set("Location", fmt.Sprint("http://localhost:0/?ST=", battleAccount.Cookie, "&accountId=", battleAccount.AccountID, "&region=", rune(battleAccount.Cookie[0]), rune(battleAccount.Cookie[1]), "&accountName=", battleAccount.Email))
	channel <- nil
	return nil, res
}

func setSystemProxy(server string, enabled bool) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	enable := 1
	if !enabled {
		enable = 0
	}
	err = key.SetDWordValue("ProxyEnable", uint32(enable))
	if err != nil {
		return err
	}
	err = key.SetStringValue("ProxyServer", server)
	if err != nil {
		return err
	}
	_, _, _ = procInternetSetOption.Call(0, 39, 0, 0)
	_, _, _ = procInternetSetOption.Call(0, 37, 0, 0)
	return nil
}
