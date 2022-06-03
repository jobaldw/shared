package client

import (
	"io"
	"testing"

	"github.com/jobaldw/shared/config"
	"github.com/jobaldw/shared/test/mock"
)

func TestClient_GetBody(t *testing.T) {
	svr := mock.Server()
	defer svr.Close()

	type client config.Client
	type args struct {
		path   string
		params map[string][]string
		body   interface{}
	}
	tests := []struct {
		name string
		conf client
		args args
		resp io.Reader
	}{
		{
			name: "healthy",
			conf: client{
				URL:    svr.URL,
				Health: "/health",
			},
			args: args{
				path:   "/save",
				params: nil,
				body:   "test payload",
			},
			resp: &io.LimitedReader{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := New(config.Client(test.conf))
			if err != nil {
				t.Errorf("Response.GetBody() error = %s", err)
				return
			}

			response, err := client.Post(test.args.path, test.args.params, test.args.body)
			if err != nil {
				t.Errorf("Response.GetBody() error = %s", err)
				return
			}

			response.GetBody()
		})
	}
}

func TestClient_GetBodyBytes(t *testing.T) {
	svr := mock.Server()
	defer svr.Close()

	type client config.Client
	type args struct {
		path   string
		params map[string][]string
		body   interface{}
	}
	tests := []struct {
		name string
		conf client
		args args
		resp io.Reader
	}{
		{
			name: "healthy",
			conf: client{
				URL:    svr.URL,
				Health: "/health",
			},
			args: args{
				path:   "/save",
				params: nil,
				body:   "test payload",
			},
			resp: &io.LimitedReader{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := New(config.Client(test.conf))
			if err != nil {
				t.Errorf("Response.GetBodyBytes() error = %s", err)
				return
			}

			response, err := client.Post(test.args.path, test.args.params, test.args.body)
			if err != nil {
				t.Errorf("Response.GetBodyBytes() error = %s", err)
				return
			}

			response.GetBodyBytes()
		})
	}
}

func TestClient_GetBodyString(t *testing.T) {
	svr := mock.Server()
	defer svr.Close()

	type client config.Client
	type args struct {
		path   string
		params map[string][]string
		body   interface{}
	}
	tests := []struct {
		name string
		conf client
		args args
		resp io.Reader
	}{
		{
			name: "healthy",
			conf: client{
				URL:    svr.URL,
				Health: "/health",
			},
			args: args{
				path:   "/save",
				params: nil,
				body:   "test payload",
			},
			resp: &io.LimitedReader{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := New(config.Client(test.conf))
			if err != nil {
				t.Errorf("Response.GetBodyString() error = %s", err)
				return
			}

			response, err := client.Post(test.args.path, test.args.params, test.args.body)
			if err != nil {
				t.Errorf("Response.GetBodyString() error = %s", err)
				return
			}

			response.GetBodyString()
		})
	}
}
