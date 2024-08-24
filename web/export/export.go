package export

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"kellnhofer.com/work-log/pkg/loc"
	vm "kellnhofer.com/work-log/web/model"
)

type writerToAdapter struct {
	f *excelize.File
}

func (wta *writerToAdapter) WriteTo(w io.Writer) (n int64, err error) {
	return wta.f.WriteTo(w)
}

// --- Base exporter ---

type export struct {
	file   *excelize.File
	styles exportStyles
	sheet  string
}

type exportStyles struct {
	base                    int
	title                   int
	textBold                int
	tableHeader             int
	tableBody               int
	tableBodyAlignmentRight int
}

type exporter struct {
}

func (e *exporter) createFile() *excelize.File {
	return excelize.NewFile()
}

func (e *exporter) createNewExport() *export {
	// Create file
	file := e.createFile()

	// Create styles
	styleBase := e.createBaseStyle(file)
	styleTitle, styleTextBold := e.createTextStyles(file)
	styleTableHeader, styleTableBody, styleTableBodyAlignmentRight := e.createTableStyles(file)

	// Create sheet
	sheet := createString("exportSheetName")

	// Create new export
	return &export{
		file: file,
		styles: exportStyles{
			base:                    styleBase,
			title:                   styleTitle,
			textBold:                styleTextBold,
			tableHeader:             styleTableHeader,
			tableBody:               styleTableBody,
			tableBodyAlignmentRight: styleTableBodyAlignmentRight,
		},
		sheet: sheet,
	}
}

func (e *exporter) createBaseStyle(file *excelize.File) int {
	style, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "left", WrapText: true},
		Font:      &excelize.Font{Size: 10},
	})
	return style
}

func (e *exporter) createTextStyles(file *excelize.File) (int, int) {
	styleTitle, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 14, Bold: true}})
	styleTextBold, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 10, Bold: true}})
	return styleTitle, styleTextBold
}

func (e *exporter) createTableStyles(file *excelize.File) (int, int, int) {
	borderLeft := excelize.Border{Type: "left", Style: 1, Color: "000000"}
	borderRight := excelize.Border{Type: "right", Style: 1, Color: "000000"}
	borderTop := excelize.Border{Type: "top", Style: 1, Color: "000000"}
	borderBottom := excelize.Border{Type: "bottom", Style: 1, Color: "000000"}
	fill := excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"EFEFEF"}}
	styleTableHeader, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "left", WrapText: true},
		Font:      &excelize.Font{Size: 10, Bold: true},
		Border:    []excelize.Border{borderLeft, borderRight, borderTop, borderBottom},
		Fill:      fill,
	})
	styleTableBody, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "left", WrapText: true},
		Font:      &excelize.Font{Size: 10},
		Border:    []excelize.Border{borderLeft, borderRight, borderTop, borderBottom},
	})
	styleTableBodyAlignmentRight, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "right", WrapText: true},
		Font:      &excelize.Font{Size: 10},
		Border:    []excelize.Border{borderLeft, borderRight, borderTop, borderBottom},
	})
	return styleTableHeader, styleTableBody, styleTableBodyAlignmentRight
}

func (e *exporter) configureDocProps(exp *export) {
	now := time.Now()
	exp.file.SetDocProps(&excelize.DocProperties{
		Created:        now.Format(time.RFC3339),
		Creator:        createString("appName"),
		Modified:       now.Format(time.RFC3339),
		LastModifiedBy: createString("appName"),
		Description:    createString("exportPropDescription", createString("appName")),
		Language:       loc.LngTag.String(),
	})
}

func (e *exporter) createWriterTo(exp *export) io.WriterTo {
	return &writerToAdapter{
		f: exp.file,
	}
}

// --- Entries exporter ---

// EntriesExporter exports the data to an Excel file.
type EntriesExporter struct {
	exporter
}

// NewEntriesExporter creates a new entries exporter.
func NewEntriesExporter() *EntriesExporter {
	return &EntriesExporter{}
}

// ExportEntries creates the Excel file for the supplied data and returns it as an io.WriterTo that
// can be used to write the file to a writer.
func (e *EntriesExporter) ExportEntries(filterDetails vm.EntryFilterDetails,
	entries *vm.ListEntries) io.WriterTo {
	exp := e.createNewExport()

	// Configure document properties
	e.configureDocProps(exp)

	// Configure work sheet
	e.configureWorkSheet(exp)

	// Write title
	e.writeTitle(exp, filterDetails)
	// Write entries
	e.writeEntries(exp, entries)

	return e.createWriterTo(exp)
}

