package decoder_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

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
}
`)

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
]
`)

	singleYAML = []byte(`
null_field: null
string_field: かきくけこ
int_field: 45678
float32_field: 9.876
bool_field: false
struct_ptr_field:
  key: hugakey
  value: hugavalue
array_string_field:
  - array_str_1
  - array_str_2
array_struct_field:
  - kkk: kkk1
    vvvv: vvv1
  - kkk: kkk2
    vvvv: vvv2
  - kkk: kkk3
    vvvv: vvv3
`)

	arrayYAML = []byte(`
- null_field: null
  string_field: かきくけこ
  int_field: 45678
  float32_field: 9.876
  bool_field: false
  struct_ptr_field:
    key: hugakey
    value: hugavalue
  array_string_field:
    - array_str_1
    - array_str_2
  array_struct_field:
    - kkk: kkk1
      vvvv: vvv1
    - kkk: kkk2
      vvvv: vvv2
    - kkk: kkk3
      vvvv: vvv3
- null_field: null
  string_field: さしすせそ
  int_field: 7890
  float32_field: 4.99
  bool_field: true
  struct_ptr_field:
    key: hugakeyXXX
    value: hugavalueXXX
  array_string_field:
    - array_str_111
    - array_str_222
  array_struct_field:
    - kkk: kkk99
      vvvv: vvv99
    - kkk: kkk999
      vvvv: vvv999
    - kkk: kkk9999
      vvvv: vvv9999
`)

	singleTOML = []byte(`
string_field = "かきくけこ,"
int_field = 45678
float32_field = "9.876,"
bool_field = false
array_string_field = ["array_str_1", "array_str_2"]

[struct_ptr_field]
  key = "hugakey"
  value = "hugavalue"

[[array_struct_field]]
  kkk = "kkk1"
  vvvv = "vvv1"

[[array_struct_field]]
  kkk = "kkk2"
  vvvv = "vvv2"

[[array_struct_field]]
  kkk = "kkk3"
  vvvv = "vvv3"
`)

	singleXML = []byte(`
<?xml version="1.0" encoding="UTF-8" ?>
<root>
    <null_field/>
    <string_field>かきくけこ</string_field>
    <int_field>45678</int_field>
    <float32_field>9.876</float32_field>
    <bool_field>false</bool_field>
    <struct_ptr_field>
        <key>hugakey</key>
        <value>hugavalue</value>
    </struct_ptr_field>
    <array_string_field>array_str_1</array_string_field>
    <array_string_field>array_str_2</array_string_field>
    <array_struct_field>
        <kkk>kkk1</kkk>
        <vvvv>vvv1</vvvv>
    </array_struct_field>
    <array_struct_field>
        <kkk>kkk2</kkk>
        <vvvv>vvv2</vvvv>
    </array_struct_field>
    <array_struct_field>
        <kkk>kkk3</kkk>
        <vvvv>vvv3</vvvv>
    </array_struct_field>
</root>
`)

	arrayXML = []byte(`
<?xml version="1.0" encoding="UTF-8" ?>
<root>
    <0>
        <null_field/>
        <string_field>かきくけこ</string_field>
        <int_field>45678</int_field>
        <float32_field>9.876</float32_field>
        <bool_field>false</bool_field>
        <struct_ptr_field>
            <key>hugakey</key>
            <value>hugavalue</value>
        </struct_ptr_field>
        <array_string_field>array_str_1</array_string_field>
        <array_string_field>array_str_2</array_string_field>
        <array_struct_field>
            <kkk>kkk1</kkk>
            <vvvv>vvv1</vvvv>
        </array_struct_field>
        <array_struct_field>
            <kkk>kkk2</kkk>
            <vvvv>vvv2</vvvv>
        </array_struct_field>
        <array_struct_field>
            <kkk>kkk3</kkk>
            <vvvv>vvv3</vvvv>
        </array_struct_field>
    </0>
    <1>
        <null_field/>
        <string_field>さしすせそ</string_field>
        <int_field>7890</int_field>
        <float32_field>4.99</float32_field>
        <bool_field>true</bool_field>
        <struct_ptr_field>
            <key>hugakeyXXX</key>
            <value>hugavalueXXX</value>
        </struct_ptr_field>
        <array_string_field>array_str_111</array_string_field>
        <array_string_field>array_str_222</array_string_field>
        <array_struct_field>
            <kkk>kkk99</kkk>
            <vvvv>vvv99</vvvv>
        </array_struct_field>
        <array_struct_field>
            <kkk>kkk999</kkk>
            <vvvv>vvv999</vvvv>
        </array_struct_field>
        <array_struct_field>
            <kkk>kkk9999</kkk>
            <vvvv>vvv9999</vvvv>
        </array_struct_field>
    </1>
</root>
`)

	singleHCL = []byte(`
"null_field" =

"string_field" = "かきくけこ,"

"int_field" = 45678

"float32_field" = "9.876,"

"bool_field" = false

"struct_ptr_field" = {
    "key" = "hugakey"

    "value" = "hugavalue"
}

"array_string_field" = ["array_str_1", "array_str_2"]

"array_struct_field" = {
    "kkk" = "kkk1"

    "vvvv" = "vvv1"
}

"array_struct_field" = {
    "kkk" = "kkk2"

    "vvvv" = "vvv2"
}

"array_struct_field" = {
    "kkk" = "kkk3"

    "vvvv" = "vvv3"
}
`)
)

