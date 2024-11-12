package main

import (
	"fmt"
	"log"

	"github.com/RafaelTauschek/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	err = cfg.SetUser("Rafael")
	if err != nil {
		log.Fatal(err)
	}

	updatedConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(updatedConfig.DBUrl)
	fmt.Println(updatedConfig.CurrentUser)

}
