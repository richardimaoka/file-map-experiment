package client1

type ForwardReverse int

type InternalFolderMap struct {
	excelFileName string
	preConfigName string
	direction     ForwardReverse
}

type ExcelFolderMap struct {
	excelFileName string
	preConfigName string
	Forward       InternalFolderMap
	Reverse       InternalFolderMap
}

func (e *ExcelFolderMap) DefaultDataJsonPath() string {
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
