// nolint
package client

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/jobaldw/shared/v2/config"
	"github.com/jobaldw/shared/v2/test/mock"
)

func TestClient_IsReady(t *testing.T) {
	svr := mock.Server()
	defer svr.Close()

	opts := cmp.Options{}

	type client config.Client
	type resp struct {
		IsReady bool
		Err     error
	}
	tests := []struct {
		name string
		conf client
		resp resp
	}{
		{
			name: "healthy",
			conf: client{
				URL:    svr.URL,
				Health: "/health",
			},
			resp: resp{
				IsReady: true,
				Err:     nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := New(config.Client(test.conf))
			if err != nil {
				t.Errorf("Client.IsReady() error = %s", err)
				return
			}

			isReady, err := client.IsReady()
			got := resp{
				IsReady: isReady,
				Err:     err,
			}

			if diff := cmp.Diff(test.resp, got, opts); diff != "" {
				t.Errorf("Client.IsReady() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClient_Post(t *testing.T) {
	svr := mock.Server()
	defer svr.Close()

	opts := cmp.Options{
		cmpopts.IgnoreFields(Client{}, "client"),
		cmpopts.IgnoreFields(Response{}, "body", "Request"),
	}

	type client config.Client
	type args struct {
		path   string
		params map[string][]string
		body   interface{}
	}
	type resp struct {
		Response *Response
		Err      error
	}
	tests := []struct {
		name string
		conf client
		args args
		resp resp
	}{
		{
			name: "success",
			conf: client{
				URL: svr.URL,
			},
			args: args{
				path:   "/save",
				params: nil,
				body:   "test payload",
			},
			resp: resp{
				Response: &Response{
					Status:     "201 Created",
					StatusCode: 201,
				},
				Err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := New(config.Client(test.conf))
			if err != nil {
				t.Errorf("Client.Post() error = %s", err)
				return
			}

			response, err := client.Post(test.args.path, test.args.params, test.args.body)
			got := resp{
				Response: response,
				Err:      err,
			}

			if diff := cmp.Diff(test.resp, got, opts); diff != "" {
				t.Errorf("Client.Post() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClient_Put(t *testing.T) {
	svr := mock.Server()
	defer svr.Close()

	opts := cmp.Options{
		cmpopts.IgnoreFields(Client{}, "client"),
		cmpopts.IgnoreFields(Response{}, "body", "Request"),
	}

	type client config.Client
	type args struct {
		path   string
		params map[string][]string
		body   interface{}
	}
	type resp struct {
		Response *Response
		Err      error
	}
	tests := []struct {
		name string
		conf client
		args args
		resp resp
	}{
		{
			name: "success",
			conf: client{
				URL: svr.URL,
			},
			args: args{
				path:   "/update",
				params: nil,
				body:   "test payload",
			},
			resp: resp{
				Response: &Response{
					Status:     "200 OK",
					StatusCode: 200,
				},
				Err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := New(config.Client(test.conf))
			if err != nil {
				t.Errorf("Client.Put() error = %s", err)
				return
			}

			response, err := client.Put(test.args.path, test.args.params, test.args.body)
			got := resp{
				Response: response,
				Err:      err,
			}

			if diff := cmp.Diff(test.resp, got, opts); diff != "" {
				t.Errorf("Client.Put() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClient_Delete(t *testing.T) {
	svr := mock.Server()
	defer svr.Close()

	opts := cmp.Options{
		cmpopts.IgnoreFields(Client{}, "client"),
		cmpopts.IgnoreFields(Response{}, "body", "Request"),
	}

	type client config.Client
	type args struct {
		path   string
		params map[string][]string
	}
	type resp struct {
		Response *Response
		Err      error
	}
	tests := []struct {
		name string
		conf client
		args args
		resp resp
	}{
		{
			name: "success",
			conf: client{
				URL: svr.URL,
			},
			args: args{
				path:   "/delete",
				params: nil,
			},
			resp: resp{
				Response: &Response{
					Status:     "204 No Content",
					StatusCode: 204,
				},
				Err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := New(config.Client(test.conf))
			if err != nil {
				t.Errorf("Client.Delete() error = %s", err)
				return
			}

			response, err := client.Delete(test.args.path, test.args.params)
			got := resp{
				Response: response,
				Err:      err,
			}

			if diff := cmp.Diff(test.resp, got, opts); diff != "" {
				t.Errorf("Client.Delete() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
