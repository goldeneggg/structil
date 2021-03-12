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
}

func TestDynamicStructJSON(t *testing.T) {
	t.Parallel()

	tests := []decoderTest{
		{
			name: "HasOnlyPrimitiveJSON",
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
			name: "HasObjJSON",
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
			name: "HasObjJSONWithTag",
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
			name: "HasObjJSONWithNest",
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
			name: "HasObjJSONWithTagNest",
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
			useTag:   true,
			wantNumF: 2,
			wantDefinition: `type DynamicStruct struct {
	ObjField struct {
		Id float64 ` + "`json:\"id\"`" + `
		Name string ` + "`json:\"name\"`" + `
	}
	StringField string ` + "`json:\"string_field\"`" + `
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCorrectCase(t, tt.data, tt.dt, tt.nest, tt.useTag, tt.wantNumF, tt.wantDefinition)
		})
	}
}

// FIXME: migrate table-driven test in TestDynamicStructJSON
func TestDynamicStructHasObjTwoNestJSON(t *testing.T) {
	t.Parallel()

	data := []byte(`
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
`)
	wantDef := `type DynamicStruct struct {
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
}`

	t.Run("TestDynamicStructHasObj", func(t *testing.T) {
		testCorrectCase(t, data, TypeJSON, true, false, 2, wantDef)
	})
}

// FIXME: migrate table-driven test in TestDynamicStructJSON
func TestDynamicStructHasArrayStringJSON(t *testing.T) {
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

// FIXME: migrate table-driven test in TestDynamicStructYAML
func TestDynamicStructHasOnlyPrimitiveYAML(t *testing.T) {
	t.Parallel()

	data := []byte(`
null_field: null
string_field: かきくけこ
int_field: 45678
float32_field: 9.876
bool_field: false
`)
	wantDef := `type DynamicStruct struct {
	BoolField bool
	Float32Field float64
	IntField int
	NullField interface {}
	StringField string
}`
	wantDefTag := `type DynamicStruct struct {
	BoolField bool ` + "`yaml:\"bool_field\"`" + `
	Float32Field float64 ` + "`yaml:\"float32_field\"`" + `
	IntField int ` + "`yaml:\"int_field\"`" + `
	NullField interface {} ` + "`yaml:\"null_field\"`" + `
	StringField string ` + "`yaml:\"string_field\"`" + `
}`

	t.Run("TestDynamicStructHasOnlyPrimitive", func(t *testing.T) {
		testCorrectCase(t, data, TypeYAML, false, false, 5, wantDef)
		testCorrectCase(t, data, TypeYAML, false, true, 5, wantDefTag)
	})
}

// FIXME: migrate table-driven test in TestDynamicStructYAML
func TestDynamicStructHasObjYAML(t *testing.T) {
	t.Parallel()

	data := []byte(`
string_field: かきくけこ
obj_field:
  id: 123
  name: Test Tarou
`)
	wantDef := `type DynamicStruct struct {
	ObjField map[string]interface {}
	StringField string
}`
	wantDefTag := `type DynamicStruct struct {
	ObjField map[string]interface {} ` + "`yaml:\"obj_field\"`" + `
	StringField string ` + "`yaml:\"string_field\"`" + `
}`
	wantDefNest := `type DynamicStruct struct {
	ObjField struct {
		Id int
		Name string
	}
	StringField string
}`
	wantDefTagNest := `type DynamicStruct struct {
	ObjField struct {
		Id int ` + "`yaml:\"id\"`" + `
		Name string ` + "`yaml:\"name\"`" + `
	}
	StringField string ` + "`yaml:\"string_field\"`" + `
}`

	t.Run("TestDynamicStructHasObj", func(t *testing.T) {
		testCorrectCase(t, data, TypeYAML, false, false, 2, wantDef)
		testCorrectCase(t, data, TypeYAML, false, true, 2, wantDefTag)
		testCorrectCase(t, data, TypeYAML, true, false, 2, wantDefNest)
		testCorrectCase(t, data, TypeYAML, true, true, 2, wantDefTagNest)
	})
}

// FIXME: migrate table-driven test in TestDynamicStructYAML
func TestDynamicStructHasObjTwoNestYAML(t *testing.T) {
	t.Parallel()

	data := []byte(`
string_field: あああ
obj_field:
  id: 45
  name: Test Jiou
  boss: true
  objobj_field:
    user_id: 678
    status: progress
`)
	wantDef := `type DynamicStruct struct {
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
}`

	t.Run("TestDynamicStructHasObj", func(t *testing.T) {
		testCorrectCase(t, data, TypeYAML, true, false, 2, wantDef)
	})
}

// FIXME: migrate table-driven test in TestDynamicStructYAML
func TestDynamicStructHasArrayStringYAML(t *testing.T) {
	t.Parallel()

	data := []byte(`
string_field: あああ
string_array_field:
  - id1
  - id2
`)
	wantDef := `type DynamicStruct struct {
	StringArrayField []string
	StringField string
}`

	t.Run("TestDynamicStructHasArrayString", func(t *testing.T) {
		testCorrectCase(t, data, TypeYAML, false, false, 2, wantDef)
	})
}

func testCorrectCase(t *testing.T, data []byte, dt DataType, nest bool, useTag bool, wantNumF int, wantDef string) {
	t.Helper()

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
