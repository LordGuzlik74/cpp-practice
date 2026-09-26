package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Film struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func ShowFilm(ctx context.Context, id int, timeout time.Duration) {
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	url := fmt.Sprintf("http://homeworksite.site/%v/info.0.json", id)
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("ошибка запроса")
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("id %d: сетевая ошибка / таймаут: %v\n", id, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println("api is not avaliable")
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("id %d: ошибка чтения: %v\n", id, err)
		return
	}
	var film Film
	err = json.Unmarshal(body, &film)
	if err != nil {
		fmt.Printf("id %d: битый json: %v\n", id, err)
		return
	}
	id, title, year, director := film.Id, film.Title, film.Year, film.Director
	fmt.Printf("%d - %s - %d - %s\n", id, title, year, director)
}

type WorkerPool struct {
	concurrency int
	FilmsChan   chan int
	wg          sync.WaitGroup
	timeout     time.Duration
}

func (wp *WorkerPool) worker(ctx context.Context) {
	for task := range wp.FilmsChan {
		if ctx.Err() != nil {
			wp.wg.Done()
			continue
		}
		ShowFilm(ctx, task, wp.timeout)
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Run(ctx context.Context, From int, To int) {

	wp.FilmsChan = make(chan int, To-From+1)
	for i := 0; i < wp.concurrency; i++ {
		go wp.worker(ctx)
	}
	wp.wg.Add(To - From + 1)
	for i := From; i < To+1; i++ {
		select {
		case <-ctx.Done():
			wp.wg.Done()
		case wp.FilmsChan <- i:
		}
	}
	close(wp.FilmsChan)
	wp.wg.Wait()
}
