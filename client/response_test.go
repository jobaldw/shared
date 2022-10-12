package client

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jobaldw/shared/v2/config"
	"github.com/jobaldw/shared/v2/internal/test/mock"
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

func TestResponse_IsSuccessful(t *testing.T) {
	opts := cmp.Options{}

	type args struct {
		StatusCode int
	}
	tests := []struct {
		name string
		args args
		resp bool
	}{
		{
			name: "successful",
			args: args{StatusCode: 200},
			resp: true,
		},
		{
			name: "unsuccessful",
			args: args{StatusCode: 300},
			resp: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := &Response{
				StatusCode: test.args.StatusCode,
			}
			got := r.IsSuccessful()
			if diff := cmp.Diff(test.resp, got, opts); diff != "" {
				t.Errorf("Response.IsSuccessful() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResponse_IsClientError(t *testing.T) {
	opts := cmp.Options{}

	type args struct {
		StatusCode int
	}
	tests := []struct {
		name string
		args args
		resp bool
	}{
		{
			name: "client error",
			args: args{StatusCode: 404},
			resp: true,
		},
		{
			name: "not client error",
			args: args{StatusCode: 501},
			resp: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := &Response{
				StatusCode: test.args.StatusCode,
			}
			got := r.IsClientError()
			if diff := cmp.Diff(test.resp, got, opts); diff != "" {
				t.Errorf("Response.IsClientError() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResponse_IsServerError(t *testing.T) {
	opts := cmp.Options{}

	type args struct {
		StatusCode int
	}
	tests := []struct {
		name string
		args args
		resp bool
	}{
		{
			name: "server error",
			args: args{StatusCode: 500},
			resp: true,
		},
		{
			name: "not server error",
			args: args{StatusCode: 0},
			resp: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := &Response{
				StatusCode: test.args.StatusCode,
			}
			got := r.IsServerError()
			if diff := cmp.Diff(test.resp, got, opts); diff != "" {
				t.Errorf("Response.IsServerError() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
