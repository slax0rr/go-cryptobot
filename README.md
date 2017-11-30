# IRC cryptocurrency Bot

Simple cryptocurrency bot for Internet Relay Chat.

## Installation

```
go get -u github.com/slax0rr/go-cryptobot
```

## Usage

```
usage: go-cryptobot --nick=NICK --server=SERVER --channel=CHANNEL [<flags>]

Simple cryptocurrency IRC bot

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
      --debug            Verbose mode
  -n, --nick=NICK        Bots nickname, required
  -u, --user=USER        IRC user, if omitted same as 'nick'
  -s, --server=SERVER    IRC server to join to, required
  -c, --channel=CHANNEL  IRC channel to join to, required
  -p, --port=6667        IRC server port, required
      --version          Show application version.
```

### On IRC

When the bot joins the channel, simply message it with the currency pair:

```
botNickname: btc usd
```

More commands will follow.
