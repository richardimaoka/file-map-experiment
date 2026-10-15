package excel

type ExcelFolderMap struct {
	excelFileName string
	preConfigName string
}

func (e *ExcelFolderMap) ExcelFilePath() string {
	return ""
}

func (e *ExcelFolderMap) PrepConfigPath() string {
	return ""
}

func (e *ExcelFolderMap) SchemaJsonPath() string {
	return ""
}

func (e *ExcelFolderMap) DataJsonPath() string {
	return ""
}
