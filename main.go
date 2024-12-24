package main

import (
	"bytes"
	"html/template"
	"maps"
	"os"
	"slices"
	"strings"

	"poe-dat-diff/pob"
	"poe-dat-diff/poedat"
)

type DiffData struct {
	Source poedat.Schema
	Target map[string]poedat.Table
}

type TableDiff struct {
	Source  poedat.Table
	Target  poedat.Table
	ColDiff []ColDiff
}

func main() {
	pobJson, _ := os.ReadFile("pob.json")
	poeDatJson, _ := os.ReadFile("poedat.json")

	poeDat, err := poedat.Parse(poeDatJson)
	if err != nil {
		panic(err)
	}

	dedupTypes := make(map[string]poedat.Table)
	for _, table := range poeDat.Tables {
		if existing, ok := dedupTypes[table.Name]; ok {
			if existing.ValidFor > table.ValidFor {
				continue
			}
		}

		dedupTypes[table.Name] = table
	}

	poeDat.Tables = slices.Collect(maps.Values(dedupTypes))

	slices.SortFunc(poeDat.Tables, func(a, b poedat.Table) int {
		return strings.Compare(a.Name, b.Name)
	})

	pobData, err := pob.UnmarshalSchema(pobJson)
	if err != nil {
		panic(err)
	}

	pobPoeDat := PobToPoe(pobData)
	mapping := make(map[string]poedat.Table)
	for _, table := range pobPoeDat.Tables {
		mapping[table.Name] = table
	}

	f, err := template.New("foo").Funcs(template.FuncMap{
		"lowercase": func(s string) string {
			return strings.ToLower(s)
		},
		"lowercaseX": func(s *string) *string {
			if s == nil {
				return nil
			}
			out := strings.ToLower(*s)
			return &out
		},
		"diff": func(a poedat.Table, b poedat.Table) TableDiff {
			return TableDiff{
				Source:  a,
				Target:  b,
				ColDiff: Diff(a, b),
			}
		},
		"strneq": func(a *string, b *string) bool {
			if a == nil && b == nil {
				return false
			}

			if a != nil && b != nil {
				return *a != *b
			}

			return true
		},
		"default": func(val *string, def string) string {
			if val == nil || *val == "" {
				return def
			}
			return *val
		},
		"color": func(override DiffType, base DiffType) string {
			if override != DiffTypeUnchanged {
				return DiffTypeColor[override]
			}
			return DiffTypeColor[base]
		},
		"colorn": func(override DiffType) string {
			if override != DiffTypeUnchanged {
				return DiffTypeColor[override]
			}
			return "neutral"
		},
		"colorname": func(val *string) string {
			if val == nil || *val == "" {
				return "neutral"
			}
			return "unset"
		},
		"colorname2": func(override DiffType, base DiffType, val *string) string {
			if val == nil || *val == "" {
				return "neutral"
			}
			if override != DiffTypeUnchanged {
				return DiffTypeColor[override]
			}
			return DiffTypeColor[base]
		},
		"missing": func(source poedat.Table, target poedat.Table) []int {
			count := len(source.Columns) - len(target.Columns)
			if count < 1 {
				return nil
			}

			out := make([]int, count)
			for i := 0; i < count; i++ {
				out[i] = len(target.Columns) + i
			}

			return out
		},
	}).ParseFiles("diff.html.tmpl")
	if err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	if err := f.ExecuteTemplate(buffer, "T", DiffData{
		Source: poeDat,
		Target: mapping,
	}); err != nil {
		panic(err)
	}

	if err := os.WriteFile("out.html", buffer.Bytes(), 0644); err != nil {
		panic(err)
	}
}

var typeMap = map[pob.Type]poedat.Type{
	pob.Bool:     poedat.Bool,
	pob.Enum:     poedat.Enumrow,
	pob.Float:    poedat.F32,
	pob.Int:      poedat.I32,
	pob.Interval: poedat.Type(pob.Interval),
	pob.Key:      poedat.Foreignrow,
	pob.ShortKey: poedat.Row,
	pob.String:   poedat.String,
	pob.UInt:     poedat.U32,
	pob.UInt16:   poedat.U16,
}

func PobToPoe(in pob.Schema) poedat.Schema {
	out := poedat.Schema{}

	for name, val := range in {
		if val.PurpleSchema != nil {
			continue
		}

		cols := val.SchemaElementArray

		columns := make([]poedat.Column, len(cols))

		for i, col := range cols {

			columns[i] = poedat.Column{
				Array: col.List,
				Type:  typeMap[col.Type],
			}

			if col.Name != nil && *col.Name != "" {
				columns[i].Name = col.Name
			}

			if col.RefTo != nil && *col.RefTo != "" {
				columns[i].References = &poedat.References{
					Table: *col.RefTo,
				}
			}
		}

		out.Tables = append(out.Tables, poedat.Table{
			ValidFor: 0,
			Name:     name,
			Columns:  columns,
		})
	}

	return out
}
