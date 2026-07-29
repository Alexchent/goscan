package duplicate

// DuplicateFile 重复文件信息结构
type DuplicateFile struct {
	FullFileName string
	MD5          string
	Size         int64
}

// FileInfo 文件信息
type FileInfo struct {
	Name          string
	Path          string
	Size          int64
	SizeFormatted string
}

// DuplicateResult 重复文件结果
type DuplicateResult struct {
	MD5           string
	Files         []FileInfo
	Size          int64
	SizeFormatted string
	Count         int
	Saved         int64
}

// PageData 页面数据
type PageData struct {
	TotalFiles      int64
	DuplicateGroups int
	DuplicateCount  int
	SavedSpace      string
	Results         []DuplicateResult
}

// 全局变量存储扫描结果
var scanResults []DuplicateResult
var totalFilesScanned int64
