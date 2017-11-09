package account

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tj/go/http/request"
)

// Card model.
type Card struct {
	ID       string `json:"id"`
	Brand    string `json:"brand"`
	LastFour string `json:"last_four"`
}

// CouponDuration is the coupon duration.
type CouponDuration string

// Durations.
const (
	Forever   CouponDuration = "forever"
	Once                     = "once"
	Repeating                = "repeating"
)

// Coupon model.
type Coupon struct {
	ID             string         `json:"id"`
	Amount         int            `json:"amount"`
	Percent        int            `json:"percent"`
	Duration       CouponDuration `json:"duration"`
	DurationPeriod int            `json:"duration_period"`
}

// Discount returns the final price from the given amount.
func (c *Coupon) Discount(n int) int {
	if c.Amount != 0 {
		return n - c.Amount
	}

	return n - int(float64(n)*(float64(c.Percent)/100))
}

// Description returns a humanized description of the savings.
func (c *Coupon) Description() (s string) {
	switch {
	case c.Amount != 0:
		n := fmt.Sprintf("%0.2f", float64(c.Amount)/100)
		s += fmt.Sprintf("$%s off", strings.Replace(n, ".00", "", 1))
	case c.Percent != 0:
		s += fmt.Sprintf("%d%% off", c.Percent)
	}

	switch c.Duration {
	case Repeating:
		s += fmt.Sprintf(" for %d months", c.DurationPeriod)
	default:
		s += fmt.Sprintf(" %s", c.Duration)
	}

	return s
}

// Plan model.
type Plan struct {
	ID         string    `json:"id"`
	Product    string    `json:"product"`
	Plan       string    `json:"plan"`
	PlanName   string    `json:"plan_name"`
	Amount     int       `json:"amount"`
	Interval   string    `json:"interval"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	CanceledAt time.Time `json:"canceled_at"`
}

// Client implementation.
type Client struct {
	url string
}

// New client.
func New(url string) *Client {
	return &Client{
		url: url,
	}
}

// GetCoupon by id.
func (c *Client) GetCoupon(id string) (coupon *Coupon, err error) {
	res, err := c.request("", "GET", "/billing/coupons/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	coupon = new(Coupon)
	err = json.NewDecoder(res.Body).Decode(coupon)
	return
}

// AddCard adds the card via stripe token.
func (c *Client) AddCard(token, cardToken string) error {
	in := struct {
		Token string `json:"token"`
	}{
		Token: cardToken,
	}

	res, err := c.requestJSON(token, "POST", "/billing/cards", in)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// GetCards returns the user's cards.
func (c *Client) GetCards(token string) (cards []Card, err error) {
	res, err := c.request(token, "GET", "/billing/cards", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&cards)
	return
}

// GetPlans returns the user's plan(s).
func (c *Client) GetPlans(token string) (plans []Plan, err error) {
	res, err := c.request(token, "GET", "/billing/plans", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&plans)
	return
}

// RemoveCard removes a user's card by id.
func (c *Client) RemoveCard(token, id string) error {
	res, err := c.request(token, "DELETE", "/billing/cards/"+id, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// AddPlan subscribes to plan.
func (c *Client) AddPlan(token, product, plan, coupon string) error {
	in := struct {
		Product string `json:"product"`
		Plan    string `json:"plan"`
		Coupon  string `json:"coupon"`
	}{
		Product: product,
		Plan:    plan,
		Coupon:  coupon,
	}

	res, err := c.requestJSON(token, "PUT", "/billing/plans", in)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// RemovePlan unsubscribes from a plan.
func (c *Client) RemovePlan(token, product, plan string) error {
	path := fmt.Sprintf("/billing/plans/%s/%s", product, plan)
	res, err := c.request(token, "DELETE", path, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// Login signs in with the given email by
// sending a verification email and returning
// a code which can be exchanged for an access key.
func (c *Client) Login(email string) (code string, err error) {
	in := struct {
		Email string `json:"email"`
	}{
		Email: email,
	}

	res, err := c.requestJSON("", "POST", "/login", in)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var out struct {
		Code string `json:"code"`
	}

	err = json.NewDecoder(res.Body).Decode(&out)
	code = out.Code
	return
}

// GetAccessKey with the given email and code.
func (c *Client) GetAccessKey(email, code string) (key string, err error) {
	in := struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}{
		Email: email,
		Code:  code,
	}

	res, err := c.requestJSON("", "POST", "/access_key", in)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// PollAccessKey polls for an access key.
func (c *Client) PollAccessKey(ctx context.Context, email, code string) (key string, err error) {
	keyC := make(chan string, 1)
	errC := make(chan error, 1)

	go func() {
		for range time.Tick(5 * time.Second) {
			key, err = c.GetAccessKey(email, code)

			if err, ok := err.(request.Error); ok && err == http.StatusUnauthorized {
				continue
			}

			if err != nil {
				errC <- err
				return
			}

			keyC <- key
		}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case e := <-errC:
		return "", e
	case k := <-keyC:
		return k, nil
	}
}

// requestJSON helper.
func (c *Client) requestJSON(token, method, path string, v interface{}) (*http.Response, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling")
	}

	return c.request(token, method, path, bytes.NewReader(b))
}

// request helper.
func (c *Client) request(token, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.url+path, body)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "requesting")
	}

	if res.StatusCode >= 400 {
		res.Body.Close()
		return nil, request.Error(res.StatusCode)
	}

	return res, nil
}
