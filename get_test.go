package gese

import (
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
			name: "empty_args",
			want: nil,
		},
		{
			name: "default_value",
			args: args{
				[]interface{}{
					[]int{},
					2,
					3,
				},
			},
			want: 3,
		},
		{
			name: "default_value_by_invalid_from",
			args: args{
				[]interface{}{
					nil,
				},
			},
			want: nil,
		},
		{
			name: "replace_zero_value",
			args: args{
				[]interface{}{
					string([]rune{0, 0, 0, 0, 0}),
					2,
					4,
					true,
				},
			},
			want: 4,
		},
		{
			name: "replace_zero_value_arr",
			args: args{
				[]interface{}{
					[]int{0, 0, 0, 0, 0},
					2,
					4,
					true,
				},
			},
			want: 4,
		},
		{
			name: "arr_in_arr",
			args: args{
				[]interface{}{
					[]interface{}{1, []int{1, 2, 3}},
					"1.1",
				},
			},
			want: 2,
		},
		{
			name: "replace_zero_val_struct",
			args: args{
				[]interface{}{
					tA{},
					"B",
					123,
					true,
				},
			},
			want: 123,
		},
		{
			name: "path_interface_str_map",
			args: args{
				[]interface{}{
					m,
					[]interface{}{"123456"},
					123,
				},
			},
			want: 123,
		},
		{
			name: "path_interface_int_map",
			args: args{
				[]interface{}{
					map[interface{}]interface{}{
						1: "abc",
					},
					[]interface{}{2},
					123,
				},
			},
			want: 123,
		},
		{
			name: "path_interface_unsupport_type_map",
			args: args{
				[]interface{}{
					map[float64]interface{}{
						1: "abc",
					},
					[]interface{}{2},
					123,
				},
			},
			want: 123,
		},
		{
			name: "replace_zero_map",
			args: args{
				[]interface{}{
					map[string]string{
						"123": "",
					},
					"123",
					456,
					true,
				},
			},
			want: 456,
		},
		{
			name: "unsupport_from",
			args: args{
				[]interface{}{
					6,
					"a.b.c",
				},
			},
			want: nil,
		},
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

type tE struct {
	F string
	V interface{}
}

type tD struct {
	E *tE
	M map[string]interface{}
}

type tC struct {
	D tD
	V interface{}
}

type tB struct {
	C *tC
	V interface{}
}

type tA struct {
	B tB
	V interface{}
}
