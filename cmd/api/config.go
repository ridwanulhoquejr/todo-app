package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const version = "1.0.0"

var environments = [2]string{"dev", "production"}

// config for our application
type config struct {
	port int
	env  string
}

func (c *config) configs() *config {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Now you can get the environment variable from the .env file
	portStr := os.Getenv("API_PORT")
	if portStr == "" {
		log.Fatal("API_PORT not set in .env file")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("cannot convert API port string to int")
	}

	return &config{
		port: port,
		env:  environments[0],
	}
}
