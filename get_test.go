package gese

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_detectPath(t *testing.T) {
	var (
		n1 = 1
	)
	type args struct {
		path interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []pathItem
		isValid bool
	}{
		{
			name: "pass_string",
			args: args{
				path: "A",
			},
			want: []pathItem{
				{
					val: "A",
					str: "A",
				},
			},
			isValid: true,
		},
		{
			name: "pass_string_dot",
			args: args{
				path: "A.B",
			},
			want: []pathItem{
				{
					val: "A",
					str: "A",
				},
				{
					val: "B",
					str: "B",
				},
			},
			isValid: true,
		},
		{
			name: "pass_string_dot_number",
			args: args{
				path: "A.B.1",
			},
			want: []pathItem{
				{
					val: "A",
					str: "A",
				},
				{
					val: "B",
					str: "B",
				},
				{
					val: "1",
					num: &n1,
					str: "1",
				},
			},
			isValid: true,
		},
		{
			name: "pass_number",
			args: args{
				path: n1,
			},
			want: []pathItem{
				{
					val: n1,
					num: &n1,
					str: "1",
				},
			},
			isValid: true,
		},
		{
			name: "pass_list_string",
			args: args{
				path: [2]string{"A", "B"},
			},
			want: []pathItem{
				{
					val: "A",
					str: "A",
				},
				{
					val: "B",
					str: "B",
				},
			},
			isValid: true,
		},
		{
			name: "pass_list_interface",
			args: args{
				path: [2]interface{}{"A", n1},
			},
			want: []pathItem{
				{
					val: "A",
					str: "A",
				},
				{
					val: n1,
					num: &n1,
					str: "1",
				},
			},
			isValid: true,
		},
		{
			name: "fail_invalid",
			args: args{
				path: [2]interface{}{"A", struct{}{}},
			},
			isValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, isValid := detectPath(tt.args.path)
			if !isValid && isValid != tt.isValid {
				t.Errorf("detectPath() isValid = %v, want %v", isValid, tt.isValid)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("detectPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPosNum(t *testing.T) {
	var (
		n1, n2 = 10, -10
	)
	type args struct {
		numPtr *int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "pass_with_pos_num",
			args: args{
				numPtr: &n1,
			},
			want: true,
		},
		{
			name: "fail_with_neg_num",
			args: args{
				numPtr: &n2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPosNum(tt.args.numPtr); got != tt.want {
				t.Errorf("isPosNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type E struct {
		F string
	}

	type D struct {
		E *E
		M map[string]interface{}
	}

	type C struct {
		D D
	}

	type B struct {
		C *C
	}

	type A struct {
		B
	}

	var (
		m = map[string]interface{}{
			"xyz": 1010,
			"abc": map[string]interface{}{
				"efg": 1919,
				"5":   "hello",
			},
		}
		f = "abc"
		e = E{F: f}
		d = D{E: &e, M: m}
		c = C{D: d}
		b = B{C: &c}
		a = A{B: b}
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
				fmt.Printf("%p, What %p\n", got, tt.want)
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
