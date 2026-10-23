// Package bbb is a small BigBlueButton API client. Every call is a GET whose
// query string is signed with sha256(call + query + secret).
package bbb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	requestTimeout  = 15 * time.Second
	maxResponseSize = 1 << 20

	returnCodeSuccess = "SUCCESS"

	callCreate    = "create"
	callIsRunning = "isMeetingRunning"
	callJoin      = "join"

	apiPath           = "api/"
	checksumSeparator = "&checksum="

	// Blindside Networks answers blockedApplication to Go's default
	// User-Agent, so every call names this service instead.
	headerUserAgent = "User-Agent"
	userAgent       = "irma-meet (+https://github.com/privacybydesign/irma-meet)"

	paramMeetingID   = "meetingID"
	paramName        = "name"
	paramAttendeePW  = "attendeePW"
	paramModeratorPW = "moderatorPW"
	paramFullName    = "fullName"
	paramPassword    = "password"
	paramRedirect    = "redirect"
)

// Room is a meeting as BigBlueButton knows it.
type Room struct {
	ID                string
	Name              string
	ModeratorPassword string
	AttendeePassword  string
}

type Client struct {
	baseURL string
	secret  string
	client  *http.Client
}

// New returns a client for the server at baseURL, which ends in a slash and
// points at the bigbluebutton/ root.
func New(baseURL, secret string, client *http.Client) *Client {
	if client == nil {
		client = &http.Client{Timeout: requestTimeout}
	}

	return &Client{baseURL: baseURL, secret: secret, client: client}
}

type response struct {
	ReturnCode string `xml:"returncode"`
	MessageKey string `xml:"messageKey"`
	Running    bool   `xml:"running"`
}

// EnsureRoom creates the room unless it is already running. It is safe to
// call on every join.
func (c *Client) EnsureRoom(ctx context.Context, room Room) error {
	running, err := c.call(ctx, callIsRunning, url.Values{paramMeetingID: {room.ID}})
	if err == nil && running.Running {
		return nil
	}

	_, err = c.call(ctx, callCreate, url.Values{
		paramMeetingID:   {room.ID},
		paramName:        {room.Name},
		paramAttendeePW:  {room.AttendeePassword},
		paramModeratorPW: {room.ModeratorPassword},
	})

	return err
}

// JoinURL returns the signed URL that puts fullName in the room, where
// BigBlueButton then redirects the browser to the client.
func (c *Client) JoinURL(roomID, fullName, password string) string {
	return c.signedURL(callJoin, url.Values{
		paramMeetingID: {roomID},
		paramFullName:  {fullName},
		paramPassword:  {password},
		paramRedirect:  {"true"},
	})
}

func (c *Client) call(ctx context.Context, name string, params url.Values) (response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.signedURL(name, params), nil)
	if err != nil {
		return response{}, err
	}

	req.Header.Set(headerUserAgent, userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		// The *url.Error quotes the request URL, which holds the room passwords.
		var uerr *url.Error
		if errors.As(err, &uerr) {
			err = uerr.Err
		}

		return response{}, fmt.Errorf("bbb %s: %w", name, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return response{}, fmt.Errorf("bbb %s: %w", name, err)
	}

	var parsed response
	if err := xml.Unmarshal(data, &parsed); err != nil {
		return response{}, fmt.Errorf("bbb %s: status %d, invalid response: %w", name, resp.StatusCode, err)
	}

	if parsed.ReturnCode != returnCodeSuccess {
		return parsed, errors.New("bbb " + name + " failed: " + parsed.MessageKey)
	}

	return parsed, nil
}

func (c *Client) signedURL(name string, params url.Values) string {
	query := params.Encode()
	sum := sha256.Sum256([]byte(name + query + c.secret))

	return c.baseURL + apiPath + name + "?" + query + checksumSeparator + hex.EncodeToString(sum[:])
}
