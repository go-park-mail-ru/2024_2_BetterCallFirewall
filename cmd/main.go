package main

import (
	"github.com/2024_2_BetterCallFirewall/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
