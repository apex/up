package table

import (
	tab "github.com/tj/simpletable"
)

// TODO: add to tj/simpletable

var style = &tab.Style{
  Border: &tab.BorderStyle{
    TopLeft:            "",
    Top:                "",
    TopRight:           "",
    Right:              "",
    BottomRight:        "",
    Bottom:             "",
    BottomLeft:         "",
    Left:               "    ",
    TopIntersection:    "",
    BottomIntersection: "",
  },
  Divider: &tab.DividerStyle{
    Left:         "",
    Center:       "=",
    Right:        "",
    Intersection: " ",
  },
  Cell: "",
}

// Cell is a single cell.
type Cell = tab.Cell

// Row is a group of cells.
type Row []*Cell

// Table is an ascii table.
type Table struct {
	*tab.Table
}

// New table.
func New() *Table {
	t := &Table{tab.New()}
  t.SetStyle(style)
  return t
}

// AddRow adds a row.
func (t *Table) AddRow(r Row)  {
  t.Body.Cells = append(t.Body.Cells, r)
}
