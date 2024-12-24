package poedat

import (
	"encoding/json"
	"fmt"
)

type Schema struct {
	Tables []Table `json:"tables,omitempty"`
	//Enumerations []Enumeration `json:"enumerations,omitempty"`
}

//type Enumeration struct {
//	ValidFor    int64     `json:"validFor,omitempty"`
//	Name        string    `json:"name,omitempty"`
//	Indexing    int64     `json:"indexing,omitempty"`
//	Enumerators []*string `json:"enumerators,omitempty"`
//}

type Table struct {
	ValidFor int64    `json:"validFor"`
	Name     string   `json:"name,omitempty"`
	Columns  []Column `json:"columns,omitempty"`
}

type Column struct {
	Name       *string     `json:"name,omitempty"`
	Array      bool        `json:"array,omitempty"`
	Type       Type        `json:"type,omitempty"`
	References *References `json:"references,omitempty"`
}

type References struct {
	Table  string  `json:"table,omitempty"`
	Column *string `json:"column,omitempty"`
}

type Type string

const (
	Array      Type = "array"
	Bool       Type = "bool"
	Enumrow    Type = "enumrow"
	F32        Type = "f32"
	Foreignrow Type = "foreignrow"
	I16        Type = "i16"
	I32        Type = "i32"
	Row        Type = "row"
	String     Type = "string"
	U16        Type = "u16"
	U32        Type = "u32"
)

func Parse(data []byte) (Schema, error) {
	out := Schema{}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, fmt.Errorf("failed parsing schema: %w", err)
	}
	return out, nil
}
