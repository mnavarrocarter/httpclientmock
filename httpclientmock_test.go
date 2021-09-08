package httpclientmock_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/mnavarrocarter/httpclientmock"
)

type FakeApiClient struct {
	client  HTTPClient
	baseUrl string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type FakeInput struct {
	Message string `json:"message"`
}

type FakeOutput struct {
	Message string `json:"message"`
}

func (cl *FakeApiClient) mustMakeRequest(ctx context.Context, method, path string, input interface{}) *http.Request {
	var body io.Reader

	if input != nil {
		b, err := json.Marshal(input)
		if err != nil {
			panic(err) // Developer error
		}
		body = bytes.NewBuffer(b)
	}

	url := cl.baseUrl + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err) // Developer error
	}

	req.Header.Add("Content-Type", "application/json")

	return req.WithContext(ctx)
}

func (cl *FakeApiClient) PostInput(ctx context.Context, input *FakeInput) (*FakeOutput, error) {
	req := cl.mustMakeRequest(ctx, "POST", "/input", input)

	res, err := cl.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	out := &FakeOutput{}

	err = json.NewDecoder(res.Body).Decode(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

var postInputTests = []struct {
	name       string
	input      *FakeInput
	mock       *httpclientmock.Mock
	assertions func(t *testing.T, output *FakeOutput, err error)
}{
	{
		name:  "test one",
		input: &FakeInput{"This is a message sent"},
		mock: &httpclientmock.Mock{
			Expect: &httpclientmock.Request{
				Method: "POST",
				Url:    "https://some.fake.service/input",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: []byte(`{"message":"This is a message sent"}`),
			},
			Return: &httpclientmock.Response{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: []byte(`{"message":"This is a message received"}`),
			},
		},
		assertions: func(t *testing.T, output *FakeOutput, err error) {
			if output == nil {
				t.Error("no output")
			}
			if err != nil {
				t.Error("an error has happened")
			}
		},
	},
}

func TestPostInput(t *testing.T) {
	client := &FakeApiClient{http.DefaultClient, "https://some.fake.service"}
	for _, test := range postInputTests {
		t.Run(test.name, func(t *testing.T) {
			restore := test.mock.InjectInClient(t, nil)
			defer restore()
			out, err := client.PostInput(context.Background(), test.input)
			test.assertions(t, out, err)
		})
	}
}
