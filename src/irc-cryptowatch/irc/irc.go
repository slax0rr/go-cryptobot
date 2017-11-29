package irc

import (
	"irc-cryptowatch/client"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/thoj/go-ircevent"
)

type IIrc interface {
	Connect() bool
	Start()
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
func (i *Irc) Start() {
	i.conn.Join(i.channel)
	i.registerEvents()
	i.conn.Loop()
}

func (i *Irc) registerEvents() {
	i.conn.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {
			log.WithFields(log.Fields{
				"msg": event.Message(),
			}).Debug("Message received")

			// todo(slax0rr): move this to separate package
			re := regexp.MustCompile("^" + i.nick + ".?\\s(.*?) (.*?)$")
			m := re.FindStringSubmatch(event.Message())
			if m == nil {
				return
			}
			log.WithFields(log.Fields{
				"parsed": m,
			}).Debug("Parsed received message")

			resp := i.client.GetTicker(m[1], m[2])
			if resp.Err != nil {
				i.conn.Privmsg(i.channel, resp.Err[0])
			}

			crypto := strings.ToUpper(m[1])
			fiat := strings.ToUpper(m[2])
			msg := crypto + " to " + fiat + ": Last: " + resp.Last +
				" High: " + resp.High +
				" Low: " + resp.Low +
				" Open: " + resp.Open
			i.conn.Privmsg(i.channel, msg)
		}(event)
	})
}
