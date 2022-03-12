package main

import (
	"testing"

	"github.com/hyperxpizza/users-service/pkg/config"
	"github.com/stretchr/testify/assert"
)

// go test -v ./tests/ --run TestLoadConfig --config=/home/hyperxpizza/dev/golang/reusable-microservices/users-service/config.json
func TestLoadConfig(t *testing.T) {
	if *configPathOpt == "" {
		t.Fail()
		return
	}

	cfg, err := config.NewConfig(*configPathOpt)
	assert.NoError(t, err)

	cfg.PrettyPrint()
}
