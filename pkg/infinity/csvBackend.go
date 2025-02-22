package infinity

import (
	"fmt"
	"strings"

	"github.com/appkube/cloud-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/yesoreyeram/grafana-framer/csvFramer"
	"github.com/yesoreyeram/grafana-framer/gframer"
)

func GetCSVBackendResponse(responseString string, query models.Query) (*data.Frame, error) {
	frame := GetDummyFrame(query)
	columns := []gframer.ColumnSelector{}
	for _, c := range query.Columns {
		columns = append(columns, gframer.ColumnSelector{
			Selector:   c.Selector,
			Alias:      c.Text,
			Type:       c.Type,
			TimeFormat: c.TimeStampFormat,
		})
	}
	csvOptions := csvFramer.CSVFramerOptions{
		FrameName:          query.RefID,
		Columns:            columns,
		Comment:            query.CSVOptions.Comment,
		Delimiter:          query.CSVOptions.Delimiter,
		SkipLinesWithError: query.CSVOptions.SkipLinesWithError,
		RelaxColumnCount:   query.CSVOptions.RelaxColumnCount,
	}
	if query.CSVOptions.Columns != "" && query.CSVOptions.Columns != "-" && query.CSVOptions.Columns != "none" {
		responseString = query.CSVOptions.Columns + "\n" + responseString
	}
	if query.CSVOptions.Columns == "-" || query.CSVOptions.Columns == "none" {
		csvOptions.NoHeaders = true
	}
	if query.Type == models.QueryTypeTSV {
		csvOptions.Delimiter = "\t"
	}
	newFrame, err := csvFramer.CsvStringToFrame(responseString, csvOptions)
	frame.Meta = &data.FrameMeta{
		Custom: &CustomMeta{
			Query: query,
		},
	}
	if err != nil {
		backend.Logger.Error("error getting response for query", "error", err.Error())
		frame.Meta.Custom = &CustomMeta{
			Query: query,
			Error: err.Error(),
		}
		return frame, err
	}
	if newFrame != nil {
		frame.Fields = append(frame.Fields, newFrame.Fields...)
	}
	frame, err = GetFrameWithComputedColumns(frame, query.ComputedColumns)
	if err != nil {
		backend.Logger.Error("error getting computed column", "error", err.Error())
		frame.Meta.Custom = &CustomMeta{Query: query, Error: err.Error()}
		return frame, err
	}
	frame, err = ApplyFilter(frame, query.FilterExpression)
	if err != nil {
		backend.Logger.Error("error applying filter", "error", err.Error())
		frame.Meta.Custom = &CustomMeta{Query: query, Error: err.Error()}
		return frame, fmt.Errorf("error applying filter. %w", err)
	}
	if strings.TrimSpace(query.SummarizeExpression) != "" {
		return GetSummaryFrame(frame, query.SummarizeExpression, query.SummarizeBy)
	}
	if query.Format == "timeseries" && frame.TimeSeriesSchema().Type == data.TimeSeriesTypeLong {
		if wFrame, err := data.LongToWide(frame, &data.FillMissing{Mode: data.FillModeNull}); err == nil {
			return wFrame, err
		}
	}
	return frame, err
}
