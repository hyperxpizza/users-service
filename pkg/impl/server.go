package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/hyperxpizza/users-service/pkg/config"
	"github.com/hyperxpizza/users-service/pkg/database"
	pb "github.com/hyperxpizza/users-service/pkg/grpc"
	"github.com/hyperxpizza/users-service/pkg/utils"
	"github.com/hyperxpizza/users-service/pkg/validator"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	UsernameAlreadyExistsError = "username: %s already exists in the database"
	EmailAlreadyExistsError    = "email: %s already exists in the database"
)

type UsersServiceServer struct {
	cfg    *config.Config
	db     *database.Database
	logger logrus.FieldLogger
	pb.UnimplementedUsersServiceServer
}

func NewUsersServiceServer(cfgPath string, logger logrus.FieldLogger) (*UsersServiceServer, error) {

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	db, err := database.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return &UsersServiceServer{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}, nil
}

func (s *UsersServiceServer) Run() {
	grpcServer := grpc.NewServer()
	pb.RegisterUsersServiceServer(grpcServer, s)
	addr := fmt.Sprintf(":%d", s.cfg.UsersService.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Fatalf("net.Listen failed: %s", err.Error())
	}

	s.logger.Infof("users service server running on %s:%d", s.cfg.UsersService.Host, s.cfg.UsersService.Host)

	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Fatalf("failed to serve: %s", err.Error())
	}
}

func (s *UsersServiceServer) GetLoginData(ctx context.Context, req *pb.LoginRequest) (*pb.LoginData, error) {

	s.logger.Infof("getting login data for: %s", req.Username)

	loginData, err := s.db.GetLoginData(req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Infof("user with username: %s was not found in the database", req.Username)
			return nil, status.Error(
				codes.NotFound,
				err.Error(),
			)
		} else {
			s.logger.Infof("getting data for: %s failed: %s", req.Username, err.Error())
			return nil, status.Error(
				codes.Internal,
				err.Error(),
			)
		}
	}

	err = utils.CompareHashAndPassword(loginData.PasswordHash, req.Password)
	if err != nil {
		s.logger.Infof("password of user: %s is not matching", req.Username)
		return nil, status.Errorf(
			codes.PermissionDenied,
			err.Error(),
		)
	}

	return loginData, nil
}

func (s *UsersServiceServer) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.UsersServiceID, error) {
	var id pb.UsersServiceID

	s.logger.Infof("inserting new user: %s", req.Username)

	//check if email and username exist in the database
	err := s.db.CheckIfUsernameExists(req.Username)
	if err != nil {
		if err.Error() == database.UsernameAlreadyExistsError {
			s.logger.Errorf("username: %s already exists in the database", req.Username)
			return nil, status.Error(
				codes.AlreadyExists,
				err.Error(),
			)
		}

		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	err = s.db.CheckIfEmailExists(req.Email)
	if err != nil {
		if err.Error() == database.EmailAlreadyExistsError {
			s.logger.Errorf("email: %s already exists in the database", req.Email)
			return nil, status.Error(
				codes.AlreadyExists,
				err.Error(),
			)
		}

		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	//validate input data
	err = validator.ValidateLoginData(req.Username, req.Email, req.Password1, req.Password2)
	if err != nil {
		s.logger.Errorf("validating data for: %s failed: %s", req.Username, err.Error())
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	//insert into the database
	idInt, err := s.db.InsertUser(req)
	if err != nil {
		s.logger.Errorf("inserting user: %s failed: %s", req.Username, err.Error())
		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	id.Id = idInt

	return &id, nil
}

func (s *UsersServiceServer) DeleteLoginData(ctx context.Context, req *pb.UsersServiceID) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
