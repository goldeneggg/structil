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

type decoderTest struct {
	name           string
	data           []byte
	dt             DataType
	nest           bool
	useTag         bool
	wantNumF       int
	wantDefinition string
	wantErrorNew   bool
	wantErrorDs    bool
}

func TestDynamicStructJSON(t *testing.T) {
	t.Parallel()

	tests := []decoderTest{
		{
			name: "HasOnlyPrimitive",
			data: []byte(`
{
	"null_field":null,
	"string_field":"かきくけこ",
	"int_field":45678,
	"float32_field":9.876,
	"bool_field":false
}
`),
			dt:       TypeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 5,
			wantDefinition: `type DynamicStruct struct {
	BoolField bool
	Float32Field float64
	IntField float64
	NullField interface {}
	StringField string
}`,
		},
		{
			name: "HasObj",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":123,
		"name":"Test Tarou"
	}
}
`),
			dt:       TypeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField map[string]interface {}
	StringField string
}`,
		},
		{
			name: "HasObjWithTag",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":123,
		"name":"Test Tarou"
	}
}
`),
			dt:       TypeJSON,
			nest:     false,
			useTag:   true,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField map[string]interface {} ` + "`json:\"obj_field\"`" + `
	StringField string ` + "`json:\"string_field\"`" + `
}`,
		},
		{
			name: "HasObjWithNest",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":123,
		"name":"Test Tarou"
	}
}
`),
			dt:       TypeJSON,
			nest:     true,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Id float64
		Name string
	}
	StringField string
}`,
		},
		{
			name: "HasObjWithTagNest",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":123,
		"name":"Test Tarou",
		"obj_array_field": [
			{
				"k1":"v1",
				"k2":"v2"
			},
			{
				"k1":"v111",
				"k2":"v222"
			}
		]
	}
}
`),
			dt:       TypeJSON,
			nest:     true,
			useTag:   true,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Id float64 ` + "`json:\"id\"`" + `
		Name string ` + "`json:\"name\"`" + `
		ObjArrayField []struct {
			K1 string ` + "`json:\"k1\"`" + `
			K2 string ` + "`json:\"k2\"`" + `
		} ` + "`json:\"obj_array_field\"`" + `
	} ` + "`json:\"obj_field\"`" + `
	StringField string ` + "`json:\"string_field\"`" + `
}`,
		},
		{
			name: "HasObjTwoNest",
			data: []byte(`
{
	"string_field":"あああ",
	"obj_field":{
		"id":45,
		"name":"Test Jiou",
		"boss":true,
		"objobj_field":{
			"user_id":678,
			"status":"progress"
		}
	}
}
`),
			dt:       TypeJSON,
			nest:     true,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Boss bool
		Id float64
		Name string
		ObjobjField struct {
			Status string
			UserId float64
		}
	}
	StringField string
}`,
		},
		{
			name: "HasArrayString1",
			data: []byte(`
{
	"string_field":"あああ",
	"string_array_field":[
		"id1"
	]
}
`),
			dt:       TypeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	StringArrayField []string
	StringField string
}`,
		},
		{
			name: "HasArrayString2",
			data: []byte(`
{
	"string_field":"あああ",
	"string_array_field":[
		"id1",
		"id2"
	]
}
`),
			dt:       TypeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	StringArrayField []string
	StringField string
}`,
		},
		{
			name: "TopLevelIsArray",
			data: []byte(`
[
	{
		"string_field":"あああ",
		"objobj_field":{
			"user_id":678,
			"status":"progress"
		},
		"string_array_field":[
			"id1",
			"id2"
		]
	},
	{
		"string_field":"いいいい",
		"string_array_field":[
			"id4",
			"id5",
			"id6"
		]
	}	
]
`),
			dt:       TypeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 3,
			wantDefinition: `type DynamicStruct struct {
	ObjobjField map[string]interface {}
	StringArrayField []string
	StringField string
}`,
		},
		{
			name:     "BracketOnly",
			data:     []byte(`{}`),
			dt:       TypeJSON,
			nest:     false,
			useTag:   false,
			wantNumF: 0,
			// FIXME: is this ok?
			wantDefinition: `type DynamicStruct struct {
}`,
		},
		{
			name:           "ArrayBracketOnly",
			data:           []byte(`[]`),
			dt:             TypeJSON,
			nest:           false,
			useTag:         false,
			wantNumF:       0,
			wantDefinition: ``,
			wantErrorDs:    true,
		},
		{
			name:           "OnlyLiteral",
			data:           []byte(`aiueo`),
			dt:             TypeJSON,
			nest:           false,
			useTag:         false,
			wantNumF:       0,
			wantDefinition: ``,
			wantErrorNew:   true,
		},
		{
			name:           "Empty",
			data:           []byte(``),
			dt:             TypeJSON,
			nest:           false,
			useTag:         false,
			wantNumF:       0,
			wantDefinition: ``,
			wantErrorNew:   true,
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCorrectCase(t, tt.data, tt.dt, tt.nest, tt.useTag, tt.wantNumF, tt.wantDefinition, tt.wantErrorNew, tt.wantErrorDs)
		})
	}
}

func TestDynamicStructYAML(t *testing.T) {
	t.Parallel()

	tests := []decoderTest{
		{
			name: "HasOnlyPrimitive",
			data: []byte(`
