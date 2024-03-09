package main

import (
	"context"
	"crebito-fiber/controllers"
	"crebito-fiber/helpers"
	"log"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load(".env")
	helpers.SetupLog()	
	logrus.Info(">>>>>>>>>>   " + os.Getenv("SERVER_NAME")  + "   <<<<<<<<<< ")
	
	helpers.GetDBConnection()
	logrus.Info("[TEST] DB Connection OK")
	// 
	resetDB()

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	setupRoutes(app)

	log.Fatal(app.Listen(":" + os.Getenv("WEB_PORT") ))
}


func setupRoutes(app *fiber.App) {

    app.Get("/clientes/:ID/extrato", func(c *fiber.Ctx) error {
		result := controllers.NewExtratoController(c).Run()
		return c.Status(result.StatusCode).JSON(result)
    })

	app.Post("/clientes/:ID/transacoes", func(c *fiber.Ctx) error {
		result := controllers.NewTransacaoController(c).Run()
		return c.Status(result.StatusCode).JSON(result)
    })

	logrus.Info("Routes OK")
}


func resetDB() {

	_, err := helpers.GetDBConnection().Exec(context.Background(), `truncate table transacoes;`)
	if err != nil {
		logrus.Error("[ResetDB.TruncateTransacoes]", err.Error())
		os.Exit(1)
	}
	_, err = helpers.GetDBConnection().Exec(context.Background(), `update clientes set saldo = 0;`)
	if err != nil {
		logrus.Error("[ResetDB.UpdateClientesSaldo]", err.Error())
		os.Exit(1)
	}

	logrus.Debug("[RESET DB] OK")
}
