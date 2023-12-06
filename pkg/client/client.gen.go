package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oapi-codegen/runtime"
)

const (
	ApiKeyAuthScopes = "ApiKeyAuth.Scopes"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Help    *string `json:"help,omitempty"`
	Message *string `json:"message,omitempty"`
	Ok      bool    `json:"ok"`
}

// Location defines model for Location.
type Location struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// MeasuringPoint defines model for MeasuringPoint.
type MeasuringPoint struct {
	Active              *bool    `json:"active,omitempty"`
	AlarmPercentage     *int     `json:"alarm_percentage,omitempty"`
	Category            *string  `json:"category,omitempty"`
	DataSaveLevel       *float32 `json:"data_save_level,omitempty"`
	DisableLed          *bool    `json:"disable_led,omitempty"`
	GuideLine           *string  `json:"guide_line,omitempty"`
	Id                  *int     `json:"id,omitempty"`
	MeasurementDuration *int     `json:"measurement_duration,omitempty"`
	MeasuringType       *string  `json:"measuring_type,omitempty"`
	Name                *string  `json:"name,omitempty"`
	TracePostTrigger    *float64 `json:"trace_post_trigger,omitempty"`
	TracePreTrigger     *float64 `json:"trace_pre_trigger,omitempty"`
	TraceSaveLevel      *float64 `json:"trace_save_level,omitempty"`
	// UserLocation        *Location `json:"user_location,omitempty"`
	VibrationType *string `json:"vibration_type,omitempty"`
}

// PeakRecord defines model for PeakRecord.
type PeakRecord struct {
	Category      *string `json:"category,omitempty"`
	GuideLine     *string `json:"guide_line,omitempty"`
	MeasuringType *string `json:"measuring_type,omitempty"`
	Timestamp     *int64  `json:"timestamp,omitempty"`
	VibrationType *string `json:"vibration_type,omitempty"`
}

// PeakRecordsResponse defines model for PeakRecordsResponse.
type PeakRecordsResponse struct {
	Ok      bool         `json:"ok"`
	Samples []PeakRecord `json:"samples"`
}

// Sensor defines model for Sensor.
type Sensor struct {
	ConnectedUsing *string         `json:"connected_using,omitempty"`
	Lastseen       *time.Time      `json:"lastseen,omitempty"`
	Location       *Location       `json:"location,omitempty"`
	MeasuringPoint *MeasuringPoint `json:"measuring_point,omitempty"`
	Name           *string         `json:"name,omitempty"`
}

// SuccessResponse defines model for SuccessResponse.
type SuccessResponse struct {
	Ok      bool      `json:"ok"`
	Sensors *[]Sensor `json:"sensors,omitempty"`
}

// GetPeakRecordsParams defines parameters for GetPeakRecords.
type GetPeakRecordsParams struct {
	// MeasuringPointId ID of the measuring point
	MeasuringPointId int `form:"measuring_point_id" json:"measuring_point_id"`

	// StartTime Start time in milliseconds since epoch
	StartTime int `form:"start_time" json:"start_time"`

	// EndTime End time in milliseconds since epoch (optional)
	EndTime *int `form:"end_time,omitempty" json:"end_time,omitempty"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetPeakRecords request
	GetPeakRecords(ctx context.Context, params *GetPeakRecordsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListSensors request
	ListSensors(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetPeakRecords(ctx context.Context, params *GetPeakRecordsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetPeakRecordsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListSensors(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListSensorsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetPeakRecordsRequest generates requests for GetPeakRecords
func NewGetPeakRecordsRequest(server string, params *GetPeakRecordsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/get_peak_records")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "measuring_point_id", runtime.ParamLocationQuery, params.MeasuringPointId); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "start_time", runtime.ParamLocationQuery, params.StartTime); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		if params.EndTime != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "end_time", runtime.ParamLocationQuery, *params.EndTime); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewListSensorsRequest generates requests for ListSensors
func NewListSensorsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/list_sensors")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, token string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	client.RequestEditors = append(client.RequestEditors,
		func(ctx context.Context, req *http.Request) error {
			// Get the existing query values
			query := req.URL.Query()

			// Add the token to the query values
			query.Set("token", token)

			// Set the modified query values back to the request URL
			req.URL.RawQuery = query.Encode()

			return nil
		})
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetPeakRecordsWithResponse request
	GetPeakRecordsWithResponse(ctx context.Context, params *GetPeakRecordsParams, reqEditors ...RequestEditorFn) (*GetPeakRecordsResponse, error)

	// ListSensorsWithResponse request
	ListSensorsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListSensorsResponse, error)
}

type GetPeakRecordsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PeakRecordsResponse
	JSON400      *ErrorResponse
	JSON500      *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r GetPeakRecordsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetPeakRecordsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListSensorsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SuccessResponse
	JSON400      *ErrorResponse
	JSON500      *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r ListSensorsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListSensorsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetPeakRecordsWithResponse request returning *GetPeakRecordsResponse
func (c *ClientWithResponses) GetPeakRecordsWithResponse(ctx context.Context, params *GetPeakRecordsParams, reqEditors ...RequestEditorFn) (*GetPeakRecordsResponse, error) {
	rsp, err := c.GetPeakRecords(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetPeakRecordsResponse(rsp)
}

// ListSensorsWithResponse request returning *ListSensorsResponse
func (c *ClientWithResponses) ListSensorsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListSensorsResponse, error) {
	rsp, err := c.ListSensors(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListSensorsResponse(rsp)
}

// ParseGetPeakRecordsResponse parses an HTTP response from a GetPeakRecordsWithResponse call
func ParseGetPeakRecordsResponse(rsp *http.Response) (*GetPeakRecordsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetPeakRecordsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PeakRecordsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseListSensorsResponse parses an HTTP response from a ListSensorsWithResponse call
func ParseListSensorsResponse(rsp *http.Response) (*ListSensorsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListSensorsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest SuccessResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}
