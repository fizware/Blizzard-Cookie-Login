package main

import (
	"context"
	"strings"
	"vAuth/authenticator"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) CookieLogin(cookie string) string {
	cookie = strings.TrimSpace(cookie)
	if err := authenticator.Login(cookie); err != nil {
		return err.Error()
	}
	return "signed in successfully"
}
