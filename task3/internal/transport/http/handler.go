package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	_ "task3/docs"
	"task3/internal/models"
	"task3/internal/service"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	service *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /last-movie", h.GetLastMovie)
	mux.HandleFunc("GET /load-movies", h.LoadMovies)
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)
	return mux
}

// GetLastMovie находит последний доступный фильм на стенде
// @Summary      Получить последний фильм
// @Description  Ищет фильм с максимальным ID на внешнем стенде
// @Tags         movies
// @Produce      json
// @Success      200  {object}  models.Movie
// @Failure      404  {string}  string  "Фильмы не найдены"
// @Failure      502  {string}  string  "Ошибка внешнего стенда"
// @Failure      504  {string}  string  "Таймаут внешнего стенда"
// @Router       /last-movie [get]
func (h *Handler) GetLastMovie(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	movie, err := h.service.GetLastMovie(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Movie not fiund", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get last movie: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(movie)
}

// LoadMovies загружает фильмы в диапазоне [from; to]
// @Summary      Загрузить диапазон фильмов
// @Description  Параллельно загружает фильмы со стенда по диапазону ID [from; to]
// @Tags         movies
// @Produce      json
// @Param        from  query     int  true  "Начальный ID фильма (>= 1)"
// @Param        to    query     int  true  "Конечный ID фильма (>= from)"
// @Success      200   {array} models.Movie
// @Failure      400   {string}  string  "Некорректные параметры from или to"
// @Router       /load-movies [get]
func (h *Handler) LoadMovies(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err1 := strconv.Atoi(fromStr)
	to, err2 := strconv.Atoi(toStr)

	if err1 != nil || err2 != nil || from < 1 || to < 1 || from > to {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	result := h.service.LoadMovies(ctx, from, to)

	movies := result.Movies
	if movies == nil {
		movies = []models.Movie{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(movies)
}
