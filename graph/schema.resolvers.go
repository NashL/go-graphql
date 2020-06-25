package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"strconv"
	"time"

	dbConnection "github.com/nashl/online-store-server/database"
	"github.com/nashl/online-store-server/graph/generated"
	"github.com/nashl/online-store-server/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	now := int(time.Now().Unix())
	user := &model.User{
		FullName:  input.FullName,
		Email:     input.Email,
		Password:  input.Password,
		CreatedAt: now,
		UpdatedAt: 0,
	}
	result, err := dbConnection.DB.Exec("INSERT INTO `users` (fullName, email, password, createdAt, updatedAt) VALUES(?, ?, ?, ?, ?)",
		user.FullName, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	fmt.Println("after the first error handler")

	lastId, err := result.LastInsertId()
	fmt.Println("after getting lastId:", lastId)
	if err != nil {
		return nil, err
	}
	fmt.Println("after second error handler")
	user.UserID = strconv.Itoa(int(lastId))
	fmt.Println("after setting userId")
	return user, nil
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.Token, error) {
	var u model.User

	err := dbConnection.DB.QueryRow("SELECT userId, email FROM `users` WHERE email = ? AND password = ? ", email, password).Scan(&u.UserID, &u.Email)


	if err != nil{
		if err == sql.ErrNoRows {
			return nil, errors.New("user or password wrong")
		} else {
			log.Fatal(err)
		}
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"userId": u.UserID,
		"exp": int(time.Now().Add(time.Hour * 1).Unix()),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		log.Fatal(err)
	}

	expiredAt := time.Now().Add(time.Hour * 1).Unix()
	obj := &model.Token{
		Token:     tokenString,
		ExpiredAt: int(expiredAt),
	}

	return obj, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	var result []*model.User
	rows, err := dbConnection.DB.Query("SELECT userId, fullname, email, password, createdAt, updatedAt FROM `users`")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // important

	for rows.Next() {
		var u model.User
		err = rows.Scan(&u.UserID, &u.FullName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, &u)
	}
	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
