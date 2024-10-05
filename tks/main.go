package main

import (
	"fmt"

	"github.com/giorgiovilardo/tksgo/internal"
)

func main() {
	conf := internal.LoadConf()
	fmt.Println("Loaded configuration:")
	for _, league := range conf.Leagues {
		fmt.Printf("%s URL: %s\n", league.Name, league.URL)
	}
}
