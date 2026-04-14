package clients

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrDoctorNotFound = errors.New("doctor not found")
var ErrDoctorServiceUnavailable = errors.New("doctor service unavailable")

type DoctorClient interface {
	DoctorExists(doctorID string) error
}

type httpDoctorClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPDoctorClient(baseURL string) DoctorClient {
	return &httpDoctorClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *httpDoctorClient) DoctorExists(doctorID string) error {
	url := fmt.Sprintf("%s/doctors/%s", c.baseURL, doctorID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return ErrDoctorServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrDoctorNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return ErrDoctorServiceUnavailable
	}
	return nil
}
