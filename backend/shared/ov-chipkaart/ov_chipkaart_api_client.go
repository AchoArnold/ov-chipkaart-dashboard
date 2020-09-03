package ovchipkaart

import (
	"bytes"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/palantir/stacktrace"

	"github.com/AchoArnold/homework/services/json"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
)

const endpointAuthentication = "https://login.ov-chipkaart.nl/oauth2/token"
const endpointAuthorisation = "https://api2.ov-chipkaart.nl/femobilegateway/v1/api/authorize"
const endpointTransactions = "https://api2.ov-chipkaart.nl/femobilegateway/v1/transactions"

const contentTypeJSON = "application/json"
const contentTypeFormURLEncoded = "application/x-www-form-urlencoded"

const responseCodeOk = 200

const dateFormat = "2006-01-02"

const transactionRequestsPerSecond = 10

const (
	// ErrCodeUnauthorized is returned when the user is not authorized
	ErrCodeUnauthorized = stacktrace.ErrorCode(401)
	// ErrCodeInternalServerError represents any other error
	ErrCodeInternalServerError = stacktrace.ErrorCode(500)
)

type authenticationTokenResponse struct {
	IDToken          string `json:"id_token"`
	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`
}

type authorisationTokenResponse struct {
	ResponseCode int         `json:"c"`
	Value        string      `json:"o"`
	Error        interface{} `json:"e"`
}

type transactionsResponse struct {
	ResponseCode int `json:"c"`
	Response     struct {
		TotalSize              int         `json:"totalSize"`
		NextOffset             int         `json:"nextOffset"`
		PreviousOffset         int         `json:"previousOffset"`
		Records                []RawRecord `json:"records"`
		TransactionsRestricted bool        `json:"transactionsRestricted"`
		NextRequestContext     struct {
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Offset    int    `json:"offset"`
		} `json:"nextRequestContext"`
	} `json:"o"`
	Error interface{} `json:"e"`
}

type transactionsPayload struct {
	AuthorisationToken string `json:"authorizationToken"`
	MediumID           string `json:"mediumId"`
	Locale             string `json:"locale"`
	Offset             string `json:"offset"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
}

// APIClient is responsible for the fetching transactions using the ov-chipkaart API
type APIClient struct {
	clientID     string
	clientSecret string
	httpClient   HTTPClient
	locale       string
}

// APIServiceConfig is the configuration for this service
type APIServiceConfig struct {
	ClientID     string
	ClientSecret string
	Locale       string
	Client       HTTPClient
}

// NewAPIService Initializes the API service.
func NewAPIService(config APIServiceConfig) APIClient {
	return APIClient{
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
		httpClient:   config.Client,
		locale:       config.Locale,
	}
}

// GetAuthorisationToken fetches the auth token based on username/password combination
func (client APIClient) GetAuthorisationToken(username string, password string) (authorisationToken string, err error) {
	authenticationToken, err := client.getAuthenticationToken(username, password)
	if err != nil {
		return authorisationToken, stacktrace.PropagateWithCode(err, ErrCodeUnauthorized, "could not fetch authentication token")
	}

	authorisationTokenResponse, err := client.getAuthorisationToken(authenticationToken)
	if err != nil {
		return authorisationToken, stacktrace.PropagateWithCode(err, ErrCodeUnauthorized, "could not fetch authorisation token")
	}

	return authorisationTokenResponse.Value, err
}

// FetchTransactions returns the transaction records based on the parameter provided.
func (client APIClient) FetchTransactions(options TransactionFetchOptions) (records []RawRecord, err error) {
	authorisationToken, err := client.GetAuthorisationToken(options.Username, options.Password)
	if err != nil {
		return records, stacktrace.PropagateWithCode(err, ErrCodeUnauthorized, "could not authenticate user")
	}

	records, err = client.getTransactions(authorisationToken, options)
	if err != nil {
		return records, stacktrace.PropagateWithCode(err, ErrCodeInternalServerError, "could not fetch transactions")
	}

	return records, nil
}

func (client APIClient) getTransactions(authorisationToken string, options TransactionFetchOptions) ([]RawRecord, error) {
	payload := transactionsPayload{
		AuthorisationToken: authorisationToken,
		MediumID:           options.CardNumber,
		Locale:             client.locale,
		StartDate:          options.StartDate.Format(dateFormat),
		EndDate:            options.EndDate.Format(dateFormat),
	}

	transactionsResponse, err := client.getTransaction(payload)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot perform transactions request: payload = %+v", payload)
	}

	records := transactionsResponse.Response.Records

	payload.StartDate = transactionsResponse.Response.NextRequestContext.StartDate
	payload.EndDate = transactionsResponse.Response.NextRequestContext.EndDate

	requestLimit := len(records)

	numberOfRequests := int(math.Ceil(float64(transactionsResponse.Response.TotalSize) / float64(requestLimit)))

	rateLimiter := ratelimit.New(transactionRequestsPerSecond)
	for i := 1; i < numberOfRequests; i++ {
		payload.Offset = strconv.Itoa(transactionsResponse.Response.NextRequestContext.Offset)

		rateLimiter.Take()

		transactionsResponse, err = client.getTransaction(payload)
		if err != nil {
			return nil, stacktrace.Propagate(err, "cannot perform transactions request: payload = %+v", payload)
		}

		records = append(records, transactionsResponse.Response.Records...)
	}

	return records, nil
}

