package userapiclient

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/benjaminkitson/bk-user-api/models"
	"go.uber.org/zap"
)

type HTTPClient struct {
	baseURL       *url.URL
	client        *http.Client
	awsConfig     aws.Config
	logger        *zap.Logger
	requestSigner *v4.Signer
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
	s := v4.NewSigner()
	return HTTPClient{
		baseURL:       u,
		client:        c,
		awsConfig:     cfg,
		logger:        logger,
		requestSigner: s,
	}, nil
}

func (c HTTPClient) CreateUser(ctx context.Context, email string) (models.User, error) {
	r := *c.baseURL
	r.Path = path.Join(r.Path, "create")
	bodyMap := map[string]string{"email": email}
	b, err := json.Marshal(bodyMap)
	body := bytes.NewReader(b)
	if err != nil {
		return models.User{}, err
	}

	c.logger.Info("building request")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.String(), body)
	if err != nil {
		return models.User{}, err
	}

	creds, err := c.awsConfig.Credentials.Retrieve(ctx)
	if err != nil {
		return models.User{}, err
	}

	h := sha256.Sum256([]byte(b))
	s := hex.EncodeToString(h[:])
	c.requestSigner.SignHTTP(ctx, creds, req, s, "execute-api", "eu-west-2", time.Now())

	c.logger.Info("sending request")
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("request error")
		return models.User{}, err
	}
	if res.StatusCode == 200 {
		c.logger.Info("request success")
		var u models.User
		err = json.NewDecoder(res.Body).Decode(&u)
		if err != nil {
			return models.User{}, err
		}
		return u, err
	}
	bodyRes, _ := io.ReadAll(res.Body)
	if res.StatusCode == 400 || res.StatusCode == 500 {
		c.logger.Error("error status code received", zap.Int("statusCode", res.StatusCode))
		return models.User{}, ClientError{StatusCode: res.StatusCode, Message: string(bodyRes)}
	}
	err = fmt.Errorf("api responded with unexpected status code %d, with body %s", res.StatusCode, string(bodyRes))
	return models.User{}, err
}

// TODO: convert this to use the DELETE method
func (c HTTPClient) DeleteUser(ctx context.Context, id string) (models.User, error) {
	r := *c.baseURL
	r.Path = path.Join(r.Path, "delete")
	bodyMap := map[string]string{"id": id}
	b, err := json.Marshal(bodyMap)
	body := bytes.NewReader(b)
	if err != nil {
		return models.User{}, err
	}

	c.logger.Info("building request")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.String(), body)
	if err != nil {
		return models.User{}, err
	}

	creds, err := c.awsConfig.Credentials.Retrieve(ctx)
	if err != nil {
		return models.User{}, err
	}

	h := sha256.Sum256([]byte(b))
	s := hex.EncodeToString(h[:])
	c.requestSigner.SignHTTP(ctx, creds, req, s, "execute-api", "eu-west-2", time.Now())

	c.logger.Info("sending request")
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("request error")
		return models.User{}, err
	}
	if res.StatusCode == 200 {
		c.logger.Info("request success")
		var u models.User
		err = json.NewDecoder(res.Body).Decode(&u)
		if err != nil {
			return models.User{}, err
		}
		return u, err
	}
	bodyRes, _ := io.ReadAll(res.Body)
	if res.StatusCode == 400 || res.StatusCode == 500 {
		c.logger.Error("error status code received", zap.Int("statusCode", res.StatusCode))
		return models.User{}, ClientError{StatusCode: res.StatusCode, Message: string(bodyRes)}
	}
	err = fmt.Errorf("api responded with unexpected status code %d, with body %s", res.StatusCode, string(bodyRes))
	return models.User{}, err
}
