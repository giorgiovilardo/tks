package main

import (
	"fmt"

	"github.com/giorgiovilardo/tksgo/internal"
)

func main() {
	conf := internal.LoadConf()
	fmt.Println("Loaded configuration:")
	fmt.Printf("%+v\n", conf)
	for _, league := range conf.Leagues {
		fmt.Printf("%s URL: %s\n", league.Name, league.URL)
	}
	csvs, err := internal.GetMatchesFromCsv(conf)
	if err != nil {
		fmt.Println("Error downloading CSV files:", err)
		return
	}
	fmt.Printf("Downloaded CSV files: %+v\n", csvs)
	fmt.Printf("Match1: %+v\n", csvs[0].IdempotentKey())
}
