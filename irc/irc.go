package irc

import (
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/slax0rr/go-cryptobot/client"
	"github.com/thoj/go-ircevent"
)

type IIrc interface {
	Connect() bool
	Start(func(string, string, []string))
	Write(string)
	WritePriv(string, string)
}

type Irc struct {
	conn    *irc.Connection
	client  *client.Client
	nick    string
	user    string
	server  string
	port    string
	channel string
}

func NewIrc(nick, user, server, channel string, port int, client *client.Client) IIrc {
	if user == "" {
		user = nick
	}

	if string(channel[0]) != "#" {
		channel = "#" + channel
	}

	i := new(Irc)
	i.conn = irc.IRC(nick, user)
	i.client = client
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

// starts the loop and starts listening to commands from the channel
func (i *Irc) Start(evHandler func(string, string, []string)) {
	i.conn.Join(i.channel)
	i.registerEvents(evHandler)
	i.conn.Loop()
}

func (i *Irc) Write(msg string) {
	i.conn.Privmsg(i.channel, msg)
}

func (i *Irc) WritePriv(nick, msg string) {
	i.conn.Privmsg(nick, msg)
}

func (i *Irc) registerEvents(evHandler func(string, string, []string)) {
	i.conn.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {
			re := regexp.MustCompile("^" + i.nick + ".?\\s+(.*?)$")
			m := re.FindStringSubmatch(event.Message())
			if m == nil {
				return
			}
			evHandler(m[1], event.Nick, event.Arguments)
		}(event)
	})
}
