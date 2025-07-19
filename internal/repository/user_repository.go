package repository

import (
	"context"
	"socialNetworkOtus/internal/api"

	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type UserRepository struct {
	db *goqu.Database
}

func NewUserRepository(db *goqu.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, req *api.PostUserRegisterJSONBody, passwordHash string) (string, error) {
	var id string
	var birthDate interface{} = nil
	if req.Birthdate != nil {
		birthDate = req.Birthdate.Time
	}
	insert := r.db.Insert("users").Rows(goqu.Record{
		"first_name":    req.FirstName,
		"second_name":   req.SecondName,
		"birth_date":    birthDate,
		"biography":     req.Biography,
		"city":          req.City,
		"password_hash": passwordHash,
	}).Returning("id")
	_, err := insert.Executor().ScanValContext(ctx, &id)
	return id, err
}

type userDB struct {
	Id         string     `db:"id"`
	FirstName  string     `db:"first_name"`
	SecondName string     `db:"second_name"`
	Birthdate  *time.Time `db:"birth_date"`
	Biography  *string    `db:"biography"`
	City       *string    `db:"city"`
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*api.User, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	var u userDB
	ds := r.db.From("users").Select(
		goqu.I("id"),
		goqu.I("first_name"),
		goqu.I("second_name"),
		goqu.I("birth_date"),
		goqu.I("biography"),
		goqu.I("city"),
	).Where(goqu.Ex{"id": id})
	sqlStr, args, _ := ds.ToSQL()
	fmt.Printf("SQL: %s, args: %v\n", sqlStr, args)
	found, err := ds.ScanStructContext(ctx, &u)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, sql.ErrNoRows
	}
	user := &api.User{
		Id:         &u.Id,
		FirstName:  &u.FirstName,
		SecondName: &u.SecondName,
		Biography:  u.Biography,
		City:       u.City,
	}
	if u.Birthdate != nil {
		bd := openapi_types.Date{Time: *u.Birthdate}
		user.Birthdate = &bd
	}
	return user, nil
}

func (r *UserRepository) SearchUsersByPrefix(ctx context.Context, firstName, lastName string) ([]api.User, error) {
	var users []userDB
	ds := r.db.From("users").Select(
		goqu.I("id"),
		goqu.I("first_name"),
		goqu.I("second_name"),
		goqu.I("birth_date"),
		goqu.I("biography"),
		goqu.I("city"),
	).Where(
		goqu.L("first_name LIKE ?", firstName+"%"),
		goqu.L("second_name LIKE ?", lastName+"%"),
	).Order(goqu.I("id").Asc())

	sqlStr, args, _ := ds.ToSQL()
	fmt.Printf("SQL: %s, args: %v\n", sqlStr, args)

	err := ds.ScanStructsContext(ctx, &users)
	if err != nil {
		return nil, err
	}

	result := make([]api.User, 0, len(users))
	for _, u := range users {
		user := api.User{
			Id:         &u.Id,
			FirstName:  &u.FirstName,
			SecondName: &u.SecondName,
			Biography:  u.Biography,
			City:       u.City,
		}
		if u.Birthdate != nil {
			bd := openapi_types.Date{Time: *u.Birthdate}
			user.Birthdate = &bd
		}
		result = append(result, user)
	}
	return result, nil
}

func (r *UserRepository) DB() *goqu.Database {
	return r.db
}
