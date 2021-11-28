package gese

import (
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	var (
		destStr1 = "abc"
		destStr2 = "xyz"
		m        = map[string]interface{}{
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
				dest: &a,
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
				dest:   &a,
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
				dest: &a,
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
				V: "1234",
			},
		},
		// {
		// 	name: "struct_str_path_5",
		// 	args: args{
		// 		dest: &a,
		// 		path: "B.C.D.M.hello_world.efg",
		// 		setVal: &tE{
		// 			F: "f123",
		// 		},
		// 	},
		// 	want: &tA{
		// 		B: tB{
		// 			C: &tC{
		// 				D: tD{
		// 					M: map[string]interface{}{
		// 						"xyz": 1010,
		// 						"hello_world": map[string]interface{}{
		// 							"efg": "f123",
		// 							"5":   "hello",
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 		V: "1234",
		// 	},
		// },
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
