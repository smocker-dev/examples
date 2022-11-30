package sdks

import (
	"context"
	"encoding/json"
	"net/http"

	"example/types"
)

type Users interface {
	GetUsers(ctx context.Context) ([]types.User, error)
	GetUserByName(ctx context.Context, name string) (types.User, error)
}

type usersClient struct {
	url string
}

func NewUsersClient(url string) Users {
	return &usersClient{url: url}
}

func (c *usersClient) GetUsers(ctx context.Context) ([]types.User, error) {
	var users []types.User
	req, err := http.NewRequestWithContext(ctx, "GET", c.url+"/users", nil)
	if err != nil {
		return users, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return users, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&users)
	return users, err
}

func (c *usersClient) GetUserByName(ctx context.Context, name string) (types.User, error) {
	var user types.User
	req, err := http.NewRequestWithContext(ctx, "GET", c.url+"/users/"+name, nil)
	if err != nil {
		return user, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&user)
	return user, err
}
