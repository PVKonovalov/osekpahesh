package osek

import (
	"fmt"
	table "github.com/PVKonovalov/dyn_table"
	"osekpahesh/internal/configuration"
	"osekpahesh/internal/currency"
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

func (o *Osek) getGrandTotal() currency.Currency {
	var grandTotal currency.Currency

	for _, transaction := range o.Config.Transaction {
		if transaction.Account != 1 {
			grandTotal.Add(*transaction.Total.Rate(transaction.Rate))
		} else {
			grandTotal.Add(transaction.Total)
		}
	}

	return grandTotal
}

func (o *Osek) PrintTransactions() {
	var grandTotal currency.Currency

	tab := table.DynTable{
		Width:   []int{4, 20, 8, 12, 12, 7, 12},
		Headers: []string{"#", "Client", "Receipt", "Date", "Total", "Rate", "Total, NIS"},
		Align:   []int{table.AlignRight, table.AlignLeft, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignLeft, table.AlignRight},
	}

	tab.WriteHeader(nil, 2)

	for idx, transaction := range o.Config.Transaction {
		rated := transaction.Total.Rate(transaction.Rate)
		if transaction.Account != 1 {
			tab.AppendRow([]string{
				fmt.Sprintf("%d", idx+1),
				fmt.Sprintf("%s", o.Config.Client[transaction.Client].Name),
				fmt.Sprintf("%d", transaction.Receipt),
				transaction.Date,
				fmt.Sprintf("%s %s", o.Config.Osek.Account[transaction.Account].Currency, transaction.Total.String()),
				fmt.Sprintf("%4s", transaction.Rate.String()),
				fmt.Sprintf("%s", rated.String()),
			})
			grandTotal.Add(*rated)
		} else {
			tab.AppendRow([]string{
				fmt.Sprintf("%d", idx+1),
				fmt.Sprintf("%s", o.Config.Client[transaction.Client].Name),
				fmt.Sprintf("%d", transaction.Receipt),
				transaction.Date,
				"",
				"",
				fmt.Sprintf("%s %s", o.Config.Osek.Account[transaction.Account].Currency, transaction.Total.String()),
			})
			grandTotal.Add(transaction.Total)
		}
	}
	tab.AppendRow([]string{
		"",
		"",
		"",
		"",
		"",
		"",
		fmt.Sprintf("%s %s", o.Config.Osek.Account[1].Currency, grandTotal.String()),
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