func (e *EntriesExporter) configureWorkSheet(exp *export) {
	f := exp.file
	sheet := exp.sheet
	styles := exp.styles

	f.SetColWidth(sheet, "A", "A", 12)
	f.SetColWidth(sheet, "B", "B", 10.5)
	f.SetColWidth(sheet, "C", "E", 7.5)
	f.SetColWidth(sheet, "F", "F", 16.5)
	f.SetColWidth(sheet, "G", "G", 16.5)
	f.SetColWidth(sheet, "H", "H", 42)
	f.SetColStyle(sheet, "A:H", styles.base)
}

func (e *EntriesExporter) writeTitle(exp *export, filterDetails vm.EntryFilterDetails) {
	f := exp.file
	sheet := exp.sheet
	styles := exp.styles

	f.MergeCell(sheet, "A1", "H1")
	f.MergeCell(sheet, "A2", "H2")
	f.MergeCell(sheet, "A3", "H3")
	f.SetCellValue(sheet, "A1", createString("entryExportTitle", createString("appName")))
	f.SetCellValue(sheet, "A2", e.buildExportDetailsString(filterDetails))
	f.SetCellStyle(sheet, "A1", "A1", styles.title)
	f.SetCellStyle(sheet, "A2", "A2", styles.textBold)
}

func (e *EntriesExporter) buildExportDetailsString(filterDetails vm.EntryFilterDetails) string {
	if filterDetails == nil {
		return ""
	}

	switch fd := filterDetails.(type) {
	case *vm.BasicEntryFilterDetails:
		return createString("formLabelText") + " " + createQuotedString(fd.Text)
	case *vm.AdvancedEntryFilterDetails:
		var details []string
		if fd.ByType {
			details = append(details, createString("formLabelType")+" "+fd.Type)
		}
		if fd.ByDate {
			details = append(details, createString("formLabelDate")+" "+fd.Date)
		}
		if fd.ByActivity {
			details = append(details, createString("formLabelActivity")+" "+fd.Activity)
		}
		if fd.ByProject {
			project := createQuotedString(fd.Project)
			details = append(details, createString("formLabelProject")+" "+project)
		}
		if fd.ByDescription {
			description := createQuotedString(fd.Description)
			details = append(details, createString("formLabelDescription")+" "+description)
		}
		if fd.ByLabels {
			labels := strings.Join(createQuotedStrings(fd.Labels), ", ")
			details = append(details, createString("formLabelLabels")+" "+labels)
		}
		return strings.Join(details, " | ")
	default:
		return ""
	}
}

func (e *EntriesExporter) writeEntries(exp *export, entries *vm.ListEntries) {
	f := exp.file
	sheet := exp.sheet
	styles := exp.styles

	// Create heading
	f.MergeCell(sheet, "A4", "H4")
	f.SetCellValue(sheet, "A4", createString("entryExportHeadingEntries"))
	f.SetCellStyle(sheet, "A4", "A4", styles.textBold)

	// Create table header
	f.SetCellValue(sheet, "A5", createString("tableColDate"))
	f.SetCellValue(sheet, "B5", createString("tableColType"))
	f.SetCellValue(sheet, "C5", createString("tableColStart"))
	f.SetCellValue(sheet, "D5", createString("tableColEnd"))
	f.SetCellValue(sheet, "E5", createString("tableColNet"))
	f.SetCellValue(sheet, "F5", createString("tableColActivity"))
	f.SetCellValue(sheet, "G5", createString("tableColProject"))
	f.SetCellValue(sheet, "H5", createString("tableColDescription"))
	f.SetCellStyle(sheet, "A5", "E5", styles.tableHeader)
	f.SetCellStyle(sheet, "F5", "H5", styles.tableHeader)

	// Create table body
	startRow := 6
	curRow := startRow
	for _, day := range entries.Days {
		for _, entry := range day.Entries {
			f.SetCellValue(sheet, getCellName("A", curRow), day.Date)
			f.SetCellValue(sheet, getCellName("B", curRow), entry.EntryType)
			f.SetCellValue(sheet, getCellName("C", curRow), entry.StartTime)
			f.SetCellValue(sheet, getCellName("D", curRow), entry.EndTime)
			f.SetCellValue(sheet, getCellName("E", curRow), entry.Duration)
			f.SetCellValue(sheet, getCellName("F", curRow), entry.EntryActivity)
			f.SetCellValue(sheet, getCellName("G", curRow), entry.Project)
			f.SetCellValue(sheet, getCellName("H", curRow), entry.Description)
			curRow++
		}
	}
	f.SetCellStyle(sheet, getCellName("A", startRow), getCellName("E", curRow-1), styles.tableBody)
	f.SetCellStyle(sheet, getCellName("F", startRow), getCellName("H", curRow-1), styles.tableBody)
}

