package decoder

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"
)

// dataType is the type of original data format
// This type provides an unified interface for marshal and unmarshal functions per data formats.
type dataType int

const (
	typeJSON dataType = iota
	typeYAML
	typeHCL

	// FIXME: futures as follows
	// TypeXML
	// TypeTOML
	// TypeCSV

	end // end of iota
)

var formats = [...]string{
	typeJSON: "json",
	typeYAML: "yaml",
	typeHCL:  "hcl",
}

func (dt dataType) string() string {
	if dt >= 0 && int(dt) < len(formats) {
		return formats[dt]
	}
	return ""
}

func (dt dataType) unmarshal(data []byte) (ret interface{}, err error) {
	switch dt {
	case typeHCL:
		// Note: hclsimple.Decode supports only pointer of map or struct.
		var m map[string]interface{}
		err = hclsimple.Decode("example.hcl", data, nil, &m)
		if err != nil {
			return nil, err
		}
		dec, err := decodeHCL(data)
		if err != nil {
			return nil, err
		}
		return interface{}(dec), nil
	default:
		err = dt.unmarshalWithPtr(data, &ret)
	}

	return
}

func (dt dataType) unmarshalWithPtr(data []byte, iptr interface{}) (err error) {
	switch dt {
	case typeJSON:
		err = json.Unmarshal(data, iptr)
	case typeYAML:
		err = yaml.Unmarshal(data, iptr)
	default:
		err = fmt.Errorf("invalid datatype for Unmarshal: %v", dt)
	}

	return
}

// TODO: add tests and examples
// func (dt dataType) marshal(v interface{}) (data []byte, err error) {
func (dt dataType) marshal(m map[string]interface{}) (data []byte, err error) {
	switch dt {
	case typeJSON:
		data, err = json.Marshal(m)
	case typeYAML:
		data, err = yaml.Marshal(m)
	default:
		err = fmt.Errorf("invalid datatype for Marshal: %v", dt)
	}

	return
}

func decodeHCL(data []byte) (map[string]interface{}, error) {
	// Note: hclsimple.Decode supports only pointer of map or struct.
	var m map[string]interface{}
	err := hclsimple.Decode("example.hcl", data, nil, &m)
	if err != nil {
		return m, err
	}

	decoded := make(map[string]interface{}, len(m))
	for k, v := range m {
		attr, ok := v.(*hcl.Attribute)
		if !ok {
			return decoded, fmt.Errorf("%q field can not cast to *hcl.Attribute. v = %v", k, v)
		}

		ctyVal, _ := attr.Expr.Value(nil)
		decoded[k], err = convCtyToGo(ctyVal)
		if err != nil {
			return decoded, err
		}
	}

	return decoded, nil
}

func convCtyToGo(ctyVal cty.Value) (interface{}, error) {
	var err error

	ctyType := ctyVal.Type()
	if ctyType == cty.String {
		return ctyVal.AsString(), nil
	} else if ctyType == cty.Number {
		return ctyVal.AsBigFloat(), nil
	} else if ctyType == cty.Bool {
		return ctyVal.True(), nil
	} else if ctyType.IsTupleType() {
		vals := ctyVal.AsValueSlice()
		ret := make([]interface{}, len(vals))
		for i, v := range vals {
			ret[i], err = convCtyToGo(v)
			if err != nil {
				return nil, err
			}
		}
		return ret, nil
	} else if ctyType.IsObjectType() {
		valM := ctyVal.AsValueMap()
		retM := make(map[string]interface{}, len(valM))
		for k, v := range valM {
			retM[k], err = convCtyToGo(v)
			if err != nil {
				return nil, err
			}
		}
		return retM, nil
	} else if ctyType == cty.DynamicPseudoType {
		// FIXME: just support only null?
		return nil, nil
	} else {
		return nil, fmt.Errorf("unsupported ctyType: %v", ctyType)
	}

}
