package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/imroc/req/v3"

	"github.com/franciscosbf/spotify-tui/internals/uri"
)

var (
	playerEndpoint   string
	devicesEndpoint  string
	playEndpoint     string
	profileEndpoint  string
	pauseEndpoint    string
	previousEndpoint string
	nextEndpoint     string
	shuffleEndpoint  string
	repeatEndpoint   string
	volumeEndpoint   string
)

func endpoint(base, path string) string {
	e, _ := url.JoinPath(base, path)

	return e
}

func init() {
	profileEndpoint = endpoint(uri.API, "v1/me")

	playerEndpoint = endpoint(uri.API, "v1/me/player")

	devicesEndpoint = endpoint(playerEndpoint, "devices")
	playEndpoint = endpoint(playerEndpoint, "play")
	pauseEndpoint = endpoint(playerEndpoint, "pause")
	previousEndpoint = endpoint(playerEndpoint, "previous")
	nextEndpoint = endpoint(playerEndpoint, "next")
	shuffleEndpoint = endpoint(playerEndpoint, "shuffle")
	repeatEndpoint = endpoint(playerEndpoint, "repeat")
	volumeEndpoint = endpoint(playerEndpoint, "volume")
}

var ErrRequestFailed = errors.New("failed to send request")

type Client struct {
	*req.Client
	token string
}

func NewClient(token string) *Client {
	client := req.C()

	return &Client{client, token}
}

func (c *Client) tokenBearer() string {
	return fmt.Sprintf("Bearer %s", c.token)
}

func (c *Client) request(method string, endpoint string, request *req.Request) error {
	bearer := c.tokenBearer()

	var er struct {
		Error ErrResponse `json:"error"`
	}

	request = request.
		SetHeader("Authorization", bearer).
		SetErrorResult(&er)

	resp, err := request.Send(method, endpoint)
	if err != nil {
		return ErrRequestFailed
	}

	if resp.IsErrorState() {
		return er.Error
	}

	return nil
}

func (c *Client) simpleRequest(method, endpoint string) error {
	request := c.R()

	if err := c.request(method, endpoint, request); err != nil {
		return err
	}

	return nil
}

func (c *Client) stateRequest(endpoint, state string) error {
	request := c.R().
		SetQueryParam("state", state)

	if err := c.request(http.MethodPut, endpoint, request); err != nil {
		return err
	}

	return nil
}

func (c *Client) shuffleRequest(shuffle bool) error {
	return c.stateRequest(shuffleEndpoint, strconv.FormatBool(shuffle))
}

func (c *Client) repeatRequest(mode string) error {
	return c.stateRequest(repeatEndpoint, mode)
}

func (c *Client) GetUserProfile() (UserProfile, error) {
	var profile UserProfile

	request := c.R().
		SetSuccessResult(&profile)

	if err := c.request(http.MethodGet, profileEndpoint, request); err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}

func (c *Client) Pause() error {
	return c.simpleRequest(http.MethodPut, pauseEndpoint)
}

func (c *Client) SkipToPrevious() error {
	return c.simpleRequest(http.MethodPost, previousEndpoint)
}

func (c *Client) SkipToNext() error {
	return c.simpleRequest(http.MethodPost, nextEndpoint)
}

func (c *Client) EnableShuffle() error {
	return c.shuffleRequest(true)
}

func (c *Client) DisableShuffle() error {
	return c.shuffleRequest(false)
}

func (c *Client) SetRepeatTrack() error {
	return c.repeatRequest("track")
}

func (c *Client) SetRepeatContext() error {
	return c.repeatRequest("context")
}

func (c *Client) DisableRepeat() error {
	return c.repeatRequest("off")
}

func (c *Client) SetToken(token string) {
	c.token = token
}
