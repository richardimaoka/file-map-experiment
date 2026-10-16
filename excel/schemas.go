package excel

type ColumnSchema struct {
	HeaderName   string `json:"header_name"`
	ExcelColumn  string `json:"excel_column"`
	IsRowGrouped bool   `json:"is_row_grouped"`
}

type TableSchema struct {
	SchemaName      string         `json:"schema_name"`
	DataStartRow    int            `json:"data_start_row"`
	DataStartColumn string         `json:"data_start_column"`
	DataEndColumn   string         `json:"data_end_column"`
	Columns         []ColumnSchema `json:"columns"`
}
