package excel

type SchemaPrep struct {
	SchemaName     string `json:"schema_name"`
	HeaderStartRow int    `json:"header_start_row"`
	HeaderEndRow   int    `json:"header_end_row"`
	DataStartRow   int    `json:"data_start_row"`
	StartColumn    string `json:"start_column"`
	EndColumn      string `json:"end_column"`
}

type Config struct {
	ExcelFileName string       `json:"excel_file_name"`
	SchemaPreps   []SchemaPrep `json:"schema_preps"`
}
