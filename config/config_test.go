// nolint
package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Unmarshal(t *testing.T) {
	opts := cmp.Comparer(func(x, y error) bool {
		return x.Error() == y.Error()
	})

	newStruct := struct {
		Field1 string `json:"field1"`
	}{}

	type args struct {
		conf interface{}
	}
	tests := []struct {
		name string
		args args
		resp error
	}{
		{
			name: "good struct",
			args: args{
				conf: &newStruct,
			},
			resp: nil,
		},
		{
			name: "bad struct",
			args: args{
				conf: newStruct,
			},
			resp: ErrNonPointerStruct,
		},
		{
			name: "no config files",
			args: args{
				conf: &newStruct,
			},
			resp: ErrConfigsNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.name == "no config files" {
				configDirectory = "config"
			}

			got := Unmarshal(test.args.conf)

			if test.name == "no config files" {
				configDirectory = "configs"
			}

			if diff := cmp.Diff(test.resp, got, opts); diff != "" {
				t.Errorf("Unmarshal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_unmarshal(t *testing.T) {
	opts := cmp.Comparer(func(x, y error) bool {
		return x.Error() == y.Error()
	})

	newStruct := struct {
		Field1 string `json:"field1"`
	}{}

	type args struct {
		path   string
		config interface{}
	}

	tests := []struct {
		name string
		args args
		resp error
	}{
		{
			name: "path_exist",
			args: args{
				path:   "./configs/test.json",
				config: &newStruct,
			},
			resp: nil,
		},
		{
			name: "path does not exist",
			args: args{
				path:   "bad path",
				config: &newStruct,
			},
			resp: ErrConfigsNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := unmarshal(test.args.path, test.args.config)
			if diff := cmp.Diff(test.resp, got, opts); diff != "" {
				t.Errorf("read() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