func TestDynamicStructHasOnlyPrimitive(t *testing.T) {
	t.Parallel()

	data := []byte(`
{
  "null_field":null,
  "string_field":"かきくけこ",
  "int_field":45678,
  "float32_field":9.876,
  "bool_field":false
}
`)
	wantDef := `type DynamicStruct struct {
	BoolField bool
	Float32Field float64
	IntField float64
	NullField interface {}
	StringField string
}`

	t.Run("TestDynamicStructHasOnlyPrimitive", func(t *testing.T) {
		testCorrectCase(t, data, TypeJSON, false, false, 5, wantDef)
	})
}

func TestDynamicStructHasObj(t *testing.T) {
	t.Parallel()

	data := []byte(`
{
  "string_field":"あああ",
  "obj_field":{
    "id":123,
    "name":"Test Tarou"
  }
}
`)
	wantDef := `type DynamicStruct struct {
	ObjField map[string]interface {}
	StringField string
}`

	t.Run("TestDynamicStructHasObj", func(t *testing.T) {
		testCorrectCase(t, data, TypeJSON, false, false, 2, wantDef)
	})
}

func TestDynamicStructHasArrayString(t *testing.T) {
	t.Parallel()

	data := []byte(`
{
  "string_field":"あああ",
  "string_array_field":[
    "id1",
    "id2"
  ]
}
`)
	wantDef := `type DynamicStruct struct {
	StringArrayField []string
	StringField string
}`

	t.Run("TestDynamicStructHasArrayString", func(t *testing.T) {
		testCorrectCase(t, data, TypeJSON, false, false, 2, wantDef)
	})
}

func TestDynamicStructHasArrayInt(t *testing.T) {
	t.Parallel()

	data := []byte(`
{
  "string_field":"あああ",
  "int_array_field":[
    3,
    4
  ]
}
`)
	wantDef := `type DynamicStruct struct {
	IntArrayField []float64
	StringField string
}`

	t.Run("TestDynamicStructHasArrayInt", func(t *testing.T) {
		testCorrectCase(t, data, TypeJSON, false, false, 2, wantDef)
	})
}

