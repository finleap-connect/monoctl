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
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Internal/Util/TableFactory", func() {
	tf := NewTableFactory()
	Expect(tf).ToNot(BeNil())

	tf.SetHeader([]string{"NAME", "VALUE"})
	tf.SetData([][]interface{}{
		{"z", 1},
		{"a", 26},
		{"h", 8},
	})

	It("can create table with default sorting ascending", func() {
		tableData := tf.formatData()
		Expect(tableData[0][0]).To(Equal("a"))
		Expect(tableData[1][0]).To(Equal("h"))
		Expect(tableData[2][0]).To(Equal("z"))
		tf.ToTable().Render()
	})

	It("can create table with default sorting descending", func() {
		tf.SetSortOrder(Descending)
		tableData := tf.formatData()
		Expect(tableData[0][0]).To(Equal("z"))
		Expect(tableData[1][0]).To(Equal("h"))
		Expect(tableData[2][0]).To(Equal("a"))
		tf.ToTable().Render()
	})

	It("can create table with sorting a specific column ascending", func() {
		tf.SetSortOrder(Ascending).SetSortColumn("value")

		tableData := tf.formatData()
		Expect(tableData[0][1]).To(Equal("1"))
		Expect(tableData[1][1]).To(Equal("8"))
		Expect(tableData[2][1]).To(Equal("26"))
		tf.ToTable().Render()
	})

	It("can create table with sorting a specific column descending", func() {
		tf.SetSortOrder(Descending).SetSortColumn("value")

		tableData := tf.formatData()
		Expect(tableData[0][1]).To(Equal("26"))
		Expect(tableData[1][1]).To(Equal("8"))
		Expect(tableData[2][1]).To(Equal("1"))
		tf.ToTable().Render()
	})

	It("can create table with a special formatter", func() {
		tf.SetColumnFormatter("value", func(i interface{}) string {
			return fmt.Sprintf("%02d", i.(int))
		})
		tableData := tf.formatData()
		Expect(tableData[0][1]).To(Equal("26"))
		Expect(tableData[1][1]).To(Equal("08"))
		Expect(tableData[2][1]).To(Equal("01"))
		tf.ToTable().Render()
	})
})
