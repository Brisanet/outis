package main

import (
	"errors"
	"time"

	"github.com/Brisanet/outis"
)

func main() {
	log, err := outis.NewLogger("v1/example")
	if err != nil {
		log.Fatal(err.Error())
	}
	// Initialize Outis to be able to add routines
	watch := outis.Watcher("8b1d6a18-5f3d-4482-a574-35d3965c8783", "v1/example") // outis.Logger(),                // Option to implement logs interface
	// outis.WithOutisInterface(nil), // Option to implement outis interface

	go watch.Go(
		// Routine identifier to perform concurrency control
		outis.WithID("422138b3-c721-4021-97ab-8cf7e174fb4f"),

		outis.WithName("Here is the name of my routine"),
		outis.WithDesc("Here is the description of my routine"),

		// It will run every 10 second
		outis.WithInterval(time.Second),
		// It will run from 12pm to 4pm.
		// by default, there are no time restrictions.
		// outis.WithHours(12, 16),

		// Here the script function that will be executed will be passed
		outis.WithScript(func(ctx *outis.Context) error {
			log.Info("aaaaaaaaaaaaaaaaaaa")
			ctx.Info("this is an information message")
			ctx.Error(errors.New("this is an error message"))

			ctx.AddSingleMetadata("client_ids", []int64{234234})
			ctx.AddMetadata(outis.Metadata{"notification": outis.Metadata{
				"client_id": 234234,
				"message":   "Hi, we are notifying you.",
				"fcm":       "3p2okrmionfiun2uni3nfin2i3f",
			}})

			ctx.Debug("Hello")

			return nil
		}),
	)

	// Method that maintains routine in the process
	watch.Wait()
}