func (client APIClient) getTransaction(payload transactionsPayload) (transactionsResponse *transactionsResponse, err error) {
	payloadAsMap, err := json.JsonToStringMap(payload)
	if err != nil {
		return transactionsResponse, stacktrace.Propagate(err, "cannot serialize request to map %#+v", payload)
	}

	request, err := client.createPostRequest(endpointTransactions, payloadAsMap)
	if err != nil {
		return transactionsResponse, stacktrace.Propagate(err, "cannot create transaction request: payload = %+#v", payloadAsMap)
	}

	response, err := client.doHTTPRequest(request)
	if err != nil {
		return transactionsResponse, stacktrace.Propagate(err, "cannot perform transaction request: payload = %+#v", request)
	}

	err = json.JsonDecode(&transactionsResponse, response.Body)
	if err != nil {
		return transactionsResponse, stacktrace.Propagate(err, "cannot decode response into transactions response: payload = %+#v", response)
	}

	if transactionsResponse != nil && transactionsResponse.ResponseCode != responseCodeOk {
		return transactionsResponse, stacktrace.NewError("Invalid response code %d: payload = %+#v", transactionsResponse.ResponseCode, payload)
	}

	return transactionsResponse, nil
}

func (client APIClient) getAuthorisationToken(authenticationTokenResponse authenticationTokenResponse) (authorisationToken authorisationTokenResponse, err error) {
	payload := map[string]string{
		"authenticationToken": authenticationTokenResponse.IDToken,
	}

	request, err := client.createPostRequest(endpointAuthorisation, payload)
	if err != nil {
		return authorisationToken, stacktrace.Propagate(err, "cannot create authorisation request")
	}

	response, err := client.doHTTPRequest(request)
	if err != nil {
		return authorisationToken, stacktrace.Propagate(err, "cannot perform authorisation request")
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return authorisationToken, stacktrace.Propagate(err, "cannot read body from response")
	}

	err = json.JsonDecode(&authorisationToken, bytes.NewBuffer(responseBody))
	if err != nil {
		return authorisationToken, stacktrace.Propagate(err, "cannot decode authorisation token response content: %s", responseBody)
	}

	if authorisationToken.ResponseCode != responseCodeOk {
		return authorisationToken, stacktrace.NewError("Response Code: %d, Error: %s", authorisationToken.ResponseCode, authorisationToken.Value)
	}

	return authorisationToken, nil
}

func (client APIClient) getAuthenticationToken(username, password string) (authenticationToken authenticationTokenResponse, err error) {
	payload := map[string]string{
		"username":      username,
		"password":      password,
		"client_id":     client.clientID,
		"client_secret": client.clientSecret,
		"grant_type":    "password",
		"scope":         "openid",
	}

	request, err := client.createPostRequest(endpointAuthentication, payload)
	if err != nil {
		return authenticationToken, stacktrace.Propagate(err, "cannot create authentication request")
	}

	response, err := client.doHTTPRequest(request)
	if err != nil {
		return authenticationToken, stacktrace.Propagate(err, "cannot perform authentication request")
	}

	err = json.JsonDecode(&authenticationToken, response.Body)
	if err != nil {
		return authenticationToken, stacktrace.Propagate(err, "cannot decode authentication token response")
	}

	if authenticationToken.Error != "" {
		return authenticationToken, stacktrace.Propagate(errors.New(authenticationToken.Error), authenticationToken.ErrorDescription)
	}

	return authenticationToken, nil
}

func (client APIClient) doHTTPRequest(request *http.Request) (*http.Response, error) {
	apiResponse, err := client.httpClient.Do(request)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot execute %s request for %s: ", request.Method, request.URL.String())
	}

	return apiResponse, nil
}

func (client APIClient) createPostRequest(endpoint string, payload map[string]string) (*http.Request, error) {
	data := url.Values{}
	for key, val := range payload {
		data.Set(key, val)
	}

	apiRequest, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot create request for URL: "+endpoint)
	}

	apiRequest.Header.Set("Accept", contentTypeJSON)
	apiRequest.Header.Set("Content-Type", contentTypeFormURLEncoded)

	return apiRequest, nil
}
