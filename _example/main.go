package main

import (
	"errors"
	"time"

	"github.com/Brisanet/outis"
)

func main() {
	// Inicializa o log
	log, err := outis.NewLogger("scriptName", outis.LogOptions{
		Level: outis.DebugLevel,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Inicializa o outis para receber rotinas
	watch := outis.Watcher("8b1d6a18-5f3d-4482-a574-35d3965c8783", "scriptName",
		// Passa o log personalizado, se não informado é criado um log padrão
		outis.Logger(log),
	)

	watch.Go(
		// Identificador de rotina para executar controle de concorrência
		outis.WithID("422138b3-c721-4021-97ab-8cf7e174fb4f"),

		outis.WithName("Here is the name of my routine"),
		outis.WithDesc("Here is the description of my routine"),

		// Executará a cada 10 segundos
		outis.WithInterval(time.Second),

		// Executará de 12pm a 4pm.
		// por padrão, não há restrições de tempo.
		// outis.WithHours(12, 16),

		// Executará somente uma vez
		// outis.WithNotUseLoop(),

		// Aqui é passada a função do script que será executada
		outis.WithScript(func(ctx outis.Context) error {
			ctx.LogInfo("this is an information message")
			ctx.LogError(errors.New("this is an error message"))

			ctx = ctx.AddSingleMetadata("client_ids", []int64{234234})
			ctx = ctx.AddMetadata(outis.Metadata{"notification": outis.Metadata{
				"client_id": 234234,
				"message":   "Hi, we are notifying you.",
				"fcm":       "231223",
			}})

			ctx.LogDebug("this is an debug message with metadata")

			return nil
		}),
	)

	// Método que mantém a rotina no processo
	watch.Wait()
}
