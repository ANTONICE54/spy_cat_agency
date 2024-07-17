package main

import (
	"spy_cat_agency/internal/server"
)

func main() {
	server := server.NewServer()
	server.Run()
}
