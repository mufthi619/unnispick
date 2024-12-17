package api

import "log"

func StartAPI() {
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}
