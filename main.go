package main

import (
	r "main/internal/route"
)

func main() {
	route := r.SetupRouter()
	route.Run("localhost:8080")
}