null_field: null
string_field: かきくけこ
int_field: 45678
float32_field: 9.876
bool_field: false
`),
			dt:       TypeYAML,
			nest:     false,
			useTag:   false,
			wantNumF: 5,
			wantDefinition: `type DynamicStruct struct {
	BoolField bool
	Float32Field float64
	IntField int
	NullField interface {}
	StringField string
}`,
		},
		{
			name: "HasObj",
			data: []byte(`
string_field: かきくけこ
obj_field:
  id: 123
  name: Test Tarou
`),
			dt:       TypeYAML,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField map[string]interface {}
	StringField string
}`,
		},
		{
			name: "HasObjWithTag",
			data: []byte(`
string_field: かきくけこ
obj_field:
  id: 123
  name: Test Tarou
`),
			dt:       TypeYAML,
			nest:     false,
			useTag:   true,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField map[string]interface {} ` + "`yaml:\"obj_field\"`" + `
	StringField string ` + "`yaml:\"string_field\"`" + `
}`,
		},
		{
			name: "HasObjWithNest",
			data: []byte(`
string_field: かきくけこ
obj_field:
  id: 123
  name: Test Tarou
`),
			dt:       TypeYAML,
			nest:     true,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Id int
		Name string
	}
	StringField string
}`,
		},
		{
			name: "HasObjWithTagNest",
			data: []byte(`
string_field: かきくけこ
obj_field:
  id: 123
  name: Test Tarou
`),
			dt:       TypeYAML,
			nest:     true,
			useTag:   true,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Id int ` + "`yaml:\"id\"`" + `
		Name string ` + "`yaml:\"name\"`" + `
	}
	StringField string ` + "`yaml:\"string_field\"`" + `
}`,
		},
		{
			name: "HasObjTwoNest",
			data: []byte(`
string_field: あああ
obj_field:
  id: 45
  name: Test Jiou
  boss: true
  objobj_field:
    user_id: 678
    status: progress
`),
			dt:       TypeYAML,
			nest:     true,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Boss bool
		Id int
		Name string
		ObjobjField struct {
			Status string
			UserId int
		}
	}
	StringField string
}`,
		},
		{
			name: "HasArrayString",
			data: []byte(`
string_field: あああ
string_array_field:
  - id1
  - id2
`),
			dt:       TypeYAML,
			nest:     false,
			useTag:   false,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	StringArrayField []string
	StringField string
}`,
		},
		{
			name:           "OnlyLiteral",
			data:           []byte(`aiueo`),
			dt:             TypeYAML,
			nest:           false,
			useTag:         false,
			wantNumF:       0,
			wantDefinition: ``,
			wantErrorDs:    true,
		},
		{
			name:           "Empty",
			data:           []byte(``),
			dt:             TypeYAML,
			nest:           false,
			useTag:         false,
			wantNumF:       0,
			wantDefinition: ``,
			wantErrorDs:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // See: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// FIXME: error cases
			testCorrectCase(t, tt.data, tt.dt, tt.nest, tt.useTag, tt.wantNumF, tt.wantDefinition, tt.wantErrorNew, tt.wantErrorDs)
		})
	}
}

func testCorrectCase(t *testing.T, data []byte, dt DataType, nest bool, useTag bool, wantNumF int, wantDef string, wantErrorNew bool, wantErrorDs bool) {
	t.Helper()

	var d *Decoder
	var err error

	switch dt {
	case TypeJSON:
		d, err = NewJSON(data)
	case TypeYAML:
		d, err = NewYAML(data)
	}
	if err != nil {
		if !wantErrorNew {
			t.Fatalf("unexpected error is returned from NewXXX: %v", err)
		}
		return
	} else if wantErrorNew {
		t.Fatalf("error is expected but it does not occur from NewXXX. data: %s", string(data))
	}

	ds, err := d.DynamicStruct(nest, useTag)
	if err != nil {
		if !wantErrorDs {
			t.Fatalf("unexpected error is returned from DynamicStruct: %v", err)
		}
		return
	} else if wantErrorDs {
		t.Fatalf("error is expected but it does not occur from DynamicStruct. data: %s", string(data))
	}

	if ds == nil {
		t.Fatalf("unexpected DynamicStruct is null. got: is null, want: is not null")
	}

	if ds.NumField() != wantNumF {
		t.Fatalf("unmatch numfield. got: %d, want: %d, ds.Definition:\n%s", ds.NumField(), wantNumF, ds.Definition())
	}

	if d := cmp.Diff(ds.Definition(), wantDef); d != "" {
		t.Fatalf("mismatch Definition: (-got +want)\n%s", d)
	}
}
