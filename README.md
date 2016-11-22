[![CircleCI](https://circleci.com/gh/tosan88/go-exercise-2.svg?style=shield)](https://circleci.com/gh/tosan88/go-exercise-2)
[![Coverage Status](https://coveralls.io/repos/github/tosan88/go-exercise-2/badge.svg?branch=master)](https://coveralls.io/github/tosan88/go-exercise-2?branch=master)
# go-exercise-2
An IRC bot which reacts to user interactions with an IRC channel.

## Install & run

```
go get -u github.com/tosan88/go-exercise-2
go-exercise-2 --server="chat.freenode.net:6667" --channel="go-test-bot" --bot-name="test-bot" --db-file-name="my.db"
```

## Test

```
go test -v -race ./...
```

### Functional testing

You could install an IRC client of your choice (e.g. Circ, a Google Chrome app) and see how the bot behaves 
by joining to the same server and channel as the bot.