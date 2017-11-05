package account

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// AddCard adds the card via stripe token.
func (c *Client) AddCard(token, cardToken string) error {
	in := struct {
		Token string `json:"token"`
	}{
		Token: cardToken,
	}

	b, err := json.Marshal(in)
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	url := fmt.Sprintf("%s/billing/cards", c.url)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "requesting")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return request.Error{Status: res.StatusCode}
	}

	return nil
}

// GetCards returns the user's cards.
func (c *Client) GetCards(token string) (cards []Card, err error) {
	req, err := http.NewRequest("GET", c.url+"/billing/cards", nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, request.Error{Status: res.StatusCode}
	}

	err = json.NewDecoder(res.Body).Decode(&cards)
	return
}

// GetPlans returns the user's plan(s).
func (c *Client) GetPlans(token string) (plans []Plan, err error) {
	req, err := http.NewRequest("GET", c.url+"/billing/plans", nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, request.Error{Status: res.StatusCode}
	}

	err = json.NewDecoder(res.Body).Decode(&plans)
	return
}

// RemoveCard removes a user's card by id.
func (c *Client) RemoveCard(token, id string) error {
	req, err := http.NewRequest("DELETE", c.url+"/billing/cards/"+id, nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return request.Error{Status: res.StatusCode}
	}

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

	b, err := json.Marshal(in)
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	req, err := http.NewRequest("PUT", c.url+"/billing/plans", bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "creating request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return request.Error{Status: res.StatusCode}
	}

	return nil
}

// RemovePlan unsubscribes from a plan.
func (c *Client) RemovePlan(token, product, plan string) error {
	url := fmt.Sprintf("%s/billing/plans/%s/%s", c.url, product, plan)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return request.Error{Status: res.StatusCode}
	}

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

	b, err := json.Marshal(in)
	if err != nil {
		return "", errors.Wrap(err, "marshaling")
	}

	res, err := http.Post(c.url+"/login", "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", request.Error{Status: res.StatusCode}
	}

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

	b, err := json.Marshal(in)
	if err != nil {
		return "", errors.Wrap(err, "marshaling")
	}

	res, err := http.Post(c.url+"/access_key", "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", request.Error{Status: res.StatusCode}
	}

	b, err = ioutil.ReadAll(res.Body)
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

			if err, ok := err.(request.Error); ok && err.Status == http.StatusUnauthorized {
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
