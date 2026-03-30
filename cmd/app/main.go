package main

import (
	"log"
	appbootstrap "sistema_venta_pasajes/configs"
)

func main() {
	if err := appbootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
