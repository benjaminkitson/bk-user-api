package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/benjaminkitson/bk-user-api/models/user"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type HTTPClient struct {
	baseURL   *url.URL
	client    *http.Client
	awsConfig aws.Config
}

type ClientError struct {
	Message    string
	StatusCode int
}

func (e ClientError) Error() string {
	return e.Message
}

func NewClient(baseURL string) (HTTPClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return HTTPClient{}, err
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return HTTPClient{}, err
	}
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	return HTTPClient{
		baseURL:   u,
		client:    c,
		awsConfig: cfg,
	}, nil
}

func (c HTTPClient) CreateUser(ctx context.Context, username string) (user.User, error) {
	r := *c.baseURL
	r.Path = path.Join(r.Path, "create")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.String(), nil)
	if err != nil {
		return user.User{}, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return user.User{}, err
	}
	if res.StatusCode == 200 {
		var u user.User
		err = json.NewDecoder(res.Body).Decode(&u)
		if err != nil {
			return user.User{}, err
		}
		return u, err
	}
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode == 400 || res.StatusCode == 500 {
		return user.User{}, ClientError{StatusCode: res.StatusCode, Message: string(body)}
	}
	err = fmt.Errorf("api responded with unexpected status code %d, with body %s", res.StatusCode, string(body))
	return user.User{}, err
}
