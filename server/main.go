package main

import (
	"log"

	"github.com/twist"

	"github.com/gofiber/fiber"
)

func main() {

	//Init Example
	twist.CreateAccount("fred", 5000)
	twist.CreateAccount("armani", 1000)

	// Fiber instance
	app := fiber.New()

	// APIs
	app.Get("/wallets", twist.Wallets)

	app.Post("/deduct", twist.DeductTry)
	app.Put("/deduct", twist.DeductConfirm)
	app.Delete("/deduct", twist.DeductCancel)

	app.Post("/deposit", twist.Deposit)

	// Start server
	log.Fatal(app.Listen(3000))

}
