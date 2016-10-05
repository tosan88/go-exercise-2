package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

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
		Desc:  "The nickname for the bot which will be joined to channel. Defaults to 'test-bot'",
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
			log.Printf("Application finished. It was active %v seconds", elapsed.Seconds())
		}(time.Now())

		client := NewClient(&conf{
			server:  *server,
			channel: *channel,
			botName: *botName,
		})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			client.Run()
			wg.Done()
		}()

		waitForQuitSignal()
		client.Stop()
		wg.Wait()
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func waitForQuitSignal() {
	shutDownCh := make(chan os.Signal) //should I move this to botClient & get rid off shouldStop?
	signal.Notify(shutDownCh, syscall.SIGINT, syscall.SIGTERM)

	<-shutDownCh
}
