package flags

import (
	"flag"
)

type Flags struct {
	pathToConfig      string // Path to .yaml configuration file
	printTransactions bool
	generateReport    bool
}

func New() *Flags {
	return &Flags{}
}

func (f *Flags) Parse() {
	flag.StringVar(&f.pathToConfig, "conf", "myosek.yml", "path to the .yml configuration file")
	flag.BoolVar(&f.printTransactions, "trans", false, "print transactions")
	flag.BoolVar(&f.generateReport, "report", false, "generate report")
	flag.Parse()
}

func (f *Flags) GetPathToConfig() string {
	return f.pathToConfig
}

func (f *Flags) IsPrintTransactions() bool {
	return f.printTransactions
}

func (f *Flags) IsGenerateReport() bool {
	return f.generateReport
}
