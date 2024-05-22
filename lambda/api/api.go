package api

import (
	"encoding/json"
	"lambdafunc/database"
	"lambdafunc/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var registerUser types.RegisterUserDto

	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "request has empty parameters",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := types.NewUser(registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "there was an error hashing the password, %s",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// does user with username already exist

	exists, err := api.dbStore.DoesUserExist(user.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "there was an error checking if user exists, %w",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if exists {
		return events.APIGatewayProxyResponse{
			Body:       "user %s already exists",
			StatusCode: http.StatusNotFound,
		}, err
	}

	// insert user
	err = api.dbStore.InsertUser(types.RegisterUserDto{
		Username: user.Username,
		Password: user.Password,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusNotFound,
		}, err
	}

	responseUser, err := json.Marshal(types.RegisterUserDto{
		Username: user.Username,
		Password: "",
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusNotFound,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseUser),
		StatusCode: http.StatusNotFound,
	}, nil
}

func (api ApiHandler) LoginUserHandler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	type LoginUserDto struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var getUser LoginUserDto
	err := json.Unmarshal([]byte(event.Body), &getUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, err
	}

	user, err := api.dbStore.GetUser(getUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	if !types.ValidatePassword(user.Password, getUser.Username) {
		return events.APIGatewayProxyResponse{
			Body:       "invalid password",
			StatusCode: http.StatusNotFound,
		}, nil
	}

	responseUser, err := json.Marshal(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseUser),
		StatusCode: http.StatusNotFound,
	}, nil
}
