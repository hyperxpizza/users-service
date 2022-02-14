package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hyperxpizza/users-service/pkg/config"
	pb "github.com/hyperxpizza/users-service/pkg/grpc"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Database struct {
	*sql.DB
}

func Connect(cfg *config.Config) (*Database, error) {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{database}, nil
}

func (db *Database) InsertLogin() {}

func (db *Database) GetLoginData(username string) (*pb.LoginData, error) {
	var data pb.LoginData

	var created time.Time
	var updated time.Time
	err := db.QueryRow(`select * from logins where username=$1`).Scan(
		&data.Id,
		&data.Username,
		&data.Email,
		&data.PasswordHash,
		&created,
		&updated,
		&data.UserID,
	)
	if err != nil {
		return nil, err
	}

	data.Created = timestamppb.New(created)
	data.Updated = timestamppb.New(updated)

	return &data, nil
}
