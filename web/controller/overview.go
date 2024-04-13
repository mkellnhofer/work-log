package controller

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view/pages"
)

type overviewFormInput struct {
	month string
}

// OverviewController handles requests for overview endpoints.
type OverviewController struct {
	baseController

	mapper *mapper.OverviewMapper
}

// NewOverviewController creates a new overview controller.
func NewOverviewController(uServ *service.UserService, eServ *service.EntryService,
) *OverviewController {
	return &OverviewController{
		baseController: baseController{
			uServ: uServ,
			eServ: eServ,
		},
		mapper: mapper.NewOverviewMapper(),
	}
}

// --- Endpoints ---

// GetOverviewHandler returns a handler for "GET /overview".
func (c *OverviewController) GetOverviewHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /overview.")
		return c.handleShowOverview(eCtx)
	}
}

// PostOverviewHandler returns a handler for "POST /overview".
func (c *OverviewController) PostOverviewHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /overview.")
		return c.handleExecuteOverviewChange(eCtx)
	}
}

// GetOverviewExportHandler returns a handler for "GET /overview/export".
func (c *OverviewController) GetOverviewExportHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /overview/export.")
		return c.handleExportOverview(eCtx)
	}
}

// --- Handler functions ---

func (c *OverviewController) handleShowOverview(eCtx echo.Context) error {
	// Get view data
	userModel, model, err := c.getOverviewViewData(eCtx)
	if err != nil {
		return err
	}

	// Render
	return web.Render(eCtx, http.StatusOK, pages.OverviewEntriesPage(userModel, model))
}

func (c *OverviewController) handleExportOverview(eCtx echo.Context) error {
	// Get view data
	_, model, err := c.getOverviewViewData(eCtx)
	if err != nil {
		return err
	}

	// Create file
	fileName := fmt.Sprintf("work-log-export-%s.xlsx", model.CurrMonth)
	file := c.exportOverviewEntries(model)

	// Write file
	return c.writeFile(eCtx.Response(), fileName, file)
}

