package irc

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/thoj/go-ircevent"
)

type IIrc interface {
	Connect() bool
}

type Irc struct {
	conn    *irc.Connection
	nick    string
	user    string
	server  string
	port    string
	channel string
}

func NewIrc(nick, user, server, channel string, port int) IIrc {
	if user == "" {
		user = nick
	}

	i := new(Irc)
	i.conn = irc.IRC(nick, user)
	i.nick = nick
	i.user = user
	i.server = server
	i.channel = channel
	i.port = strconv.Itoa(port)
	return i
}

func (i *Irc) Connect() bool {
	err := i.conn.Connect(i.server + ":" + i.port)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error occured while trying to connect IRC server")
		return false
	}
	return true
}
