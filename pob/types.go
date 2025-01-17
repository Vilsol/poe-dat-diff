// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    schema, err := UnmarshalSchema(bytes)
//    bytes, err = schema.Marshal()

package pob

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Schema map[string]*SchemaValue

func UnmarshalSchema(data []byte) (Schema, error) {
	var r Schema
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Schema) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SchemaElement struct {
	RefTo    *string `json:"refTo,omitempty"`
	Width    int64   `json:"width"`
	List     bool    `json:"list"`
	Type     Type    `json:"type"`
	Name     *string `json:"name,omitempty"`
	EnumBase *int64  `json:"enumBase,omitempty"`
}

type PurpleSchema struct {
}

type Type string

const (
	Bool     Type = "Bool"
	Enum     Type = "Enum"
	Float    Type = "Float"
	Int      Type = "Int"
	Interval Type = "Interval"
	Key      Type = "Key"
	ShortKey Type = "ShortKey"
	String   Type = "String"
	UInt     Type = "UInt"
	UInt16   Type = "UInt16"
)

type SchemaValue struct {
	PurpleSchema       *PurpleSchema
	SchemaElementArray []SchemaElement
}

func (x *SchemaValue) UnmarshalJSON(data []byte) error {
	x.SchemaElementArray = nil
	x.PurpleSchema = nil
	var c PurpleSchema
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.SchemaElementArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.PurpleSchema = &c
	}
	return nil
}

func (x *SchemaValue) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.SchemaElementArray != nil, x.SchemaElementArray, x.PurpleSchema != nil, x.PurpleSchema, false, nil, false, nil, false)
}

func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
		*pi = nil
	}
	if pf != nil {
		*pf = nil
	}
	if pb != nil {
		*pb = nil
	}
	if ps != nil {
		*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return false, err
	}

	switch v := tok.(type) {
	case json.Number:
		if pi != nil {
			i, err := v.Int64()
			if err == nil {
				*pi = &i
				return false, nil
			}
		}
		if pf != nil {
			f, err := v.Float64()
			if err == nil {
				*pf = &f
				return false, nil
			}
			return false, errors.New("Unparsable number")
		}
		return false, errors.New("Union does not contain number")
	case float64:
		return false, errors.New("Decoder should not return float64")
	case bool:
		if pb != nil {
			*pb = &v
			return false, nil
		}
		return false, errors.New("Union does not contain bool")
	case string:
		if haveEnum {
			return false, json.Unmarshal(data, pe)
		}
		if ps != nil {
			*ps = &v
			return false, nil
		}
		return false, errors.New("Union does not contain string")
	case nil:
		if nullable {
			return false, nil
		}
		return false, errors.New("Union does not contain null")
	case json.Delim:
		if v == '{' {
			if haveObject {
				return true, json.Unmarshal(data, pc)
			}
			if haveMap {
				return false, json.Unmarshal(data, pm)
			}
			return false, errors.New("Union does not contain object")
		}
		if v == '[' {
			if haveArray {
				return false, json.Unmarshal(data, pa)
			}
			return false, errors.New("Union does not contain array")
		}
		return false, errors.New("Cannot handle delimiter")
	}
	return false, errors.New("Cannot unmarshal union")

}

func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
		return json.Marshal(*pi)
	}
	if pf != nil {
		return json.Marshal(*pf)
	}
	if pb != nil {
		return json.Marshal(*pb)
	}
	if ps != nil {
		return json.Marshal(*ps)
	}
	if haveArray {
		return json.Marshal(pa)
	}
	if haveObject {
		return json.Marshal(pc)
	}
	if haveMap {
		return json.Marshal(pm)
	}
	if haveEnum {
		return json.Marshal(pe)
	}
	if nullable {
		return json.Marshal(nil)
	}
	return nil, errors.New("Union must not be null")
}
