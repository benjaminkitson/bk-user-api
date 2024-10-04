package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/benjaminkitson/bk-user-api/models/user"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type HTTPClient struct {
	baseURL   *url.URL
	client    *http.Client
	awsConfig aws.Config
	logger    *zap.Logger
}

type ClientError struct {
	Message    string
	StatusCode int
}

func (e ClientError) Error() string {
	return e.Message
}

func NewClient(baseURL string, logger *zap.Logger) (HTTPClient, error) {
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
		logger:    logger,
	}, nil
}

func (c HTTPClient) CreateUser(ctx context.Context, email string) (user.User, error) {
	r := *c.baseURL
	r.Path = path.Join(r.Path, "create")
	bodyMap := map[string]string{"email": email}
	b, err := json.Marshal(bodyMap)
	body := bytes.NewReader(b)
	if err != nil {
		return user.User{}, err
	}
	c.logger.Info("building request")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.String(), body)
	if err != nil {
		return user.User{}, err
	}
	c.logger.Info("sending request")
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("request error")
		return user.User{}, err
	}
	if res.StatusCode == 200 {
		c.logger.Info("request success")
		var u user.User
		err = json.NewDecoder(res.Body).Decode(&u)
		if err != nil {
			return user.User{}, err
		}
		return u, err
	}
	bodyRes, _ := io.ReadAll(res.Body)
	if res.StatusCode == 400 || res.StatusCode == 500 {
		c.logger.Error("error status code received", zap.Int("statusCode", res.StatusCode))
		return user.User{}, ClientError{StatusCode: res.StatusCode, Message: string(bodyRes)}
	}
	err = fmt.Errorf("api responded with unexpected status code %d, with body %s", res.StatusCode, string(bodyRes))
	return user.User{}, err
}
