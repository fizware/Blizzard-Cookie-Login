package battleinfo

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var verified = false

func GetBattleInfo(cookie string) (*BattleAccount, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://account.battle.net/oauth2/authorization/account-settings", strings.NewReader(""))
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	u, _ := url.Parse("https://us.battle.net")
	jar.SetCookies(u, []*http.Cookie{{
		Name:     "BA-tassadar",
		Value:    cookie,
		Path:     "/login",
		Domain:   ".battle.net",
		SameSite: http.SameSiteNoneMode,
	}})
	client := &http.Client{
		Jar: jar,
	}
	res, err := client.Do(req)
	if err != nil {
		if res != nil {
			location, _ := res.Location()
			if strings.Contains(location.String(), "challenge") && !verified {
				values := url.Values{}
				values.Add("csrf", "")
				values.Add("csrf", "")
				_, err := http.PostForm(location.String(), values)
				if err != nil {
					return nil, err
				}
				verified = true
				return GetBattleInfo(cookie)
			}
		}
		verified = false
		return nil, err
	}
	verified = false
	overviewReq, _ := http.NewRequest(http.MethodGet, "https://account.battle.net/api/overview", strings.NewReader(""))
	res, err = client.Do(overviewReq)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if body == nil || len(body) == 0 {
		return nil, errors.New("invalid cookie")
	}
	var data map[string]any
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	if data["authenticated"] != nil && !data["authenticated"].(bool) {
		return nil, errors.New("invalid cookie")
	}
	accountDetails := data["accountDetails"].(map[string]any)
	accountInfo := new(BattleAccount)
	accountInfo.Email = accountDetails["email"].(string)
	if accountInfo.Email == "" {
		return nil, errors.New("invalid cookie")
	}
	accountInfo.AccountID = int64(accountDetails["accountId"].(float64))
	accountInfo.Cookie = cookie
	return accountInfo, nil
}
