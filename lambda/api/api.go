package api

import (
	"fmt"
	"lambdafunc/database"
	"lambdafunc/types"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("request has empty parameters")
	}

	user, err := types.NewUser(event)
	if err != nil {
		return fmt.Errorf("there was an error hashing the password, %s", user)
	}

	// does user with username already exist

	exists, err := api.dbStore.DoesUserExist(event.Username)
	if err != nil {
		return fmt.Errorf("there was an error checking if user exists, %w", err)
	}

	if exists {
		return fmt.Errorf("user %s already exists", event.Username)
	}

	// insert user

	err = api.dbStore.InsertUser(event)

	if err != nil {
		return fmt.Errorf("there was an error inserting the user, %w", err)
	}

	return nil
}
