package decoder_test

import (
	"testing"

	. "github.com/goldeneggg/structil/dynamicstruct/decoder"
)

var (
	singleJSON = []byte(`
{
	"null_field":null,
	"string_field":"かきくけこ",
	"int_field":45678,
	"float32_field":9.876,
	"bool_field":false,
	"struct_ptr_field":{
		"key":"hugakey",
		"value":"hugavalue"
	},
	"array_string_field":[
		"array_str_1",
		"array_str_2"
	],
	"array_struct_field":[
		{
			"kkk":"kkk1",
			"vvvv":"vvv1"
		},
		{
			"kkk":"kkk2",
			"vvvv":"vvv2"
		},
		{
			"kkk":"kkk3",
			"vvvv":"vvv3"
		}
	]
}`)

	arrayJSON = []byte(`
[
	{
		"null_field":null,
		"string_field":"かきくけこ",
		"int_field":45678,
		"float32_field":9.876,
		"bool_field":false,
		"struct_ptr_field":{
			"key":"hugakey",
			"value":"hugavalue"
		},
		"array_string_field":[
			"array_str_1",
			"array_str_2"
		],
		"array_struct_field":[
			{
				"kkk":"kkk1",
				"vvvv":"vvv1"
			},
			{
				"kkk":"kkk2",
				"vvvv":"vvv2"
			},
			{
				"kkk":"kkk3",
				"vvvv":"vvv3"
			}
		]
	},
	{
		"null_field":null,
		"string_field":"さしすせそ",
		"int_field":7890,
		"float32_field":4.99,
		"bool_field":true,
		"struct_ptr_field":{
			"key":"hugakeyXXX",
			"value":"hugavalueXXX"
		},
		"array_string_field":[
			"array_str_111",
			"array_str_222"
		],
		"array_struct_field":[
			{
				"kkk":"kkk99",
				"vvvv":"vvv99"
			},
			{
				"kkk":"kkk999",
				"vvvv":"vvv999"
			},
			{
				"kkk":"kkk9999",
				"vvvv":"vvv9999"
			}
		]
	}
]`)
)

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	type args struct {
		jsonData []byte
	}
	tests := []struct {
		name         string
		args         args
		wantError    bool
		wantDsIsNull bool
		numField     int
	}{
		{
			name: "JSON does not have null field",
			args: args{
				jsonData: singleJSON,
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     8,
		},
		{
			name: "JSON is valid array",
			args: args{
				jsonData: arrayJSON,
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     8,
		},
		{
			name: "Only one null field",
			args: args{
				jsonData: []byte(`{"nullfield":null}`),
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     1,
		},
		{
			name: "Empty JSON",
			args: args{
				jsonData: []byte(`{}`),
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     0,
		},
		{
			name: "Empty array JSON",
			args: args{
				jsonData: []byte(`[]`),
			},
			wantError:    false,
			wantDsIsNull: true,
			numField:     0,
		},
		{
			name: "empty",
			args: args{
				jsonData: []byte(``),
			},
			wantError:    true,
			wantDsIsNull: true,
			numField:     0,
		},
		{
			name: "null",
			args: args{
				jsonData: []byte(`null`),
			},
			wantError:    true,
			wantDsIsNull: true,
			numField:     0,
		},
		{
			name: "Invalid string",
			args: args{
				jsonData: []byte(`invalid`),
			},
			wantError:    true,
			wantDsIsNull: true,
			numField:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dr, err := Decode(tt.args.jsonData, TypeJSON)
			if err == nil {
				if tt.wantError {
					t.Errorf("error did not occur. DecodedInterface: %#v", dr.DecodedInterface)
					return
				}

				if dr.DynamicStruct != nil {
					if tt.wantDsIsNull {
						t.Errorf("unexpected DynamicStruct is null. got: is not null, want: is null, ds.Definition:\n%s", dr.DynamicStruct.Definition())
						return
					}
					if dr.DynamicStruct.NumField() != tt.numField {
						t.Errorf("unmatch numfield. got: %d, want: %d, ds.Definition:\n%s", dr.DynamicStruct.NumField(), tt.numField, dr.DynamicStruct.Definition())
						return
					}
				} else {
					if !tt.wantDsIsNull {
						t.Errorf("unexpected DynamicStruct is null. got: is null, want: is not null, DecodedInterface:\n%#v", dr.DecodedInterface)
						return
					}
				}

			} else if !tt.wantError {
				t.Errorf("unexpected error occurred. wantError %v, err: %v", tt.wantError, err)
			}
		})
	}
}

// benchmark tests

func BenchmarkSingleJSONDecode(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode(singleJSON, TypeJSON)
	}
}

func BenchmarkArrayJSONDecode(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode(arrayJSON, TypeJSON)
	}
}
