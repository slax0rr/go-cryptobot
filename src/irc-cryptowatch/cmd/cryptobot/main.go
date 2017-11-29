package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var build = "0.1.0-dev"

var app = kingpin.New("cryptobot", "Simple cryptocurrency IRC bot")
var verb = app.Flag("debug", "Verbose mode").Default("false").Bool()

var nick = app.Flag("nick", "Bots nickname, required").Short('n').Required().String()
var server = app.Flag("server", "IRC server to join to, required").Short('s').Required().String()
var channel = app.Flag("channel", "IRC channel to join to, required").Short('c').Required().String()
var port = app.Flag("port", "IRC server port").Short('p').Default("6667").Int()

func main() {
	app.Author("Tomaz Lovrec <tomaz.lovrec@gmail.com>")
	app.Version(build)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *verb {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")
	}

	log.Info("Foo")
}
