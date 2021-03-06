package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	database "github.com/TitusW/team4-kargo-excellerate/db"

	"github.com/TitusW/team4-kargo-excellerate/handlers"
	gohandlers "github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	l := log.New(os.Stdout, "", log.LstdFlags)

	err := initEnv(l)
	if err != nil {
		log.Fatalf("Error ENV: %s", err)
	}

	err = checkEnv(l)
	if err != nil {
		log.Fatalf("Error ENV: %s", err)
	}

	db := database.Connect(l)
	defer db.Close(l)

	port := ":" + os.Getenv("PORT")

	CORSHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	router := mux.NewRouter()
	driverh := handlers.NewDriverHandler(db.Connection)
	truckh := handlers.NewTruckHandler(db.Connection)

	router.HandleFunc("/api/v1/drivers", driverh.GetDrivers).Methods("GET")
	router.HandleFunc("/api/v1/driver/{id}", driverh.GetDriverById).Methods("GET")
	router.HandleFunc("/api/v1/driver", driverh.CreateDriver).Methods("POST")
	router.HandleFunc("/api/v1/driver/{id}", driverh.UpdateDriver).Methods("PUT")

	router.HandleFunc("/api/v1/trucks", truckh.GetTrucks).Methods("GET")
	router.HandleFunc("/api/v1/truck/{id}", truckh.GetTruckById).Methods("PUT")
	router.HandleFunc("/api/v1/truck", truckh.CreateTruck).Methods("POST")
	router.HandleFunc("/api/v1/truck/{id}", truckh.UpdateTruck).Methods("PUT")

	srv := &http.Server{
		Addr:         port,
		Handler:      CORSHandler(router),
		ErrorLog:     l,
		ReadTimeout:  time.Duration(5) * time.Second,
		WriteTimeout: time.Duration(10) * time.Second,
		IdleTimeout:  time.Duration(120) * time.Second,
	}

	go func() {

		l.Printf("Starting server on port %s\n", port)

		err := srv.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	sig := <-c
	l.Println("Got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	l.Println("Shutting down...")
	os.Exit(0)
}

func initEnv(l *log.Logger) error {
	l.Println("Checking .env file")

	if err := godotenv.Load(); err != nil {
		return err
	}

	return nil
}

func checkEnv(l *log.Logger) error {
	l.Println("Checking environment")

	if _, found := os.LookupEnv("DATABASE_HOST"); !found {
		return errors.New("DATABASE_HOST value does not exist")
	}

	if _, found := os.LookupEnv("DATABASE_PORT"); !found {
		return errors.New("DATABASE_PORT value does not exist")
	}

	if _, found := os.LookupEnv("DATABASE_USER"); !found {
		return errors.New("DATABASE_USER value does not exist")
	}

	if _, found := os.LookupEnv("DATABASE_PASS"); !found {
		return errors.New("DATABASE_PASS value does not exist")
	}

	if _, found := os.LookupEnv("PORT"); !found {
		return errors.New("PORT value does not exist")
	}

	l.Println("Environment check passed")
	return nil
}
