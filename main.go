package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/LuisFlahan4051/basic-counter-app/database"
	routesExpenses "github.com/LuisFlahan4051/basic-counter-app/routes/expenses"
	routesIncomes "github.com/LuisFlahan4051/basic-counter-app/routes/incomes"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/gorilla/mux"
)

var (
	restServerPort *string
	URLs           []string
	urlGui         string
	host           = "http://localhost:"
	devMode        = true
)

func main() {
	initFlags()

	database.InitDatabaseIfNotExists()

	// Adding multi routes in the same server
	router := mux.NewRouter()
	routesIncomes.SetIncomesHandleActions(router)
	routesExpenses.SetExpensesHandleActions(router)
	setUIHandleStaticFiles(router, "5173")

	// Running the async server
	go http.ListenAndServe(":"+*restServerPort, router)

	runUI()
}
func setUIHandleStaticFiles(router *mux.Router, develPort string) {
	uiPrefix := "/"
	urlGui = host + *restServerPort + uiPrefix
	if develPort != "" && devMode {
		urlGui = host + develPort + uiPrefix
	}
	router.PathPrefix(uiPrefix).Handler(
		http.StripPrefix(uiPrefix, http.FileServer(http.Dir("./front-devel-react/dist"))),
	)
}

func runUI() {
	//App initialization
	logEntry := log.New(log.Writer(), log.Prefix(), log.Flags())
	app, err := astilectron.New(logEntry, astilectron.Options{
		AppName:            "Basic Counter App",
		AppIconDefaultPath: "icon.png",  // If path is relative, it must be relative to the data directory
		AppIconDarwinPath:  "icon.icns", // Same here
		BaseDirectoryPath:  "dependencies",
	})
	if err != nil {
		logEntry.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer app.Close()
	app.Start()

	//MainWindow initialization
	var mainWindow *astilectron.Window
	if mainWindow, err = app.NewWindow(urlGui, &astilectron.WindowOptions{
		Center:    astikit.BoolPtr(true),
		Height:    astikit.IntPtr(680),
		MinHeight: astikit.IntPtr(670),
		//Width:     astikit.IntPtr(500),
		Width:    astikit.IntPtr(1150),
		MinWidth: astikit.IntPtr(1140),
		//Frame:     astikit.BoolPtr(false),
		Resizable: astikit.BoolPtr(true),
	}); err != nil {
		logEntry.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}
	if err = mainWindow.Create(); err != nil {
		logEntry.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	app.Wait()
}

func initFlags() {
	restServerPort = flag.String("restServerPort", "8080", "RestServerPort to use")

	flag.Parse()
}
