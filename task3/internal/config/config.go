package config

import (
	"os"
)

type Config struct {
	MovieBaseUrl string
	HttpAddr     string
	GrpcAddr     string
}

func Load() Config {
	baseUrl := os.Getenv("MOVIE_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://homeworksite.site"
	}

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":9090"
	}

	return Config{
		MovieBaseUrl: baseUrl,
		HttpAddr:     httpAddr,
		GrpcAddr:     grpcAddr,
	}
}