func (c *OverviewController) getOverviewViewData(eCtx echo.Context) (*vm.UserInfo,
	*vm.OverviewEntries, error) {
	// Get context
	ctx := getContext(eCtx)

	// Get current user and user contract
	user, userContract, err := c.getUserAndUserContract(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Get year and month
	year, month, err := c.getOverviewParams(eCtx)
	if err != nil {
		return nil, nil, err
	}

	// Get entries
	entries, err := c.eServ.GetMonthEntriesByUserId(ctx, user.Id, year, month)
	if err != nil {
		return nil, nil, err
	}
	// Get entry master data
	entryTypesMap, entryActivitiesMap, err := c.getEntryMasterDataMap(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Create view model
	userModel := c.mapper.CreateUserInfoViewModel(user)
	prevUrl := getPreviousUrl(eCtx)
	model := c.mapper.CreateOverviewEntriesViewModel(prevUrl, year, month, userContract, entries,
		entryTypesMap, entryActivitiesMap)

	return userModel, model, nil
}

func (c *OverviewController) getOverviewParams(eCtx echo.Context) (int, int, error) {
	// Get year and month
	y, m, avail, err := getMonthQueryParam(eCtx)
	if err != nil {
		return 0, 0, err
	}

	// Was a year and month provided?
	if !avail {
		// Get current year/month
		t := time.Now()
		return t.Year(), int(t.Month()), nil
	} else {
		// Use these
		return y, m, nil
	}
}

func (c *OverviewController) handleExecuteOverviewChange(eCtx echo.Context) error {
	// Get form inputs
	input := c.getOverviewFormInput(eCtx)

	// Validate month param
	_, _, _, err := parseMonth(input.month)
	if err != nil {
		return err
	}

	// Redirect
	return eCtx.Redirect(http.StatusFound, "/overview?month="+input.month)
}

// --- Form input retrieval functions ---

func (c *OverviewController) getOverviewFormInput(eCtx echo.Context) *overviewFormInput {
	i := overviewFormInput{}
	i.month = eCtx.FormValue("month")
	return &i
}

// --- Export functions ---

func (c *OverviewController) exportOverviewEntries(overviewEntries *vm.OverviewEntries,
) *excelize.File {
	f := excelize.NewFile()

	// Configure work book
	now := time.Now()
	f.SetDocProps(&excelize.DocProperties{
		Created:        now.Format(time.RFC3339),
		Creator:        loc.CreateString("appName"),
		Modified:       now.Format(time.RFC3339),
		LastModifiedBy: loc.CreateString("appName"),
		Title: loc.CreateString("exportPropTitle", loc.CreateString("appName"),
			overviewEntries.CurrMonthName),
		Description: loc.CreateString("exportPropDescription", loc.CreateString("appName")),
		Language:    loc.LngTag.String(),
	})

	sheet := "Sheet1"

	// Create default style
	styleDefault, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "left", WrapText: true},
		Font:      &excelize.Font{Size: 10},
	})

	// Create text styles
	styleTitle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 14, Bold: true}})
	styleTextBold, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 10, Bold: true}})

	// Creat tables styles
	borderLeft := excelize.Border{Type: "left", Style: 1, Color: "000000"}
	borderRight := excelize.Border{Type: "right", Style: 1, Color: "000000"}
	borderTop := excelize.Border{Type: "top", Style: 1, Color: "000000"}
	borderBottom := excelize.Border{Type: "bottom", Style: 1, Color: "000000"}
	fill := excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"EFEFEF"}}
	styleTableHeader, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "left", WrapText: true},
		Font:      &excelize.Font{Size: 10, Bold: true},
		Border:    []excelize.Border{borderLeft, borderRight, borderTop, borderBottom},
		Fill:      fill,
	})
	styleTableBody, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "left", WrapText: true},
		Font:      &excelize.Font{Size: 10},
		Border:    []excelize.Border{borderLeft, borderRight, borderTop, borderBottom},
	})
	styleTableBodyAlignmentRight, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top", Horizontal: "right", WrapText: true},
		Font:      &excelize.Font{Size: 10},
		Border:    []excelize.Border{borderLeft, borderRight, borderTop, borderBottom},
	})

	// Configure work sheet
	f.SetColWidth(sheet, "A", "A", 12)
	f.SetColWidth(sheet, "B", "B", 10.5)
	f.SetColWidth(sheet, "C", "F", 7.5)
	f.SetColWidth(sheet, "G", "G", 16.5)
	f.SetColWidth(sheet, "H", "H", 42)
	f.SetColStyle(sheet, "A:H", styleDefault)

	// Write title
	f.MergeCell(sheet, "A1", "H1")
	f.MergeCell(sheet, "A2", "H2")
	f.MergeCell(sheet, "A3", "H3")
	f.SetCellValue(sheet, "A1", loc.CreateString("exportTitle", loc.CreateString("appName")))
	f.SetCellValue(sheet, "A2", overviewEntries.CurrMonthName)
	f.SetCellStyle(sheet, "A1", "A1", styleTitle)
	f.SetCellStyle(sheet, "A2", "A2", styleTextBold)

	// Write summary
	f.MergeCell(sheet, "A4", "H4")
	f.MergeCell(sheet, "B5", "C5")
	f.MergeCell(sheet, "E5", "F5")
	f.MergeCell(sheet, "B6", "C6")
	f.MergeCell(sheet, "E6", "F6")
	f.MergeCell(sheet, "B7", "C7")
	f.MergeCell(sheet, "E7", "F7")
	f.MergeCell(sheet, "B8", "C8")
	f.MergeCell(sheet, "E8", "F8")
	f.MergeCell(sheet, "B9", "C9")
	f.MergeCell(sheet, "E9", "F9")
	f.MergeCell(sheet, "B10", "C10")
	f.MergeCell(sheet, "E10", "F10")
	f.MergeCell(sheet, "A11", "H11")
	f.MergeCell(sheet, "E11", "F11")
	// Create heading
	f.SetCellValue(sheet, "A4", loc.CreateString("overviewHeadingSummary"))
	f.SetCellStyle(sheet, "A4", "A4", styleTextBold)
	// Create target/actual table
	f.SetCellValue(sheet, "A5", loc.CreateString("overviewSummaryLabelTargetHours"))
	f.SetCellValue(sheet, "A6", loc.CreateString("overviewSummaryLabelActualHours"))
	f.SetCellValue(sheet, "A7", loc.CreateString("overviewSummaryLabelBalanceHours"))
	f.SetCellValue(sheet, "B5", overviewEntries.Summary.TargetHours)
	f.SetCellValue(sheet, "B6", overviewEntries.Summary.ActualHours)
	f.SetCellValue(sheet, "B7", overviewEntries.Summary.BalanceHours)
	f.SetCellStyle(sheet, "A5", "A10", styleTableHeader)
	f.SetCellStyle(sheet, "B5", "C10", styleTableBodyAlignmentRight)
	// Create types table
	f.SetCellValue(sheet, "E5", loc.CreateString("entryTypeWork"))
	f.SetCellValue(sheet, "E6", loc.CreateString("entryTypeTravel"))
	f.SetCellValue(sheet, "E7", loc.CreateString("entryTypeVacation"))
	f.SetCellValue(sheet, "E8", loc.CreateString("entryTypeHoliday"))
	f.SetCellValue(sheet, "E9", loc.CreateString("entryTypeIllness"))
	f.SetCellValue(sheet, "G5", overviewEntries.Summary.ActualWorkHours)
	f.SetCellValue(sheet, "G6", overviewEntries.Summary.ActualTravelHours)
	f.SetCellValue(sheet, "G7", overviewEntries.Summary.ActualVacationHours)
	f.SetCellValue(sheet, "G8", overviewEntries.Summary.ActualHolidayHours)
	f.SetCellValue(sheet, "G9", overviewEntries.Summary.ActualIllnessHours)
	f.SetCellValue(sheet, "G10", overviewEntries.Summary.ActualHours)
	f.SetCellStyle(sheet, "E5", "E10", styleTableHeader)
	f.SetCellStyle(sheet, "G5", "G10", styleTableBodyAlignmentRight)

	// Write entries
	// Create heading
	f.MergeCell(sheet, "A12", "H12")
	f.SetCellValue(sheet, "A12", loc.CreateString("overviewHeadingEntries"))
	f.SetCellStyle(sheet, "A12", "A12", styleTextBold)
	// Create table header
	f.SetCellValue(sheet, "A13", loc.CreateString("tableColDate"))
	f.SetCellValue(sheet, "B13", loc.CreateString("tableColType"))
	f.SetCellValue(sheet, "C13", loc.CreateString("tableColStart"))
	f.SetCellValue(sheet, "D13", loc.CreateString("tableColEnd"))
	f.SetCellValue(sheet, "E13", loc.CreateString("tableColNet"))
	f.SetCellValue(sheet, "F13", loc.CreateString("tableColActivity"))
	f.SetCellValue(sheet, "G13", loc.CreateString("tableColDescription"))
	f.SetCellStyle(sheet, "A13", "E13", styleTableHeader)
	f.SetCellStyle(sheet, "F13", "G13", styleTableHeader)
	// Create table body
	row := 14
	for _, day := range overviewEntries.Days {
		f.SetCellValue(sheet, c.getCellName("A", row), day.Weekday+" "+day.Date)
		if len(day.Entries) == 0 {
			f.SetCellValue(sheet, c.getCellName("B", row), "-")
			f.SetCellValue(sheet, c.getCellName("C", row), "-")
			f.SetCellValue(sheet, c.getCellName("D", row), "-")
			f.SetCellValue(sheet, c.getCellName("E", row), "-")
			row++
		} else {
			for _, entry := range day.Entries {
				f.SetCellValue(sheet, c.getCellName("B", row), entry.EntryType)
				f.SetCellValue(sheet, c.getCellName("C", row), entry.StartTime)
				f.SetCellValue(sheet, c.getCellName("D", row), entry.EndTime)
				f.SetCellValue(sheet, c.getCellName("E", row), entry.Duration)
				f.SetCellValue(sheet, c.getCellName("F", row), entry.EntryActivity)
				f.SetCellValue(sheet, c.getCellName("G", row), entry.Description)
				row++
			}
		}
		if len(day.Entries) > 1 {
			f.SetCellValue(sheet, c.getCellName("E", row), day.WorkDuration)
			row++
		}
	}
	f.SetCellStyle(sheet, "A14", c.getCellName("E", row-1), styleTableBody)
	f.SetCellStyle(sheet, "F14", c.getCellName("G", row-1), styleTableBody)

	return f
}

func (c *OverviewController) getCellName(col string, row int) string {
	return col + strconv.Itoa(row)
}

func (c *OverviewController) writeFile(r *echo.Response, fileName string, wt io.WriterTo) error {
	// Write header
	r.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	r.Header().Set("Content-Type", "application/octet-stream")
	r.Header().Set("Content-Transfer-Encoding", "binary")
	r.Header().Set("Expires", "0")

	// Write body
	_, wErr := wt.WriteTo(r.Writer)
	if wErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not write response.", wErr)
		log.Debug(err.StackTrace())
		return err
	}

	return nil
}
