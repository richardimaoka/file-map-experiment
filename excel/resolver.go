package excel

type PathResolver struct {
	dirPath       string
	excelFileName string
}

func baseDir() string {
	appBaseDir := os.GetEnv("APP_BASE_DIR")
	if appBaseDir == "" {
		panic("APP_BASE_DIR is not set")
	}
	return appBaseDir
}

func NewPathResolver(dirName, excelFileName string) *PathResolver {
	return &PathResolver{
		dirPath:       dirName,
		excelFileName: excelFileName,
	}
}

// func (p *PathResolver) ExcelPath() string {
// 	return
// }
