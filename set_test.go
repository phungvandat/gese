package gese

import (
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	var (
		destStr1 = "abc"
		destStr2 = "xyz"
	)

	type args struct {
		dest   interface{}
		path   interface{}
		setVal interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		want    interface{}
	}{
		{
			name: "invalid_path",
			args: args{
				dest:   tATest(),
				path:   "a..b",
				setVal: "abc",
			},
			wantErr: ErrPathNotExists,
		},
		{
			name: "path_str_not_exists",
			args: args{
				dest: func() *string {
					val := "abc"
					return &val
				}(),
				path:   "abc",
				setVal: "123",
			},
			wantErr: ErrPathNotExists,
		},
		{
			name: "bad_value_for_str",
			args: args{
				dest: func() *string {
					val := "abc"
					return &val
				}(),
				path:   "1",
				setVal: "123",
			},
			wantErr: ErrBadValue,
		},
		{
			name: "path_struct_not_exists_with_number",
			args: args{
				dest:   tATest(),
				path:   "1",
				setVal: "123",
			},
			wantErr: ErrPathNotExists,
		},
		{
			name: "path_struct_not_exists",
			args: args{
				dest:   tATest(),
				path:   "xyz",
				setVal: "123",
			},
			wantErr: ErrPathNotExists,
		},
		{
			name: "bad_type_struct",
			args: args{
				dest:   tATest(),
				path:   "B",
				setVal: "123",
			},
			wantErr: ErrBadValue,
		},
		{
			name: "failed_by_not_ptr",
			args: args{
				dest: "abc",
			},
			wantErr: ErrInvalidDest,
		},
		{
			name: "str_num_path",
			args: args{
				dest:   &destStr1,
				path:   0,
				setVal: rune('1'),
			},
			want: func() *string {
				var v = "1bc"
				return &v
			}(),
		},
		{
			name: "str_str_path",
			args: args{
				dest:   &destStr2,
				path:   "1",
				setVal: rune('5'),
			},
			want: func() *string {
				var v = "x5z"
				return &v
			}(),
		},
		{
			name: "struct_str_path_1",
			args: args{
				dest: tATest(),
				path: "B",
				setVal: tB{
					C: &tC{
						V: "hello",
					},
				},
			},
			want: &tA{
				B: tB{
					C: &tC{
						V: "hello",
					},
				},
				V: "1234",
			},
		},
		{
			name: "struct_str_path_2",
			args: args{
				dest:   tATest(),
				path:   "B.C",
				setVal: &tC{},
			},
			want: &tA{
				B: tB{
					C: &tC{},
				},
				V: "1234",
			},
		},
		{
			name: "struct_str_path_3",
			args: args{
				dest: &tA{},
				path: "B.C.D.E",
				setVal: &tE{
					F: "f123",
				},
			},
			want: &tA{
				B: tB{
					C: &tC{
						D: tD{
							E: &tE{
								F: "f123",
							},
						},
					},
				},
			},
		},
		{
			name: "struct_str_path_4",
			args: args{
				dest: tATest(),
				path: "B.C.D.E",
				setVal: &tE{
					F: "f123",
				},
			},
			want: &tA{
				B: tB{
					C: &tC{
						D: tD{
							E: &tE{
								F: "f123",
							},
							M: map[string]interface{}{
								"xyz": 1010,
								"hello_world": map[string]interface{}{
									"efg": 1919,
									"5":   "hello",
								},
							},
						},
					},
				},
				V: "1234",
			},
		},
		{
			name: "struct_str_path_5",
			args: args{
				dest: tATest(),
				path: "B.C.D.M.hello_world.efg",
				setVal: &tE{
					F: "f123",
				},
			},
			want: &tA{
				B: tB{
					C: &tC{
						D: tD{
							E: &tE{F: "abc"},
							M: map[string]interface{}{
								"xyz": 1010,
								"hello_world": map[string]interface{}{
									"efg": &tE{
										F: "f123",
									},
									"5": "hello",
								},
							},
						},
					},
				},
				V: "1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Set(tt.args.dest, tt.args.path, tt.args.setVal)
			if err != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(tt.args.dest, tt.want) {
				t.Errorf("Set() = %v, want %v", tt.args.dest, tt.want)
			}
		})
	}
}

func tATest() *tA {
	var (
		m = map[string]interface{}{
			"xyz": 1010,
			"hello_world": map[string]interface{}{
				"efg": 1919,
				"5":   "hello",
			},
		}
		f = "abc"
		e = tE{F: f}
		d = tD{E: &e, M: m}
		c = tC{D: d}
		b = tB{C: &c}
		a = tA{B: b, V: "1234"}
	)
	return &a
}
