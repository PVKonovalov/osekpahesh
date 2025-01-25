package report

import (
	"fmt"
	"github.com/dustin/go-humanize"
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
	service := myOsek.Service
	account := myOsek.Account

	// Create a new PDF document
	pdf := &gopdf.GoPdf{}
	// Start the PDF with a custom page size (we'll adjust it later)
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
		TrimBox: gopdf.Box{
			Left:   10,
			Top:    10,
			Right:  10,
			Bottom: 10,
		}})

	var err error

	if err = pdf.AddTTFFont("font1", fmt.Sprintf("./asset/fonts/%s", r.fontName)); err != nil {
		return err
	}

	if err = pdf.SetFont("font1", "", 10); err != nil {
		return err
	}

	// Add a new page to the document
	pdf.AddPage()
	pdf.SetY(0.0)
	pdf.SetFillColor(255, 155, 255)

	// Checking for the existence of the output directory and creating it if necessary
	if _, err = os.Stat(r.pathToReport); os.IsNotExist(err) {
		err = os.MkdirAll(r.pathToReport, 0722)
		if err != nil {
			return err
		}
	}

	err = pdf.Image("./asset/images/stamp.png", 10, 10, nil) //print image
	if err != nil {
		return err
	}

	cellWidth := 200.0

	tableHead := pdf.NewTableLayout(gopdf.PageSizeA4.W-2*cellWidth-20, 10, 20, 4)
	tableHead.AddColumn("", cellWidth, "right")
	tableHead.AddColumn("", cellWidth, "right")

	tableHead.AddRow([]string{reverse("מסמך ממוחשב"), reverse(myOsek.Name)})
	tableHead.AddRow([]string{reverse(fmt.Sprintf("קבלה %s (מקור)", reverse(fmt.Sprintf("%d", transaction.Receipt)))), myOsek.Address})
	tableHead.AddRow([]string{fmt.Sprintf("%s :%s", myOsek.Tz, reverse("עוסק פטור")), ""})
	tableHead.AddRow([]string{fmt.Sprintf("%s :%s", transaction.Date, reverse("תאריך")), fmt.Sprintf("%s :%s", myOsek.Email, reverse("דוא\"ל"))})

	tableHead.SetHeaderStyle(gopdf.CellStyle{
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
	})
	tableHead.SetTableStyle(gopdf.CellStyle{
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
	})

	tableHead.SetCellStyle(gopdf.CellStyle{
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "font1",
		FontSize:  10,
	})

	if err = tableHead.DrawTable(); err != nil {
		return err
	}

	// Client
	pdf.SetY(pdf.GetY() + 20.0)
	pdf.SetX(20)
	if err = pdf.SetFont("font1", "", 11); err != nil {
		return err
	}
	var clientName string
	if len(client[transaction.Client-1].Tz) == 0 {
		clientName = fmt.Sprintf("%s    :%s  ", client[transaction.Client-1].Name, reverse("לכבוד"))
	} else {
		clientName = fmt.Sprintf("%s  :%s    %s    :%s  ", client[transaction.Client-1].Tz, reverse("ת.ז./ע.מ."), client[transaction.Client-1].Name, reverse("לכבוד"))
	}
	pdf.SetGrayFill(0.1)
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

	//
	if err = pdf.SetFont("font1", "", 10); err != nil {
		return err
	}
	tableTransaction := pdf.NewTableLayout(20, pdf.GetY()+40.0, 20, 1)
	tableTransaction.AddColumn(reverse("סה\"כ"), 80, "right")
	tableTransaction.AddColumn(reverse("כמות"), 80, "right")
	tableTransaction.AddColumn(reverse("מחיר יחידה"), 80, "right")
	tableTransaction.AddColumn(reverse("פירוט"), gopdf.PageSizeA4.W-280, "right")

	tableTransaction.AddRow([]string{
		fmt.Sprintf("%s %s", account[transaction.Account-1].Currency, humanize.FormatFloat("#,###.##", transaction.Total)),
		fmt.Sprintf("%d", transaction.Amount),
		fmt.Sprintf("%s %s", account[transaction.Account-1].Currency, humanize.FormatFloat("#,###.##", transaction.Total)),
		service[transaction.Service-1].Name,
	})
	if transaction.Account != 1 {
		tableTransaction.AddRow([]string{fmt.Sprintf("%.4f %s", transaction.Rate, reverse("לפי שער")), reverse("דולר"), reverse("מטבע"), ""})
	}

	tableTransaction.SetHeaderStyle(gopdf.CellStyle{
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		Font:      "font1",
		FontSize:  13,
	})
	tableTransaction.SetTableStyle(gopdf.CellStyle{})

	tableTransaction.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top: true,
		},
		FillColor: gopdf.RGBColor{R: 240, G: 240, B: 240},
		Font:      "font1",
		FontSize:  10,
	})

	if err = tableTransaction.DrawTable(); err != nil {
		return err
	}

	//

	tableGrandTotal := pdf.NewTableLayout(20, pdf.GetY()+40.0, 20, 2)
	tableGrandTotal.AddColumn(reverse("סכום"), 80, "right")
	tableGrandTotal.AddColumn(reverse("תאריך"), 80, "right")
	tableGrandTotal.AddColumn(reverse(""), 180, "right")
	tableGrandTotal.AddColumn(reverse("שולם באמצעות"), gopdf.PageSizeA4.W-380, "right")

	if transaction.Account != 1 {
		tableGrandTotal.AddRow([]string{
			fmt.Sprintf("%s %s", account[transaction.Account-1].Currency, humanize.FormatFloat("#,###.##", transaction.Total)),
			transaction.Date,
			fmt.Sprintf("%s :%s", account[transaction.Account-1].Number, reverse("הופקד לחשבון")),
			reverse("העברה בנקאית")})
	} else {
		tableGrandTotal.AddRow([]string{
			fmt.Sprintf("₪ %s", humanize.FormatFloat("#,###.##", transaction.Total*transaction.Rate)),
			transaction.Date,
			fmt.Sprintf("%s :%s", account[transaction.Account-1].Number, reverse("הופקד לחשבון")),
			reverse("העברה בנקאית")})
	}

	tableGrandTotal.AddRow([]string{
		fmt.Sprintf("₪ %s", humanize.FormatFloat("#,###.##", transaction.Total*transaction.Rate)),
		reverse("סה\"כ שולם"),
		"",
		""})

	tableGrandTotal.SetHeaderStyle(gopdf.CellStyle{
		Font:     "font1",
		FontSize: 13,
	})
	tableGrandTotal.SetTableStyle(gopdf.CellStyle{})

	tableGrandTotal.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{},
		Font:        "font1",
		FontSize:    10,
	})

	if err = tableGrandTotal.DrawTable(); err != nil {
		return err
	}
	//

	err = pdf.Image("./asset/images/sign.png", gopdf.PageSizeA4.W-120, pdf.GetY()+40.0, nil) //print image
	if err != nil {
		return err
	}

	if err = pdf.SetFont("font1", "", 8); err != nil {
		return err
	}
	pdf.SetGrayFill(0.1)
	pdf.SetXY(20, gopdf.PageSizeA4.H-20)
	err = pdf.Cell(nil, fmt.Sprintf("%s - %d %s", reverse("מקור"), transaction.Receipt, reverse("קבלה")))
	if err != nil {
		return err
	}

	pdf.SetXY(gopdf.PageSizeA4.W-80, gopdf.PageSizeA4.H-20)
	err = pdf.Cell(nil, reverse("עמוד 1 מתוך 1"))
	if err != nil {
		return err
	}

	if err = pdf.WritePdf(fmt.Sprintf(filepath.Join(r.pathToReport, r.filenameTemplate), transaction.Receipt)); err != nil {
		return err
	}

	return nil
}
