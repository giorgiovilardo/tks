package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/labstack/echo/v4"

	"github.com/giorgiovilardo/tksgo/internal"
)

//go:embed assets
var embeddedFiles embed.FS

func main() {
	conf := internal.LoadConf()
	matches, err := internal.GetMatchesFromCsv(conf)
	if err != nil {
		fmt.Println("Error downloading CSV files:", err)
		return
	}

	e := echo.New()

	assetHandler := echo.WrapHandler(http.FileServer(getFileSystem()))
	e.GET("/", assetHandler)
	e.GET("/*", assetHandler)

	e.GET("/all_teams_json", internal.TeamsHandler(matches))
	e.GET("/all_teams", internal.TeamsHtmlHandler(matches))
	e.GET("/last_goals_json", internal.LastGoalsHandler(matches))
	e.GET("/last_goals", internal.LastGoalsHtmlHandler(matches))
	e.GET("/last_matches_json", internal.LastMatchesHandler(matches))
	e.GET("/result_matrix", internal.ResultMatrixHandler)

	go func() {
		url := "http://localhost:1323"
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		case "darwin":
			cmd = exec.Command("open", url)
		default:
			return
		}
		err := cmd.Start()
		if err != nil {
			fmt.Printf("Error opening browser: %v\n", err)
		}
	}()

	e.Logger.Fatal(e.Start(":1323"))
}

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "assets")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}
