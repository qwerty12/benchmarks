package report

import (
	"github.com/maypok86/benchmarks/simulator/internal/report/chart"
	"github.com/maypok86/benchmarks/simulator/internal/report/simulation"
	"github.com/maypok86/benchmarks/simulator/internal/report/table"
)

type Reporter struct {
	chart *chart.Chart
	table *table.Table
}

func NewReporter(name string, t [][]simulation.Result) *Reporter {
	return &Reporter{
		chart: chart.NewChart(name, t),
		table: table.NewTable(t),
	}
}

func (r *Reporter) Report() {
	if r == nil {
		return
	}

	r.table.Report()
	r.chart.Report()
}
