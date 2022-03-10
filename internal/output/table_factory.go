// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import (
	"errors"
	"fmt"
	"github.com/finleap-connect/monoctl/internal/util"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"k8s.io/apimachinery/pkg/util/duration"
	"vbom.ml/util/sortorder"
)

// TableFactory to print a table with sorting
type TableFactory struct {
	sortOrder        SortOrder
	sortIndex        uint
	sortColumn       string
	exportFormat     ExportFormat
	exportFile       string
	header           []string
	data             [][]interface{}
	columnFormatters map[string]func(interface{}) string
}

// NewTableFactory creates a new TableFactory to render a sorted table
func NewTableFactory() *TableFactory {
	tf := new(TableFactory)
	tf.sortIndex = 0
	tf.sortOrder = Ascending
	tf.exportFormat = CSV
	tf.exportFile = ""
	tf.columnFormatters = make(map[string]func(interface{}) string)
	return tf
}

func (tf *TableFactory) findColumn(column string) (int, error) {
	if tf.header == nil {
		return -1, errors.New("header not set")
	}
	for idx, header := range tf.header {
		if strings.EqualFold(header, column) {
			return idx, nil
		}
	}
	return -1, errors.New("column not found")
}

// SetSortOrder sets the SortOrder when rendering the table
func (tf *TableFactory) SetSortOrder(sortOrder SortOrder) *TableFactory {
	tf.sortOrder = sortOrder
	return tf
}

// SetSortIndex sets the index of the column after which the data should be sorted when rendering the table
func (tf *TableFactory) SetSortIndex(sortIndex uint) *TableFactory {
	tf.sortIndex = sortIndex
	return tf
}

// SetSortColumn sets the index of the column after which the data should be sorted when rendering the table by column name
func (tf *TableFactory) SetSortColumn(column string) *TableFactory {
	tf.sortColumn = column
	return tf
}

// SetExportFormat sets the file format in which the data will be written. CSV is set by default
func (tf *TableFactory) SetExportFormat(exportFormat ExportFormat) *TableFactory {
	tf.exportFormat = exportFormat
	return tf
}

// SetExportFile sets the file to which the data will be written
func (tf *TableFactory) SetExportFile(exportFile string) *TableFactory {
	tf.exportFile = exportFile
	return tf
}

// SetHeader sets the header row of the table
func (tf *TableFactory) SetHeader(header []string) *TableFactory {
	tf.header = header
	tf.sortIndex = 0
	return tf
}

// SetData sets the data rows of the table
func (tf *TableFactory) SetData(data [][]interface{}) *TableFactory {
	tf.data = data
	return tf
}

// SetColumnFormatter sets a new formatter for a specific column to render data
func (tf *TableFactory) SetColumnFormatter(column string, columnFormatter func(interface{}) string) *TableFactory {
	tf.columnFormatters[strings.ToLower(column)] = columnFormatter
	return tf
}

// ToTable creates a tablewriter.Table instance with sorted data and default rendering settings
func (tf *TableFactory) ToTable() (*tablewriter.Table, error) {
	if tf.exportFile != "" {
		return tf.newCsvTable()
	}
	return tf.newStdoutTable()
}

func (tf *TableFactory) newStdoutTable() (*tablewriter.Table, error) {
	tbl := tablewriter.NewWriter(os.Stdout)
	tbl.SetAutoWrapText(false)
	tbl.SetAutoFormatHeaders(true)
	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	tbl.SetAlignment(tablewriter.ALIGN_LEFT)
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator("")
	tbl.SetRowSeparator("")
	tbl.SetHeaderLine(false)
	tbl.SetBorder(false)
	tbl.SetTablePadding("\t") // pad with tabs
	tbl.SetNoWhiteSpace(true)
	tbl.SetHeader(tf.header)
	tbl.AppendBulk(tf.formatData())
	return tbl, nil
}

func (tf *TableFactory) newCsvTable() (*tablewriter.Table, error) {
	file, err := util.NewFileSafe(tf.exportFile)
	if err != nil {
		return nil, errors.New("failed to export. Please ensure file doesn't exits or try another path")
	}

	tbl := tablewriter.NewWriter(file)
	tbl.SetAutoWrapText(false)
	tbl.SetAutoFormatHeaders(true)
	tbl.SetAlignment(tablewriter.ALIGN_LEFT)
	tbl.SetHeaderLine(false)
	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	tbl.SetAutoMergeCells(true)
	tbl.SetColumnSeparator(",")
	tbl.SetBorder(false)
	tbl.SetHeader(tf.header)
	tbl.AppendBulk(tf.formatData())
	return tbl, nil
}

// Len implements sort.Sorter interface
func (tf *TableFactory) Len() int {
	return len(tf.data)
}

// Swap implements sort.Sorter interface
func (tf *TableFactory) Swap(i, j int) {
	tf.data[i], tf.data[j] = tf.data[j], tf.data[i]
}

func isLess(i, j reflect.Value) bool {
	switch i.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return i.Int() < j.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return i.Uint() < j.Uint()
	case reflect.Float32, reflect.Float64:
		return i.Float() < j.Float()
	case reflect.String:
		return sortorder.NaturalLess(i.String(), j.String())
	case reflect.Ptr:
		return isLess(i.Elem(), j.Elem())
	case reflect.Struct:
		in := i.Interface()
		if t, ok := in.(time.Time); ok {
			time := j.Interface().(time.Time)
			return t.Before(time)
		}
		return false
	default:
		return false
	}
}

// Less implements sort.Sorter interface
func (tf *TableFactory) Less(i, j int) bool {
	iData := reflect.ValueOf(tf.data[i][tf.sortIndex])
	jData := reflect.ValueOf(tf.data[j][tf.sortIndex])

	if tf.sortOrder == Descending {
		return isLess(jData, iData)
	}
	return isLess(iData, jData)
}

func (tf *TableFactory) formatData() [][]string {
	// Find column to sort after
	idx, err := tf.findColumn(tf.sortColumn)
	if err != nil {
		tf.sortIndex = 0
	} else {
		tf.sortIndex = uint(idx)
	}
	sort.Sort(tf)

	result := make([][]string, len(tf.data))
	for rowIdx, row := range tf.data {
		if result[rowIdx] == nil {
			result[rowIdx] = make([]string, len(tf.data[rowIdx]))
		}
		for columnIdx, value := range row {
			if formatter, ok := tf.columnFormatters[strings.ToLower(tf.header[columnIdx])]; ok {
				result[rowIdx][columnIdx] = formatter(value)
			} else {
				result[rowIdx][columnIdx] = toString(reflect.ValueOf(value))
			}
		}
	}
	return result
}

func toString(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%v", val.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", val.Float())
	case reflect.Ptr:
		return toString(val.Elem())
	case reflect.String:
		return val.String()
	case reflect.Struct:
		in := val.Interface()
		if t, ok := in.(time.Duration); ok {
			return duration.HumanDuration(t)
		}
	}
	return val.String()
}

func DefaultAgeColumnFormatter() func(i interface{}) string {
	return func(i interface{}) string {
		return duration.HumanDuration(i.(time.Duration))
	}
}
