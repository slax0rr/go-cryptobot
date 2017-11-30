package main

import (
	"irc-cryptowatch/client"
	"irc-cryptowatch/irc"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
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
	re := regexp.MustCompile("^(.*?)\\s+(.*?)$")
	m := re.FindStringSubmatch(message)
	if m == nil {
		log.Error("Conversion did not parse properly.")
		return
	}

	curr1 := strings.ToLower(m[1])
	curr2 := strings.ToLower(m[2])
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

	re := regexp.MustCompile("^(.*?)\\s+(.*?)$")
	m := re.FindStringSubmatch(message)
	if m == nil {
		return
	}
	log.WithFields(log.Fields{
		"parsed": m,
	}).Debug("Parsed received message")

	command, ok := cb.commands[m[1]]
	if !ok {
		for _, curr := range cb.currencies {
			if curr == m[1] {
				command = cb.commands["conv"]
			}
		}
	}
	command(m[2], nick, args)
}
