package gotil

import (
	"encoding/csv"
	"encoding/json"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// parsing json file to struct / map / slice
func ParseJSONFile(filename string, result any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	var _result JSON
	err = json.NewDecoder(f).Decode(&_result)
	if err != nil {
		return err
	}
	return ParseStruct(&result, _result, "json")
}

// parsing csv file to struct / map / slice
func ParseCSVFile(filename string, result any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	reader := csv.NewReader(f)
	data, err := reader.ReadAll()
	if err != nil {
		return err
	}
	headers := data[0]
	columns := data[1:]

	var _result []JSON
	for _, col := range columns {
		_data := make(JSON)
		for k, val := range col {
			_data[headers[k]] = val
		}
		_result = append(_result, _data)
	}

	return ParseStruct(&result, _result, "json")
}

// parsing excel file to struct / map / slice
func ParseExcelFile(filename string, result any) error {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}

	sheetName := f.GetSheetName(1)
	rows := f.GetRows(sheetName)

	var _result []JSON
	headers := rows[0]

	for _, row := range rows[1:] {
		record := make(JSON)
		for i, cell := range row {
			if i < len(headers) {
				record[headers[i]] = cell
			}
		}
		_result = append(_result, record)
	}

	return ParseStruct(&result, _result, "json")
}
