package grpcserver

import (
	"context"
	"errors"
	"task3/internal/models"
	"task3/internal/service"
	"task3/proto/moviepb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	moviepb.UnimplementedMovieServiceServer
	service *service.Service
}

func New(s *service.Service) *Server {
	return &Server{service: s}
}

func (s *Server) GetLastMovie(ctx context.Context, req *moviepb.GetLastMovieRequest) (*moviepb.MovieResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	m, err := s.service.GetLastMovie(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "external movie stand timeout")
		}
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "no movies found on stand")
		}
		return nil, status.Error(codes.Unavailable, "failed to get last movie: "+err.Error())
	}

	return &moviepb.MovieResponse{
		Movie: toProtoMovie(m),
	}, nil
}

func (s *Server) LoadMovies(ctx context.Context, req *moviepb.LoadMoviesRequest) (*moviepb.LoadMoviesResponse, error) {
	from := int(req.GetFrom())
	to := int(req.GetTo())

	if from < 1 || to < 1 || from > to {
		return nil, status.Error(codes.InvalidArgument, "invalid 'from' or 'to' parameters")
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res := s.service.LoadMovies(ctx, from, to)

	pbMovies := make([]*moviepb.Movie, 0, len(res.Movies))
	for _, m := range res.Movies {
		pbMovies = append(pbMovies, toProtoMovie(m))
	}

	pbErrors := make([]*moviepb.MovieError, 0, len(res.Errors))
	for _, e := range res.Errors {
		pbErrors = append(pbErrors, &moviepb.MovieError{
			Id:    int32(e.ID),
			Error: e.Error,
		})
	}

	return &moviepb.LoadMoviesResponse{
		Movies: pbMovies,
		Errors: pbErrors,
	}, nil
}

func toProtoMovie(m models.Movie) *moviepb.Movie {
	return &moviepb.Movie{
		Id:       int32(m.Id),
		Title:    m.Title,
		Year:     int32(m.Year),
		Director: m.Director,
	}
}
