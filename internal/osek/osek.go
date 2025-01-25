package osek

import (
	"fmt"
	table "github.com/PVKonovalov/dyn_table"
	"osekpahesh/internal/configuration"
	"osekpahesh/internal/report"
)

type Osek struct {
	Config configuration.OsekPaHesh
}

func New() *Osek {
	return &Osek{}
}

func (o *Osek) LoadConfiguration(pathToConfig string) error {
	return configuration.ReadConfigFromYMLFile(pathToConfig, &o.Config)
}

func (o *Osek) GetGrandTotal() float64 {
	var grandTotal float64

	for _, transaction := range o.Config.Transaction {
		if transaction.Account != 1 {
			grandTotal += transaction.Rate * transaction.Total
		} else {
			grandTotal += transaction.Total
		}
	}

	return grandTotal
}

func (o *Osek) PrintTransactions() {
	var grandTotal float64

	tab := table.DynTable{
		Width:   []int{4, 8, 12, 10, 7, 10},
		Headers: []string{"#", "Receipt", "Date", "Total", "Rate", "Total, NIS"},
		Align:   []int{table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignLeft, table.AlignRight},
	}

	tab.WriteHeader(nil, 2)

	for idx, transaction := range o.Config.Transaction {
		if transaction.Account != 1 {
			tab.AppendRow([]string{
				fmt.Sprintf("%d", idx+1),
				fmt.Sprintf("%d", transaction.Receipt),
				transaction.Date,
				fmt.Sprintf("%.2f", transaction.Total),
				fmt.Sprintf("%.4f", transaction.Rate),
				fmt.Sprintf("%.2f", transaction.Total*transaction.Rate),
			})
			grandTotal += transaction.Rate * transaction.Total
		} else {
			tab.AppendRow([]string{
				fmt.Sprintf("%d", idx+1),
				fmt.Sprintf("%d", transaction.Receipt),
				transaction.Date,
				"",
				"",
				fmt.Sprintf("%.2f", transaction.Total),
			})
			grandTotal += transaction.Total
		}
	}
	tab.AppendRow([]string{
		"",
		"",
		"",
		"",
		"",
		fmt.Sprintf("%.2f", grandTotal),
	})
}

func (o *Osek) CreateReports() {
	osekReport := report.New(&o.Config)
	for idx := range o.Config.Transaction {
		if err := osekReport.GenerateReport(idx); err != nil {
			fmt.Println(err)
		}
	}
}
