package service

import (
	"context"
	"errors"
	"sort"
	"sync"
	"task3/internal/client"
	"task3/internal/models"
)

type Service struct {
	client *client.Client
}

func New(c *client.Client) *Service {
	return &Service{client: c}
}

func (s *Service) GetLastMovie(ctx context.Context) (models.Movie, error) {
	id := 1
	LastFoundMovie, err := s.client.GetMovie(ctx, id)
	if err != nil {
		return models.Movie{}, err
	}
	low := 1
	high := 2

	for {
		movie, err := s.client.GetMovie(ctx, high)
		if errors.Is(err, models.ErrNotFound) {
			break
		} else if err != nil {
			return models.Movie{}, err
		} else if err == nil {
			low = high
			high *= 2
			LastFoundMovie = movie
		}
	}
	for low+1 < high {
		mid := low + (high-low)/2
		movie, err := s.client.GetMovie(ctx, mid)
		if errors.Is(err, models.ErrNotFound) {
			high = mid
		} else if err != nil {
			return models.Movie{}, err
		} else if err == nil {
			low = mid
			LastFoundMovie = movie
		}
	}

	return LastFoundMovie, nil
}

const defaultConcurrency = 10

func (s *Service) LoadMovies(ctx context.Context, from int, to int) models.LoadMoviesResult {
	total := to - from + 1
	jobs := make(chan int, total)

	var res models.LoadMoviesResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	workers := defaultConcurrency
	if total < workers {
		workers = total
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				if err := ctx.Err(); err != nil {
					mu.Lock()
					res.Errors = append(res.Errors, models.MovieError{
						ID:    id,
						Error: err.Error(),
					})
					mu.Unlock()
					continue
				}

				movie, err := s.client.GetMovie(ctx, id)
				mu.Lock()
				if err != nil {
					res.Errors = append(res.Errors, models.MovieError{
						ID:    id,
						Error: err.Error(),
					})
				} else {
					res.Movies = append(res.Movies, movie)
				}
				mu.Unlock()
			}
		}()
	}

	for id := from; id <= to; id++ {
		jobs <- id
	}
	close(jobs)

	wg.Wait()
	sort.Slice(res.Movies, func(i, j int) bool {
		return res.Movies[i].Id < res.Movies[j].Id
	})
	return res
}
