package main

import (
	"regexp"
	"strconv"
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
		"btc usd",
		"btc eur",
		"xrp usd",
		"xrp eur",
		"xrp btc",
		"ltc usd",
		"ltc eur",
		"ltc btc",
		"eth usd",
		"eth eur",
		"eth btc",
	}
	cb.commands = map[string]func(string, string, []string){
		"conv":  cb.conv,
		"help":  cb.printHelp,
		"pairs": cb.printPairs,
	}
	return cb
}

func (cb *CryptoBot) printPairs(message, nick string, args []string) {
	cb.irc.WritePriv(nick, "Available currency pairs for conversion:")
	cb.irc.WritePriv(nick, strings.Join(cb.pairs, ", "))
}

func (cb *CryptoBot) printHelp(message, nick string, args []string) {
	var cmds string
	for k := range cb.commands {
		cmds = cmds + k + ","
	}
	cmds = strings.TrimRight(cmds, ",")
	cb.irc.WritePriv(nick, "Available commands:")
	cb.irc.WritePriv(nick, cmds)
}

func (cb *CryptoBot) conv(message, nick string, args []string) {
	log.WithFields(log.Fields{
		"message": message,
	}).Debug("Received conversion message")

	re := regexp.MustCompile("^(.+)\\s+(.+)$")
	m := re.FindStringSubmatch(message)
	if m == nil {
		log.Error("Conversion did not parse properly.")
		return
	}

	amount, curr1 := parseAmount(m[1])
	curr2 := m[2]
	if cb.isValidPair(curr1, curr2) == false {
		log.WithFields(log.Fields{
			"currency1": curr1,
			"currency2": curr2,
		}).Error("Received data is not valid for currency conversion")
		return
	}

	resp := cb.client.GetTicker(curr1, curr2)
	if resp.Err != nil {
		log.WithFields(log.Fields{
			"error": resp.Err,
		}).Debug("Ticker responded with error")
		cb.irc.Write(nick + ": Unable to retrieve ticker data due to an unknown error")
		return
	}

	var msg string
	if amount == 0 {
		msg = nick + ": " + curr1 + " to " + curr2 + ": Last: " + resp.Last +
			" High: " + resp.High +
			" Low: " + resp.Low +
			" Open: " + resp.Open
	} else {
		last, _ := strconv.ParseFloat(resp.Last, 64)
		strAmount := strconv.FormatFloat(amount, 'f', 8, 64)
		val := strconv.FormatFloat(last*amount, 'f', 2, 64)
		msg = nick + ": " + strAmount + " " + curr1 + " to " + curr2 + ": " + val + " @ " + resp.Last
	}
	cb.irc.Write(msg)
}

func (cb *CryptoBot) evHandler(message, nick string, args []string) {
	log.WithFields(log.Fields{
		"msg": message,
	}).Debug("Message received")

	message = strings.ToLower(message)
	re := regexp.MustCompile("^(\\w+)\\s*(.*?)$")
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
		command = cb.commands["conv"]
		cmdMsg = m[0]
	}

	log.WithFields(log.Fields{
		"command": cmdMsg,
	}).Debug("Executing command with message")

	command(cmdMsg, nick, args)
}

func (cb *CryptoBot) isCurrency(currency string) bool {
	for _, curr := range cb.currencies {
		if curr == currency {
			return true
		}
	}
	return false
}

func (cb *CryptoBot) isValidPair(curr1, curr2 string) bool {
	if cb.isCurrency(curr1) == false || cb.isCurrency(curr2) == false {
		return false
	}
	p := curr1 + " " + curr2
	for _, pair := range cb.pairs {
		if pair == p {
			return true
		}
	}
	return false
}

func parseAmount(data string) (float64, string) {
	re := regexp.MustCompile("^([\\d.]+)\\s+(.+)$")
	m := re.FindStringSubmatch(data)
	if m == nil {
		return 0, data
	}
	f, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		f = 0
	}
	return f, m[2]
}
