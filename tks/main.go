package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/giorgiovilardo/tksgo/internal"
)

func main() {
	conf := internal.LoadConf()
	fmt.Println("Loaded configuration:")
	fmt.Printf("%+v\n", conf)
	for _, league := range conf.Leagues {
		fmt.Printf("%s URL: %s\n", league.Name, league.URL)
	}
	matches, err := internal.GetMatchesFromCsv(conf)
	if err != nil {
		fmt.Println("Error downloading CSV files:", err)
		return
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/all_teams", internal.TeamsHandler(matches))
	e.GET("/last_goals", internal.LastGoalsHandler(matches))
	e.Logger.Fatal(e.Start(":1323"))

}
