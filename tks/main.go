package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/giorgiovilardo/tksgo/internal"
)

//go:embed assets
var embeddedFiles embed.FS

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

	// Serve embedded files
	assetHandler := echo.WrapHandler(http.FileServer(getFileSystem()))
	e.GET("/", assetHandler)
	e.GET("/*", assetHandler)

	e.GET("/all_teams_json", internal.TeamsHandler(matches))
	e.GET("/all_teams", internal.TeamsHtmlHandler(matches))
	e.GET("/last_goals_json", internal.LastGoalsHandler(matches))
	e.GET("/last_goals", internal.LastGoalsHtmlHandler(matches))
	e.Logger.Fatal(e.Start(":1323"))
}

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "assets")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}
