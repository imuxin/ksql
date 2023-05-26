package pretty

import (
	"fmt"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func TestExample(_ *testing.T) {
	t := table.NewWriter()
	// t.SetOutputMirror(os.Stdout) // If set, Render() will output to stdout
	t.AppendHeader(table.Row{"#", "First Name", "Last Name", "Salary"})
	t.AppendRows([]table.Row{
		{1, "Arya", "Stark", 3000},
		{20, "Jon", "Snow", 2000, "You know nothing, Jon Snow!"},
	})
	t.AppendSeparator()
	t.AppendRow([]interface{}{300, "Tyrion", "Lannister", 5000})
	t.AppendFooter(table.Row{"", "", "Total", 10000})

	for _, style := range []table.Style{table.StyleDefault, table.StyleLight, table.StyleColoredBright} {
		t.SetStyle(style)
		t.Render()
	}
}

func TestAutoMerge(_ *testing.T) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"Node IP", "Pods", "Namespace", "Container", "RCE", "RCE"}, rowConfigAutoMerge)
	t.AppendHeader(table.Row{"", "", "", "", "EXE", "RUN"})
	t.AppendRow(table.Row{"1.1.1.1", "Pod 1A", "NS 1A", "C 1", "Y", "Y"}, rowConfigAutoMerge)
	t.AppendRow(table.Row{"1.1.1.1", "Pod 1A", "NS 1A", "C 2", "Y", "N"}, rowConfigAutoMerge)
	t.AppendRow(table.Row{"1.1.1.1", "Pod 1A", "NS 1B", "C 3", "N", "N"}, rowConfigAutoMerge)
	t.AppendRow(table.Row{"1.1.1.1", "Pod 1B", "NS 2", "C 4", "N", "N"}, rowConfigAutoMerge)
	t.AppendRow(table.Row{"1.1.1.1", "Pod 1B", "NS 2", "C 5", "Y", "N"}, rowConfigAutoMerge)
	t.AppendRow(table.Row{"2.2.2.2", "Pod 2", "NS 3", "C 6", "Y", "Y"}, rowConfigAutoMerge)
	t.AppendRow(table.Row{"2.2.2.2", "Pod 2", "NS 3", "C 7", "Y", "Y"}, rowConfigAutoMerge)
	t.AppendFooter(table.Row{"", "", "", 7, 5, 3})
	t.SetAutoIndex(true)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
		{Number: 4, AutoMerge: true},
		{Number: 5, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 6, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
	})
	// t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	fmt.Println(t.Render())
}
