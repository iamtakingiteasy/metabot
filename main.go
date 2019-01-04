//go:generate go generate ./model
package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/iamtakingiteasy/metabot/bot/impl"
)

func main() {
	if len(os.Args) <= 1 {
		pref := ""
		if len(os.Args) > 0 {
			pref = os.Args[0] + " "
		}
		log.Fatalln(pref + "filename.yaml required")
	}

	ctx := impl.NewContext(os.Args[1])

	err := ctx.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = ctx.SaveConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = ctx.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Connected %s v%d as %s#%s\n", ctx.Session.State.SessionID, ctx.Session.State.Version, ctx.Session.State.User.Username, ctx.Session.State.User.Discriminator)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch
	ctx.Stop()
	close(ch)
}
