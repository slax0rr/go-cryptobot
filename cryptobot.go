package main

import (
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/slax0rr/go-cryptobot/client"
	"github.com/slax0rr/go-cryptobot/irc"
)

type CryptoBot struct {
	irc    irc.IIrc
	client *client.Client

	pairs      []string
	currencies []string
	commands   map[string]func(string, string, []string)
}

func NewCryptoBot(irc irc.IIrc, client *client.Client) *CryptoBot {
	cb := new(CryptoBot)
	cb.irc = irc
	cb.client = client

	cb.currencies = []string{
		"btc",
		"eur",
		"usd",
		"xrp",
		"ltc",
		"eth",
	}
	cb.pairs = []string{
		"btcusd",
		"btceur",
		"xrpusd",
		"xrpeur",
		"xrpbtc",
		"ltcusd",
		"ltceur",
		"ltcbtc",
		"ethusd",
		"etheur",
		"ethbtc",
	}
	cb.commands = map[string]func(string, string, []string){
		"conv": cb.conv,
	}
	return cb
}

func (cb *CryptoBot) conv(message, nick string, args []string) {
	log.WithFields(log.Fields{
		"message": message,
	}).Debug("Received conversion message")

	re := regexp.MustCompile("^(.{3})\\s+(.{3})$")
	m := re.FindStringSubmatch(message)
	if m == nil {
		log.Error("Conversion did not parse properly.")
		return
	}

	curr1 := m[1]
	curr2 := m[2]
	ok := false
	for _, pair := range cb.pairs {
		if pair == curr1+curr2 {
			ok = true
			break
		}
	}
	if ok == false {
		cb.irc.Write(nick + ": Unknown currency pair: " + curr1 + " " + curr2)
		return
	}

	resp := cb.client.GetTicker(curr1, curr2)
	if resp.Err != nil {
		cb.irc.Write(nick + ": " + resp.Err[0])
	}

	msg := nick + ": " + curr1 + " to " + curr2 + ": Last: " + resp.Last +
		" High: " + resp.High +
		" Low: " + resp.Low +
		" Open: " + resp.Open
	cb.irc.Write(msg)
}

func (cb *CryptoBot) evHandler(message, nick string, args []string) {
	log.WithFields(log.Fields{
		"msg": message,
	}).Debug("Message received")

	message = strings.ToLower(message)
	re := regexp.MustCompile("^(.*?)\\s+(.*?)$")
	m := re.FindStringSubmatch(message)
	if m == nil {
		return
	}
	log.WithFields(log.Fields{
		"parsed": m,
	}).Debug("Parsed received message")

	command, ok := cb.commands[m[1]]
	cmdMsg := m[2]
	if !ok {
		for _, curr := range cb.currencies {
			if curr == m[1] {
				command = cb.commands["conv"]
				cmdMsg = m[0]
			}
		}
	}

	log.WithFields(log.Fields{
		"command": cmdMsg,
	}).Debug("Executing command with message")

	command(cmdMsg, nick, args)
}
