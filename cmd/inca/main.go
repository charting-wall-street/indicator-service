package main

import (
	"fmt"
	"inca/internal/config"
	"inca/internal/web"
)

func main() {
	fmt.Println("|- Indicator Service -|")
	config.LoadConfig()
	web.Start()
}