func testCorrectCase(t *testing.T, data []byte, dt DataType, nest bool, useTag bool, wantNumF int, wantDef string) {
	d, err := New(data, dt)
	if err != nil {
		t.Errorf("unexpected error is returned from New: %v", err)
		return
	}

	ds, err := d.DynamicStruct(nest, useTag)
	if err != nil {
		t.Errorf("unexpected error is returned from Decode: %v", err)
		return
	}

	if ds == nil {
		t.Errorf("unexpected DynamicStruct is null. got: is null, want: is not null")
		return
	}

	if ds.NumField() != wantNumF {
		t.Errorf("unmatch numfield. got: %d, want: %d, ds.Definition:\n%s", ds.NumField(), wantNumF, ds.Definition())
		return
	}

	if d := cmp.Diff(ds.Definition(), wantDef); d != "" {
		t.Errorf("mismatch Definition: (-got +want)\n%s", d)
		return
	}
}

/*
func TestDecode(t *testing.T) {
	t.Parallel()

	type args struct {
		data     []byte
		dataType DataType
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
				data:     singleJSON,
				dataType: TypeJSON,
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     8,
		},
		{
			name: "JSON is valid array",
			args: args{
				data:     arrayJSON,
				dataType: TypeJSON,
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     8,
		},
		{
			name: "Only one null field",
			args: args{
				data:     []byte(`{"nullfield":null}`),
				dataType: TypeJSON,
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     1,
		},
		{
			name: "Empty JSON",
			args: args{
				data:     []byte(`{}`),
				dataType: TypeJSON,
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     0,
		},
		{
			name: "Empty array JSON",
			args: args{
				data:     []byte(`[]`),
				dataType: TypeJSON,
			},
			wantError:    false,
			wantDsIsNull: true,
			numField:     0,
		},
		{
			name: "empty with TypeJSON",
			args: args{
				data:     []byte(``),
				dataType: TypeJSON,
			},
			wantError:    true,
			wantDsIsNull: true,
			numField:     0,
		},
		{
			name: "null with TypeJSON",
			args: args{
				data:     []byte(`null`),
				dataType: TypeJSON,
			},
			wantError:    true,
			wantDsIsNull: true,
			numField:     0,
		},
		{
			name: "Invalid string with TypeJSON",
			args: args{
				data:     []byte(`invalid`),
				dataType: TypeJSON,
			},
			wantError:    true,
			wantDsIsNull: true,
			numField:     0,
		},
		{
			name: "Invalid DataType",
			args: args{
				data:     singleJSON,
				dataType: -1,
			},
			wantError:    true,
			wantDsIsNull: true,
			numField:     0,
		},
		{
			name: "YAML does not have null field",
			args: args{
				data:     singleYAML,
				dataType: TypeYAML,
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     8,
		},
		{
			name: "YAML is valid array",
			args: args{
				data:     arrayYAML,
				dataType: TypeYAML,
			},
			wantError:    false,
			wantDsIsNull: false,
			numField:     8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dr, err := Decode(tt.args.data, tt.args.dataType)
			if err == nil {
				if tt.wantError {
					t.Errorf("error did not occur. Interface: %#v", dr.Interface)
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
						t.Errorf("unexpected DynamicStruct is null. got: is null, want: is not null, Interface:\n%#v", dr.Interface)
						return
					}
				}

			} else if !tt.wantError {
				t.Errorf("unexpected error occurred. wantError %v, err: %v", tt.wantError, err)
			}
		})
	}
}

func TestDecodeInvalidType(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "YAML",
			args: args{
				data: singleYAML,
			},
		},
		{
			name: "Invalid string",
			args: args{
				data: []byte(`invalid`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr, err := Decode(tt.args.data, -1)
			if err == nil {
				t.Errorf("error did not occur. Interface: %#v", dr.Interface)
				return
			}
		})
	}
}
*/

// benchmark tests

func BenchmarkDynamicStructSingleJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := New(singleJSON, TypeJSON)
		_, _ = d.DynamicStruct(false, false)
	}
}

func BenchmarkDynamicStructArrayJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, _ := New(arrayJSON, TypeJSON)
		_, _ = d.DynamicStruct(false, false)
	}
}
