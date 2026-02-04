package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	BaseURL        = "https://api.irail.be/v1"
	DefaultTimeout = 60 * time.Second
	UserAgent      = "irail-cli (https://github.com/dedene/irail-cli)"
)

// Client is the iRail API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	lang       string
}

// NewClient creates a new iRail API client.
func NewClient(lang string) *Client {
	return &Client{
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout:   DefaultTimeout,
			Transport: NewRetryTransport(nil),
		},
		lang: lang,
	}
}

// SetLanguage sets the language for API requests.
func (c *Client) SetLanguage(lang string) {
	c.lang = lang
}

// get performs a GET request to the API.
func (c *Client) get(ctx context.Context, endpoint string, params url.Values) (*http.Response, error) {
	if params == nil {
		params = url.Values{}
	}

	params.Set("format", "json")
	if c.lang != "" {
		params.Set("lang", c.lang)
	}

	reqURL := fmt.Sprintf("%s/%s/?%s", c.baseURL, endpoint, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

// handleResponse checks for errors and decodes the response.
func handleResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &NotFoundError{Resource: "resource"}
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := 0
		if v := resp.Header.Get("Retry-After"); v != "" {
			retryAfter, _ = strconv.Atoi(v)
		}

		return nil, &RateLimitError{RetryAfter: retryAfter}
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

// GetStations fetches all stations.
func (c *Client) GetStations(ctx context.Context) (*StationsResponse, error) {
	resp, err := c.get(ctx, "stations", nil)
	if err != nil {
		return nil, err
	}

	return handleResponse[StationsResponse](resp)
}

// GetLiveboard fetches departures or arrivals for a station.
func (c *Client) GetLiveboard(ctx context.Context, station string, arrivals bool, date, time string) (*LiveboardResponse, error) {
	params := url.Values{}
	params.Set("station", station)

	if arrivals {
		params.Set("arrdep", "arrival")
	} else {
		params.Set("arrdep", "departure")
	}

	if date != "" {
		params.Set("date", date)
	}

	if time != "" {
		params.Set("time", time)
	}

	resp, err := c.get(ctx, "liveboard", params)
	if err != nil {
		return nil, err
	}

	result, err := handleResponse[LiveboardResponse](resp)
	if err != nil {
		var notFound *NotFoundError
		if errors.As(err, &notFound) {
			notFound.Resource = "station"
			notFound.ID = station

			return nil, notFound
		}

		return nil, err
	}

	return result, nil
}

// GetConnections fetches connections between two stations.
func (c *Client) GetConnections(ctx context.Context, from, to string, date, time string, arriveBy bool, results int) (*ConnectionsResponse, error) {
	params := url.Values{}
	params.Set("from", from)
	params.Set("to", to)

	if date != "" {
		params.Set("date", date)
	}

	if time != "" {
		params.Set("time", time)
	}

	if arriveBy {
		params.Set("timesel", "arrival")
	} else {
		params.Set("timesel", "departure")
	}

	if results > 0 {
		params.Set("results", strconv.Itoa(results))
	}

	resp, err := c.get(ctx, "connections", params)
	if err != nil {
		return nil, err
	}

	return handleResponse[ConnectionsResponse](resp)
}

// GetVehicle fetches information about a specific vehicle/train.
func (c *Client) GetVehicle(ctx context.Context, id string, date string) (*VehicleResponse, error) {
	params := url.Values{}
	params.Set("id", id)

	if date != "" {
		params.Set("date", date)
	}

	resp, err := c.get(ctx, "vehicle", params)
	if err != nil {
		return nil, err
	}

	result, err := handleResponse[VehicleResponse](resp)
	if err != nil {
		var notFound *NotFoundError
		if errors.As(err, &notFound) {
			notFound.Resource = "vehicle"
			notFound.ID = id

			return nil, notFound
		}

		return nil, err
	}

	return result, nil
}

// GetComposition fetches the composition of a train.
func (c *Client) GetComposition(ctx context.Context, id string) (*CompositionResponse, error) {
	params := url.Values{}
	params.Set("id", id)

	resp, err := c.get(ctx, "composition", params)
	if err != nil {
		return nil, err
	}

	result, err := handleResponse[CompositionResponse](resp)
	if err != nil {
		var notFound *NotFoundError
		if errors.As(err, &notFound) {
			notFound.Resource = "composition"
			notFound.ID = id

			return nil, notFound
		}

		return nil, err
	}

	return result, nil
}

// GetDisturbances fetches current service disturbances.
func (c *Client) GetDisturbances(ctx context.Context) (*DisturbancesResponse, error) {
	resp, err := c.get(ctx, "disturbances", nil)
	if err != nil {
		return nil, err
	}

	return handleResponse[DisturbancesResponse](resp)
}
