package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"os"
	"time"
)

//TODO use it
type conf struct {
	server  string
	channel string
	botName string
}

func main() {
	log.Printf("Application starting with args %s", os.Args)
	app := cli.App("go-exercise-2", "Exercising go skills 2")

	server := app.String(cli.StringOpt{
		Name:  "server",
		Value: "",
		Desc:  "The IRC server address",
	})

	channel := app.String(cli.StringOpt{
		Name:  "channel",
		Value: "",
		Desc:  "The channel name to join to",
	})

	botName := app.String(cli.StringOpt{
		Name:  "bot-name",
		Value: "test-bot",
		Desc:  "The nickname for the bot which will be joined to channel",
	})

	app.Before = func() {
		if *server == "" || *channel == "" {
			app.PrintHelp()
			log.Fatalln("Server or channel paramaters are not set!")
		}
	}

	app.Action = func() {
		defer func(start time.Time) {
			elapsed := time.Since(start)
			log.Printf("Application finished. Took %v seconds", elapsed.Seconds())
		}(time.Now())

		log.Printf("Application started with bot name: %v\n", *botName)

	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
