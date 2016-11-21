package main

import (
	"database/sql"
	"github.com/jawher/mow.cli"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tosan88/go-exercise-2/bot"
	"github.com/tosan88/go-exercise-2/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"
)

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

	dbFileName := app.String(cli.StringOpt{
		Name:  "db-file-name",
		Value: "my.db",
		Desc:  "The DB file name to save internal state",
	})

	app.Before = func() {
		if *server == "" || *channel == "" {
			app.PrintHelp()
			log.Fatalln("Server or channel paramaters are not set!")
		}
		if !regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]{2,9}$").MatchString(*botName) {
			app.PrintHelp()
			log.Fatalf("Bot name %v is invalid.\n"+
				"It must contain only alphanumberic characters, dash or underscore; "+
				"the first letter should be a letter from alphabet and the bot name "+
				"should be between 3 and 10 characters\n", *botName)
		}
	}

	app.Action = func() {
		defer func(start time.Time) {
			elapsed := time.Since(start)
			log.Printf("Application finished. It was active %v seconds", elapsed.Seconds())
		}(time.Now())

		sqliteDB, err := sql.Open("sqlite3", *dbFileName)
		if err != nil {
			log.Fatalf("Error opening DB: %v\n", err)
		}
		defer sqliteDB.Close()

		db := storage.DB{DB: sqliteDB}
		err = db.CreateSchema()
		if err != nil {
			log.Fatalf("Error by creating schema: %v\n", err)
			return
		}

		client := bot.New(
			&bot.Conf{
				Server:  *server,
				Channel: *channel,
				BotName: *botName,
			},
			db,
			&http.Client{
				Timeout: 10 * time.Second,
			},
		)

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
