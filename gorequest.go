package go_httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type HttpClient interface {
	Post(url string, body interface{}) (rep *string, err error)
	Get(url string, body interface{}) (rep *string, err error)
	Send() (body *string, err error)
}

func New() *SuperAgent {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &SuperAgent{
		Url:               "",
		Method:            "",
		Header:            http.Header{},
		Data:              make(map[string]interface{}),
		SliceData:         []interface{}{},
		FormData:          url.Values{},
		QueryData:         url.Values{},
		BounceToRawString: false,
		RawString:         "",
		ForceType:         "",
		TargetType:        TypeJSON,
		Cookies:           make([]*http.Cookie, 0),
		Errors:            nil,
		Client: &http.Client{
			Timeout:   time.Second * 30,
			Transport: netTransport,
		},
	}
}

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SuperAgent struct {
	Url               string
	Method            string
	Header            http.Header
	TargetType        string
	ForceType         string
	Data              map[string]interface{}
	SliceData         []interface{}
	FormData          url.Values
	QueryData         url.Values
	BounceToRawString bool
	RawString         string
	Client            *http.Client
	Transport         *http.Transport
	Cookies           []*http.Cookie
	Errors            []error
	BasicAuth         struct{ Username, Password string }
	Debug             bool
	CurlCommand       bool
	Body              io.Reader
}

func (s *SuperAgent) Post(url string, body string) (rep *http.Response, by []byte, err error) {
	s.Method = MethodPost
	s.Url = url
	s.Body = bytes.NewBuffer([]byte(body))
	s.Header = map[string][]string{
		"Content-Type": []string{"application/json"},
	}
	rep, by, err = s.Send()
	return
}

func (s *SuperAgent) Get(url string, body interface{}) (rep *http.Response, by []byte, err error) {
	s.Method = MethodGet
	s.Url = url
	rep, by, err = s.Send()
	return
}

func (s *SuperAgent) Send() (rep *http.Response, body []byte, err error) {
	request, err := http.NewRequest(s.Method, s.Url, s.Body)
	if err != nil {
		return
	}
	request.Header = s.Header
	resp, err := s.Client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fmt.Sprintf("rep :%v", resp)
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return resp, body, nil
	} else {
		return resp, nil, nil
	}
}



func changeMapToURLValues(data map[string]interface{}) url.Values {
	var newUrlValues = url.Values{}
	for k, v := range data {
		switch val := v.(type) {
		case string:
			newUrlValues.Add(k, val)
		case bool:
			newUrlValues.Add(k, strconv.FormatBool(val))
		// if a number, change to string
		// json.Number used to protect against a wrong (for GoRequest) default conversion
		// which always converts number to float64.
		// This type is caused by using Decoder.UseNumber()
		case json.Number:
			newUrlValues.Add(k, string(val))
		case int:
			newUrlValues.Add(k, strconv.FormatInt(int64(val), 10))
		// TODO add all other int-Types (int8, int16, ...)
		case float64:
			newUrlValues.Add(k, strconv.FormatFloat(float64(val), 'f', -1, 64))
		case float32:
			newUrlValues.Add(k, strconv.FormatFloat(float64(val), 'f', -1, 64))
		// following slices are mostly needed for tests
		case []string:
			for _, element := range val {
				newUrlValues.Add(k, element)
			}
		case []int:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatInt(int64(element), 10))
			}
		case []bool:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatBool(element))
			}
		case []float64:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatFloat(float64(element), 'f', -1, 64))
			}
		case []float32:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatFloat(float64(element), 'f', -1, 64))
			}
		// these slices are used in practice like sending a struct
		case []interface{}:

			if len(val) <= 0 {
				continue
			}

			switch val[0].(type) {
			case string:
				for _, element := range val {
					newUrlValues.Add(k, element.(string))
				}
			case bool:
				for _, element := range val {
					newUrlValues.Add(k, strconv.FormatBool(element.(bool)))
				}
			case json.Number:
				for _, element := range val {
					newUrlValues.Add(k, string(element.(json.Number)))
				}
			}
		default:
			// TODO add ptr, arrays, ...
		}
	}
	return newUrlValues
}