// --- Overview exporter ---

// OverviewExporter exports the overview data to an Excel file.
type OverviewExporter struct {
	exporter
}

// NewOverviewExporter creates a new overview exporter.
func NewOverviewExporter() *OverviewExporter {
	return &OverviewExporter{}
}

// ExportOverviewEntries creates the Excel file for the supplied data and returns it as an
// io.WriterTo that can be used to write the file to a writer.
func (e *OverviewExporter) ExportOverviewEntries(overviewEntries *vm.OverviewEntries) io.WriterTo {
	exp := e.createNewExport()

	// Configure document properties
	e.configureDocProps(exp)

	// Configure work sheet
	e.configureWorkSheet(exp)

	// Write title
	e.writeTitle(exp, overviewEntries)
	// Write summary
	e.writeSummary(exp, overviewEntries)
	// Write entries
	e.writeEntries(exp, overviewEntries)

	return e.createWriterTo(exp)
}

func (e *OverviewExporter) configureWorkSheet(exp *export) {
	f := exp.file
	sheet := exp.sheet
	styles := exp.styles

	f.SetColWidth(sheet, "A", "A", 12)
	f.SetColWidth(sheet, "B", "B", 10.5)
	f.SetColWidth(sheet, "C", "E", 7.5)
	f.SetColWidth(sheet, "F", "F", 16.5)
	f.SetColWidth(sheet, "G", "G", 16.5)
	f.SetColWidth(sheet, "H", "H", 42)
	f.SetColStyle(sheet, "A:H", styles.base)
}

func (e *OverviewExporter) writeTitle(exp *export, overviewEntries *vm.OverviewEntries) {
	f := exp.file
	sheet := exp.sheet
	styles := exp.styles

	f.MergeCell(sheet, "A1", "H1")
	f.MergeCell(sheet, "A2", "H2")
	f.MergeCell(sheet, "A3", "H3")
	f.SetCellValue(sheet, "A1", createString("overviewExportTitle", createString("appName")))
	f.SetCellValue(sheet, "A2", overviewEntries.CurrMonthName)
	f.SetCellStyle(sheet, "A1", "A1", styles.title)
	f.SetCellStyle(sheet, "A2", "A2", styles.textBold)
}

func (e *OverviewExporter) writeSummary(exp *export, overviewEntries *vm.OverviewEntries) {
	f := exp.file
	sheet := exp.sheet
	styles := exp.styles

	// Prepare cells
	f.MergeCell(sheet, "A4", "H4")
	f.MergeCell(sheet, "B5", "C5")
	f.MergeCell(sheet, "D5", "H5")
	f.MergeCell(sheet, "B6", "C6")
	f.MergeCell(sheet, "D6", "H6")
	f.MergeCell(sheet, "B7", "C7")
	f.MergeCell(sheet, "D7", "H7")
	f.MergeCell(sheet, "A8", "H8")
	f.MergeCell(sheet, "B9", "C9")
	f.MergeCell(sheet, "D9", "H9")
	f.MergeCell(sheet, "B10", "C10")
	f.MergeCell(sheet, "D10", "H10")
	f.MergeCell(sheet, "B11", "C11")
	f.MergeCell(sheet, "D11", "H11")
	f.MergeCell(sheet, "B12", "C12")
	f.MergeCell(sheet, "D12", "H12")
	f.MergeCell(sheet, "B13", "C13")
	f.MergeCell(sheet, "D13", "H13")
	f.MergeCell(sheet, "B14", "C14")
	f.MergeCell(sheet, "D14", "H14")
	f.MergeCell(sheet, "A15", "H15")
	f.MergeCell(sheet, "D15", "H15")

	// Create heading
	f.SetCellValue(sheet, "A4", createString("overviewExportHeadingSummary"))
	f.SetCellStyle(sheet, "A4", "A4", styles.textBold)

	// Create target/actual table
	f.SetCellValue(sheet, "A5", createString("overviewExportSummaryLabelTarget"))
	f.SetCellValue(sheet, "A6", createString("overviewExportSummaryLabelActual"))
	f.SetCellValue(sheet, "A7", createString("overviewExportSummaryLabelBalance"))
	f.SetCellValue(sheet, "B5", overviewEntries.Summary.MonthTargetHours)
	f.SetCellValue(sheet, "B6", overviewEntries.Summary.MonthActualHours)
	f.SetCellValue(sheet, "B7", overviewEntries.Summary.MonthBalanceHours)
	f.SetCellStyle(sheet, "A5", "A7", styles.tableHeader)
	f.SetCellStyle(sheet, "B5", "C7", styles.tableBodyAlignmentRight)

	// Create types table
	f.SetCellValue(sheet, "A9", createString("entryTypeWork"))
	f.SetCellValue(sheet, "A10", createString("entryTypeTravel"))
	f.SetCellValue(sheet, "A11", createString("entryTypeVacation"))
	f.SetCellValue(sheet, "A12", createString("entryTypeHoliday"))
	f.SetCellValue(sheet, "A13", createString("entryTypeIllness"))
	f.SetCellValue(sheet, "B9", overviewEntries.Summary.TypeHours[vm.EntryTypeIdWork])
	f.SetCellValue(sheet, "B10", overviewEntries.Summary.TypeHours[vm.EntryTypeIdTravel])
	f.SetCellValue(sheet, "B11", overviewEntries.Summary.TypeHours[vm.EntryTypeIdVacation])
	f.SetCellValue(sheet, "B12", overviewEntries.Summary.TypeHours[vm.EntryTypeIdHoliday])
	f.SetCellValue(sheet, "B13", overviewEntries.Summary.TypeHours[vm.EntryTypeIdIllness])
	f.SetCellValue(sheet, "B14", overviewEntries.Summary.MonthActualHours)
	f.SetCellStyle(sheet, "A9", "A14", styles.tableHeader)
	f.SetCellStyle(sheet, "B9", "C14", styles.tableBodyAlignmentRight)
}

