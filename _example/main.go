package main

import (
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
		outis.WithID("422138b3-c721-4021-97ab-8cf7e174fb4f"),
		outis.WithInterval(
			outis.NewInterval(
				outis.WithHours(18, 19),
				outis.WithMinutes(0, 10),
				// outis.WithEvery(time.Second*10),
			),
		),
		outis.WithName("-PlanUpdater"),
		outis.WithDesc("Atualiza planos com data de validade expiradas, tornando a venda nao visivel."),
		outis.WithScript(func(ctx outis.Context) error {
			ctx.LogInfo("this is an information message")

			return nil
		}),
	)

	// Método que mantém a rotina no processo
	watch.Wait()
}
