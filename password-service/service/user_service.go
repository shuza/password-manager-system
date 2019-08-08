package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"password-service/error_tracer"
	"password-service/model"
)

type UserService struct{}

func (s *UserService) GetUserId(token string) (int, error) {
	url := fmt.Sprintf("%s/api/v1/auth/token", os.Getenv("TOKEN_SERVICE_HOST"))
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	query := req.URL.Query()
	query.Add("token", token)
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		error_tracer.Client.ErrorLog("userService", "httpCall", fmt.Sprintf("%s http call Error : %v", url, err))
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		error_tracer.Client.InfoLog("userService", "httpCall", "invalid token response")
		return 0, errors.New("Invalid token")
	}

	var response struct {
		Message string     `json:"message"`
		Data    model.User `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		error_tracer.Client.ErrorLog("userService", "responseBody", fmt.Sprintf("%s response Error :  %v", url, err))
		return 0, err
	}

	return response.Data.Id, nil
}
