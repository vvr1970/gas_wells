package exporter

import "fmt"

// internal/pkg/exporter/excel.go
func ExportToExcel(wells []entity.Well, filename string) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "ID")
	f.SetCellValue(sheet, "B1", "Name")

	for i, well := range wells {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), well.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), well.Name)
	}

	return f.SaveAs(filename)
}
