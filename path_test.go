package gese

import (
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