func (e *OverviewExporter) writeEntries(exp *export, overviewEntries *vm.OverviewEntries) {
	f := exp.file
	sheet := exp.sheet
	styles := exp.styles

	// Create heading
	f.MergeCell(sheet, "A16", "H16")
	f.SetCellValue(sheet, "A16", createString("overviewExportHeadingEntries"))
	f.SetCellStyle(sheet, "A16", "A16", styles.textBold)

	// Create table header
	f.SetCellValue(sheet, "A17", createString("tableColDate"))
	f.SetCellValue(sheet, "B17", createString("tableColType"))
	f.SetCellValue(sheet, "C17", createString("tableColStart"))
	f.SetCellValue(sheet, "D17", createString("tableColEnd"))
	f.SetCellValue(sheet, "E17", createString("tableColNet"))
	f.SetCellValue(sheet, "F17", createString("tableColActivity"))
	f.SetCellValue(sheet, "G17", createString("tableColProject"))
	f.SetCellValue(sheet, "H17", createString("tableColDescription"))
	f.SetCellStyle(sheet, "A17", "E17", styles.tableHeader)
	f.SetCellStyle(sheet, "F17", "H17", styles.tableHeader)

	// Create table body
	startRow := 18
	curRow := startRow
	for _, day := range overviewEntries.EntriesDays {
		f.SetCellValue(sheet, getCellName("A", curRow), day.Weekday+" "+day.Date)
		if len(day.Entries) == 0 {
			f.SetCellValue(sheet, getCellName("B", curRow), "-")
			f.SetCellValue(sheet, getCellName("C", curRow), "-")
			f.SetCellValue(sheet, getCellName("D", curRow), "-")
			f.SetCellValue(sheet, getCellName("E", curRow), "-")
			curRow++
		} else {
			for _, entry := range day.Entries {
				if entry.TypeId == 0 {
					continue
				}
				f.SetCellValue(sheet, getCellName("B", curRow), entry.Type)
				f.SetCellValue(sheet, getCellName("C", curRow), entry.StartTime)
				f.SetCellValue(sheet, getCellName("D", curRow), entry.EndTime)
				f.SetCellValue(sheet, getCellName("E", curRow), entry.Duration)
				f.SetCellValue(sheet, getCellName("F", curRow), entry.Activity)
				f.SetCellValue(sheet, getCellName("G", curRow), entry.Project)
				f.SetCellValue(sheet, getCellName("H", curRow), entry.Description)
				curRow++
			}
		}
		if len(day.Entries) > 1 {
			f.SetCellValue(sheet, getCellName("E", curRow), day.Hours)
			curRow++
		}
	}
	f.SetCellStyle(sheet, getCellName("A", startRow), getCellName("E", curRow-1), styles.tableBody)
	f.SetCellStyle(sheet, getCellName("F", startRow), getCellName("H", curRow-1), styles.tableBody)
}

// --- Helper functions ---

func createString(key string, args ...any) string {
	return loc.CreateString(key, args...)
}

func createQuotedString(value string) string {
	return "\"" + value + "\""
}

func createQuotedStrings(values []string) []string {
	quotedValues := make([]string, len(values))
	for i, value := range values {
		quotedValues[i] = createQuotedString(value)
	}
	return quotedValues
}

func getCellName(col string, row int) string {
	return col + strconv.Itoa(row)
}
