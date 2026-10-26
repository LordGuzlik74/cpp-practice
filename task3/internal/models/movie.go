package models

import (
	"errors"
)

type Movie struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}
type MovieError struct {
	ID    int    `json:"id"`
	Error string `json:"error"`
}
type LoadMoviesResult struct {
	Movies []Movie      `json:"movies"`
	Errors []MovieError `json:"errors"`
}

var ErrNotFound = errors.New("movie not found")
