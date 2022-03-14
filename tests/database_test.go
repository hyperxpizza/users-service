package main

import (
	"testing"

	"github.com/hyperxpizza/users-service/pkg/config"
	"github.com/hyperxpizza/users-service/pkg/database"
	pb "github.com/hyperxpizza/users-service/pkg/grpc"
	"github.com/stretchr/testify/assert"
)

// go test -v ./tests/ --run TestDatabaseOperations --config=/home/hyperxpizza/dev/golang/reusable-microservices/users-service/config.json
func TestDatabaseOperations(t *testing.T) {

	getLoginData := func() *pb.RegisterUserRequest {
		return &pb.RegisterUserRequest{
			Username:  "pizza",
			Email:     "pizza@kernelpanic.live",
			Password1: "passworD1#",
			Password2: "passworD1#",
		}
	}

	if *configPathOpt == "" {
		t.Fail()
		return
	}

	cfg, err := config.NewConfig(*configPathOpt)
	assert.NoError(t, err)

	db, err := database.Connect(cfg)
	assert.NoError(t, err)

	defer db.Close()

	id, err := db.InsertUser(getLoginData())
	assert.NoError(t, err)

	if *deleteOpt {
		err := db.DeleteLoginData(id)
		assert.NoError(t, err)
	}

	/*
		t.Run("Insert LoginData", func(t *testing.T) {


		})
	*/
}
