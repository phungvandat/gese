package gese

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	var (
		m = map[string]interface{}{
			"xyz": 1010,
			"abc": map[string]interface{}{
				"efg": 1919,
				"5":   "hello",
			},
		}
		f = "abc"
		e = tE{F: f}
		d = tD{E: &e, M: m}
		c = tC{D: d}
		b = tB{C: &c}
		a = tA{B: b}
	)
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "only_from_value",
			args: args{
				[]interface{}{
					10,
				},
			},
			want: nil,
		},
		{
			name: "from_value_path_exists_string",
			args: args{
				[]interface{}{
					"abcd",
					2,
				},
			},
			want: []rune("c")[0],
		},
		{
			name: "from_value_path_not_exists_string",
			args: args{
				[]interface{}{
					"abcd",
					10,
				},
			},
			want: nil,
		},
		{
			name: "from_value_path_exists_slice",
			args: args{
				[]interface{}{
					[]int{1, 2, 3, 4, 5},
					3,
				},
			},
			want: 4,
		},
		{
			name: "from_value_path_not_exists_slice",
			args: args{
				[]interface{}{
					b,
					10,
				},
			},
			want: nil,
		},
		{
			name: "from_value_string_path_exists_interface_1",
			args: args{
				[]interface{}{
					a,
					"B.C.D.E",
				},
			},
			want: &e,
		},
		{
			name: "from_value_slice_path_exists_interface_1",
			args: args{
				[]interface{}{
					a,
					[]string{"B", "C", "D", "E"},
				},
			},
			want: &e,
		},
		{
			name: "from_value_string_path_exists_interface_2",
			args: args{
				[]interface{}{
					m,
					"abc.efg",
				},
			},
			want: 1919,
		},
		{
			name: "from_value_string_path_exists_interface_3",
			args: args{
				[]interface{}{
					m,
					"abc.5",
				},
			},
			want: "hello",
		},
		{
			name: "from_value_string_path_exists_interface_4",
			args: args{
				[]interface{}{
					m,
					[]interface{}{"abc", 5},
				},
			},
			want: "hello",
		},
		{
			name: "from_value_string_path_not_exists_interface",
			args: args{
				[]interface{}{
					a,
					"B.C.D.G",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.args...); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_get(t *testing.T) {
	type args struct {
		from           interface{}
		path           interface{}
		defaultVal     interface{}
		replaceZeroVal bool
		isFirst        bool
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get(tt.args.from, tt.args.path, tt.args.defaultVal, tt.args.replaceZeroVal, tt.args.isFirst); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}
