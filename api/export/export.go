package export

import (
	"encoding/csv"
	"io"
	"strings"
	"time"

	"kellnhofer.com/work-log/pkg/model"
)

type csvWriterToAdapter struct {
	data [][]string
}

func (cta *csvWriterToAdapter) WriteTo(w io.Writer) (n int64, err error) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	totalBytes := int64(0)
	for _, record := range cta.data {
		if err := writer.Write(record); err != nil {
			return totalBytes, err
		}
		for _, field := range record {
			totalBytes += int64(len(field)) + 1 // +1 for delimiter/newline
		}
	}

	return totalBytes, writer.Error()
}

// EntriesExporter exports entries to a CSV file.
type EntriesExporter struct {
}

// NewEntriesExporter creates a new entries exporter.
func NewEntriesExporter() *EntriesExporter {
	return &EntriesExporter{}
}

// ExportEntries creates the CSV file for the supplied data and returns it as an io.WriterTo that
// can be used to write the file to a writer.
func (e *EntriesExporter) ExportEntries(entries []*model.Entry, entryTypes []*model.EntryType,
	entryActivities []*model.EntryActivity) io.WriterTo {
	// Create maps for lookups
	entryTypesMap := make(map[int]string)
	for _, entryType := range entryTypes {
		entryTypesMap[entryType.Id] = entryType.Description
	}
	entryActivitiesMap := make(map[int]string)
	for _, entryActivity := range entryActivities {
		entryActivitiesMap[entryActivity.Id] = entryActivity.Description
	}

	var data [][]string

	// Create and append header
	header := []string{
		"Start Time",
		"End Time",
		"Type",
		"Activity",
		"Description",
		"Labels",
	}
	data = append(data, header)

	// Create and append records
	for _, entry := range entries {
		record := []string{
			entry.StartTime.Format(time.RFC3339),
			entry.EndTime.Format(time.RFC3339),
			e.getEntryTypeDescription(entryTypesMap, entry.TypeId),
			e.getEntryActivityDescription(entryActivitiesMap, entry.ActivityId),
			entry.Description,
			strings.Join(entry.Labels, " "),
		}
		data = append(data, record)
	}

	return &csvWriterToAdapter{
		data: data,
	}
}

func (e *EntriesExporter) getEntryTypeDescription(entryTypes map[int]string, id int) string {
	et, ok := entryTypes[id]
	if ok {
		return et
	}
	return ""
}

func (e *EntriesExporter) getEntryActivityDescription(entryActivities map[int]string, id int) string {
	et, ok := entryActivities[id]
	if ok {
		return et
	}
	return ""
}
