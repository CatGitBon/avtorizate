package main

import (
	"context"
	"errors"
	"fmt"
	avtorizate "github.com/CatGitBon/avtorizate/pkg"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

// Пример ключа для подписи JWT
var jwtKey = []byte("secret_key")

// Структура для хранения данных о пользователе
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type server struct {
	avtorizate.UnimplementedAvtorizateServer
}

func main() {

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем наш avtorizate server
	avtorizate.RegisterAvtorizateServer(grpcServer, &server{})

	// Запускаем сервер
	fmt.Println("avtorizate Service listening on port :50051")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
func (s *server) AvtorizateMethod(ctx context.Context, req *avtorizate.AvtorizateRequest) (*avtorizate.AvtorizateResponse, error) {

	tokenString := req.Token

	if tokenString == "" {
		return &avtorizate.AvtorizateResponse{
			Success: false,
			Message: "",
		}, errors.New("Missing token")
	}

	// Обрезаем "Bearer " из заголовка (если он есть)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return &avtorizate.AvtorizateResponse{
			Success: false,
			Message: "",
		}, errors.New("Invalid token")
	}

	// Извлекаем информацию о пользователе из токена
	_, ok := token.Claims.(*Claims)
	if !ok {
		return &avtorizate.AvtorizateResponse{
			Success: false,
			Message: "",
		}, errors.New("Invalid token claims")
	}

	return &avtorizate.AvtorizateResponse{
		Success: true,
		Message: "",
	}, nil
}
