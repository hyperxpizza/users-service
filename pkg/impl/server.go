package impl

import (
	"context"

	"github.com/hyperxpizza/users-service/pkg/config"
	"github.com/hyperxpizza/users-service/pkg/database"
	pb "github.com/hyperxpizza/users-service/pkg/grpc"
)

type UsersServiceServer struct {
	cfg *config.Config
	db  *database.Database
	pb.UnimplementedUsersServiceServer
}

func NewUsersServiceServer(cfgPath string) (*UsersServiceServer, error) {
	return &UsersServiceServer{}, nil
}

func (s *UsersServiceServer) Run() {

}

func (s *UsersServiceServer) GetLoginData(ctx context.Context, req *pb.LoginRequest) (*pb.LoginData, error) {
	var loginData pb.LoginData

	return &loginData, nil
}