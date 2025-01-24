package report

import (
	"fmt"
	"github.com/signintech/gopdf"
	"os"
	"osekpahesh/internal/configuration"
	"path/filepath"
)

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] == ')' {
			runes[i] = '('
		} else if runes[i] == '(' {
			runes[i] = ')'
		}

		if runes[j] == ')' {
			runes[j] = '('
		} else if runes[j] == '(' {
			runes[j] = ')'
		}

		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type ReportPdf struct {
	config           *configuration.OsekPaHesh
	pathToReport     string
	filenameTemplate string
	fontName         string
}

func New(config *configuration.OsekPaHesh) *ReportPdf {
	return &ReportPdf{
		config:           config,
		filenameTemplate: "receipt-%d.pdf",
		pathToReport:     "./reports",
		fontName:         "Arial.ttf",
	}
}

func (r *ReportPdf) GenerateReport(transactionId int) error {
	transaction := r.config.Transaction[transactionId]
	myOsek := r.config.Osek
	client := r.config.Client

	// Create a new PDF document
	pdf := &gopdf.GoPdf{}
	// Start the PDF with a custom page size (we'll adjust it later)
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	// Add a new page to the document
	pdf.AddPage()

	// Checking for the existence of the output directory and creating it if necessary
	if _, err := os.Stat(r.pathToReport); os.IsNotExist(err) {
		err = os.MkdirAll(r.pathToReport, 0722)
		if err != nil {
			return err
		}
	}

	var err error

	if err = pdf.AddTTFFont("font1", fmt.Sprintf("./asset/fonts/%s", r.fontName)); err != nil {
		return err
	}

	if err = pdf.SetFont("font1", "", 10); err != nil {
		return err
	}

	err = pdf.Image("./asset/images/stamp.png", 10, 10, nil) //print image
	if err != nil {
		return err
	}

	cellWidth := 200.0

	tableHead := pdf.NewTableLayout(gopdf.PageSizeA4.W-2*cellWidth-20, 10, 20, 2)
	tableHead.AddColumn("", cellWidth, "right")
	tableHead.AddColumn("", cellWidth, "right")

	tableHead.AddRow([]string{reverse("מסמך ממוחשב"), reverse(myOsek.Name)})
	tableHead.AddRow([]string{reverse(fmt.Sprintf("קבלה %s (מקור)", reverse(fmt.Sprintf("%d", transaction.Receipt)))), myOsek.Address})
	tableHead.AddRow([]string{fmt.Sprintf("%s :%s", myOsek.Tz, reverse("עוסק פטור")), ""})
	tableHead.AddRow([]string{fmt.Sprintf("%s :%s", transaction.Date, reverse("תאריך")), fmt.Sprintf("%s :%s", myOsek.Email, reverse("דוא\"ל"))})

	tableHead.SetHeaderStyle(gopdf.CellStyle{})
	tableHead.SetTableStyle(gopdf.CellStyle{})

	tableHead.SetCellStyle(gopdf.CellStyle{
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "font1",
		FontSize:  10,
	})

	if err = tableHead.DrawTable(); err != nil {
		return err
	}

	pdf.SetY(pdf.GetY() + 20.0)
	pdf.SetX(20)

	var clientName string
	if len(client[transaction.Client-1].Tz) == 0 {
		clientName = fmt.Sprintf("%s    :%s  ", client[transaction.Client-1].Name, reverse("לכבוד"))
	} else {
		clientName = fmt.Sprintf("%s  :%s    %s    :%s  ", client[transaction.Client-1].Tz, reverse("ת.ז./ע.מ."), client[transaction.Client-1].Name, reverse("לכבוד"))
	}
	pdf.SetGrayFill(0.8)
	err = pdf.CellWithOption(&gopdf.Rect{
		W: gopdf.PageSizeA4.W - 40.0,
		H: 30,
	},
		clientName,
		gopdf.CellOption{
			Align:  gopdf.Right | gopdf.Middle,
			Border: gopdf.Left | gopdf.Right | gopdf.Bottom | gopdf.Top,
		})

	if err != nil {
		return err
	}

	if err = pdf.WritePdf(fmt.Sprintf(filepath.Join(r.pathToReport, r.filenameTemplate), transaction.Receipt)); err != nil {
		return err
	}

	return nil
}
