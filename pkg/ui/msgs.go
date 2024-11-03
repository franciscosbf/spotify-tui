package ui

import (
	"math/rand/v2"

	"github.com/franciscosbf/spotify-tui/internals/api"
	"github.com/franciscosbf/spotify-tui/internals/auth"
)

type (
	welcomeColorMsg     int
	initMsg             struct{}
	configReadMsg       struct{}
	awaitDotsMsg        int
	expiredTokenMsg     struct{}
	genTokenMsg         struct{}
	failedRegenTokenMsg struct{}
	removeClickMsg      struct{}
	errMsg              error
	ackedAuthMsg        auth.Token
	userInfoMsg         api.UserProfile
	dismissWarnErrMsg   int
)

type newTokenMsg struct {
	token     auth.Token
	refreshed bool
}

type warnErrMsg struct {
	err error
	id  int
}

func newWarnErrMsg(err error) warnErrMsg {
	id := rand.Int()

	return warnErrMsg{err, id}
}

func newNoWarnErrMsg() warnErrMsg {
	return warnErrMsg{nil, -1}
}

func (w warnErrMsg) warn() bool {
	return w.id != -1
}
