package main

import (
	"github.com/2024_2_BetterCallFirewall/internal/app/auth"
)

func main() {
	if err := auth.Run(); err != nil {
		panic(err)
	}
}
