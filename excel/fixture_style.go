package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// Light header fill and darker grid lines for fixture workbooks (not black).
const (
	fixtureHeaderFillColor = "#EDEDED"
	fixtureGridBorderColor = "#9E9E9E"
	fixtureBorderLineStyle = 1
)

func applyFixtureWorkbookStyle(f *excelize.File, sheet string, prep SchemaPrep) error {
	headerStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{fixtureHeaderFillColor},
			Pattern: 1,
		},
		Border: fixtureGrayBorders(),
	})
	if err != nil {
		return fmt.Errorf("applyFixtureWorkbookStyle: header style: %w", err)
	}
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: fixtureGrayBorders(),
	})
	if err != nil {
		return fmt.Errorf("applyFixtureWorkbookStyle: data style: %w", err)
	}

	hTL, hBR, err := cellRangeAxis(prep.HeaderRange())
	if err != nil {
		return err
	}
	if err := f.SetCellStyle(sheet, hTL, hBR, headerStyle); err != nil {
		return fmt.Errorf("applyFixtureWorkbookStyle: header range: %w", err)
	}

	dTL, dBR, err := cellRangeAxis(prep.DataRange())
	if err != nil {
		return err
	}
	if err := f.SetCellStyle(sheet, dTL, dBR, dataStyle); err != nil {
		return fmt.Errorf("applyFixtureWorkbookStyle: data range: %w", err)
	}
	return nil
}

func fixtureGrayBorders() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: fixtureGridBorderColor, Style: fixtureBorderLineStyle},
		{Type: "right", Color: fixtureGridBorderColor, Style: fixtureBorderLineStyle},
		{Type: "top", Color: fixtureGridBorderColor, Style: fixtureBorderLineStyle},
		{Type: "bottom", Color: fixtureGridBorderColor, Style: fixtureBorderLineStyle},
	}
}

func cellRangeAxis(r CellRange) (topLeft, bottomRight string, err error) {
	topLeft, err = excelize.JoinCellName(r.StartColumnName(), r.StartRowNumber())
	if err != nil {
		return "", "", fmt.Errorf("cellRangeAxis: %w", err)
	}
	bottomRight, err = excelize.JoinCellName(r.EndColumnName(), r.EndRowNumber())
	if err != nil {
		return "", "", fmt.Errorf("cellRangeAxis: %w", err)
	}
	return topLeft, bottomRight, nil
}
