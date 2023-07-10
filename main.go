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
	"github.com/rs/cors"
)

var (
	host              = "http://localhost:"
	restServerPort    *string
	uiServerPort      *string
	urlGui            string
	devMode           *bool
	uiDevelServerPort *string
	develHost         *string
)

// This define the default values of the flags
func initFlags() {
	restServerPort = flag.String("restServerPort", "8080", "RestServerPort to use")
	uiServerPort = flag.String("uiServerPort", "3000", "UiServerPort to use")
	develHost = flag.String("develHost", "http://127.0.0.1:", "DevelHost to use")
	uiDevelServerPort = flag.String("uiDevelServerPort", "5173", "UiDevelServerPort to use")
	devMode = flag.Bool("devMode", false, "DevMode to use")
	flag.Parse()
}

func main() {
	initFlags()

	database.InitDatabaseIfNotExists()

	// Adding multi routes in the same server of the API
	router := mux.NewRouter()
	routesIncomes.SetIncomesHandleActions(router)
	routesExpenses.SetExpensesHandleActions(router)

	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
			*develHost + *uiDevelServerPort,
		},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	fmt.Println("Server API: http://" + host + ":" + *restServerPort + "/")
	fmt.Println("Server UI: http://" + host + ":" + *uiServerPort + "/")

	// Running the async server for the API
	go http.ListenAndServe(":"+*restServerPort, router)

	// This run another server for the UI and launch the UI
	runUI()
}

func corsConfigure(router *mux.Router) {
	//Use this for enable all origins of requests
	//router.Use(cors.AllowAll().Handler)

	//Use this for enable specific origins
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
			*develHost + *uiDevelServerPort,
		},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
}

// This function needs to be at the end of the main function, because it's a blocking function.
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

	//Launch the UI server
	router := mux.NewRouter()
	setUIHandleStaticFiles(router)
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
			*develHost + *uiDevelServerPort,
		},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	if !*devMode {
		go http.ListenAndServe(":"+*uiServerPort, router)
	} else {
		urlGui = *develHost + *uiDevelServerPort + "/"
	}

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

func setUIHandleStaticFiles(router *mux.Router) {
	uiPrefix := "/"
	urlGui = host + *uiServerPort + uiPrefix

	router.PathPrefix(uiPrefix).Handler(
		http.StripPrefix(uiPrefix, http.FileServer(http.Dir("./front-devel-react/dist"))),
	)
}
