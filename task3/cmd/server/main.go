package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "task3/docs"
	"task3/internal/client"
	"task3/internal/config"
	"task3/internal/service"
	grpcserver "task3/internal/transport/grpc"
	httpserver "task3/internal/transport/http"
	"task3/proto/moviepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// @title           Movie Stand Proxy API
// @version         1.0
// @description     HTTP API для загрузки фильмов с внешнего стенда.
// @BasePath        /
func main() {
	// 1. Загружаем конфигурацию
	cfg := config.Load()

	// 2. Инициализируем клиент и сервисный слой
	movieClient := client.New(cfg.MovieBaseUrl)
	movieService := service.New(movieClient)

	// 3. Настраиваем HTTP-сервер
	httpHandler := httpserver.New(movieService)
	httpSrv := &http.Server{
		Addr:    cfg.HttpAddr,
		Handler: httpHandler.Router(),
	}

	// 4. Настраиваем gRPC-сервер
	lis, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		log.Fatalf("не удалось занять порт для gRPC (%s): %v", cfg.GrpcAddr, err)
	}

	grpcSrv := grpc.NewServer()
	moviepb.RegisterMovieServiceServer(grpcSrv, grpcserver.New(movieService))
	// Включаем reflection, чтобы gRPC-сервер можно было легко тестировать утилитой grpcurl или Postman
	reflection.Register(grpcSrv)

	// 5. Контекст для перехвата Ctrl+C (SIGINT / SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 6. Запускаем HTTP-сервер в отдельной горутине
	go func() {
		fmt.Println("HTTP сервер запущен на", cfg.HttpAddr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ошибка HTTP сервера: %v", err)
		}
	}()

	// 7. Запускаем gRPC-сервер в отдельной горутине
	go func() {
		fmt.Println("gRPC сервер запущен на", cfg.GrpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("ошибка gRPC сервера: %v", err)
		}
	}()

	// 8. Ждём сигнала остановки (Ctrl+C)
	<-ctx.Done()
	fmt.Println("\nПолучен сигнал остановки, завершаем работу серверов...")

	// Даём 5 секунд на завершение текущих запросов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("ошибка остановки HTTP сервера: %v", err)
	}
	grpcSrv.GracefulStop()

	fmt.Println("Серверы успешно остановлены.")
}
