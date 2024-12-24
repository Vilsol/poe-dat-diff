package main

import (
	"strings"

	"poe-dat-diff/poedat"
)

type DiffType string

const (
	DiffTypeAdded     = DiffType("Added")
	DiffTypeChanged   = DiffType("Changed")
	DiffTypeRemoved   = DiffType("Removed")
	DiffTypeUnchanged = DiffType("Unchanged")
)

type ColDiff struct {
	Full      DiffType
	Type      DiffType
	Name      DiffType
	Array     DiffType
	Reference DiffType
}

var DiffTypeColor = map[DiffType]string{
	DiffTypeAdded:     "success",
	DiffTypeChanged:   "changed",
	DiffTypeRemoved:   "missing",
	DiffTypeUnchanged: "unset",
}

func Diff(a poedat.Table, b poedat.Table) []ColDiff {
	diffs := make([]ColDiff, len(b.Columns))

	for i, colB := range b.Columns {
		cd := ColDiff{
			Full:      DiffTypeAdded,
			Type:      DiffTypeUnchanged,
			Name:      DiffTypeUnchanged,
			Array:     DiffTypeUnchanged,
			Reference: DiffTypeUnchanged,
		}

		colA := colB
		if len(a.Columns) > i {
			cd.Full = DiffTypeUnchanged
			colA = a.Columns[i]
		}

		if colA.Type != colB.Type {
			cd.Type = DiffTypeChanged
		}

		if colA.Name != nil && colB.Name != nil {
			if strings.ToLower(*colA.Name) != strings.ToLower(*colB.Name) {
				cd.Name = DiffTypeChanged
			}
		} else {
			if colA.Name != nil && colB.Name == nil {
				cd.Name = DiffTypeRemoved
			} else if colA.Name == nil && colB.Name != nil {
				cd.Name = DiffTypeAdded
			}
		}

		if colA.Array != colB.Array {
			cd.Array = DiffTypeChanged
		}

		if colA.References != nil && colB.References != nil {
			if colA.References.Table != colB.References.Table {
				cd.Reference = DiffTypeChanged
			}
		} else {
			if colA.References != nil && colB.References == nil {
				cd.Reference = DiffTypeRemoved
			} else if colA.References == nil && colB.References != nil {
				cd.Reference = DiffTypeAdded
			}
		}

		diffs[i] = cd
	}

	return diffs
}

func strneq(a *string, b *string) bool {
	if a == nil && b == nil {
		return false
	}

	if a != nil && b != nil {
		return *a != *b
	}

	return true
}
