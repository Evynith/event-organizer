package main

import (
	router "main/internal/route"
)

func main() {
	route := router.SetupRouter()
	route.Run(":8000")
}
