package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hyperxpizza/users-service/pkg/config"
	pb "github.com/hyperxpizza/users-service/pkg/grpc"
	"github.com/hyperxpizza/users-service/pkg/utils"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	UsernameAlreadyExistsError = "username already exists in the database"
	EmailAlreadyExistsError    = "email already exists in the database"
)

var ctx = context.Background()

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

func (db *Database) InsertUser(data *pb.RegisterUserRequest) (int64, error) {
	var id int64

	passwordHash, err := utils.GenerateHashAndSalt(data.Password1)
	if err != nil {
		return 0, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil
	}

	stmt, err := tx.Prepare(`
	insert into logins(id, username, email, passwordHash, created, updated)
	values(default, $1, $2, $3, $4, $5)
	returning id`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var lID int
	err = stmt.QueryRow(data.Username, data.Email, passwordHash, time.Now(), time.Now()).Scan(&lID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmt2, err := tx.Prepare(`insert into users(id, created, updated, loginid) values (default, $1, $2, $3) returning id`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = stmt2.QueryRow(time.Now(), time.Now(), lID).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

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

func (db *Database) DeleteLoginData(id int64) error {
	stmt, err := db.Prepare(`delete from users where id=$1`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) CheckIfUsernameExists(username string) error {

	err := db.QueryRow(`select username from logins where username=$1`, username).Scan(&username)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil
		}

		return err
	}

	return errors.New(UsernameAlreadyExistsError)
}

func (db *Database) CheckIfEmailExists(email string) error {
	err := db.QueryRow(`select email from logins where email=$1`, email).Scan(&email)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil
		}

		return err
	}

	return errors.New(EmailAlreadyExistsError)
}
