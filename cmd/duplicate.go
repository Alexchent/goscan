package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Alexchent/goscan/help"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	consumer     = 100
	webPort      = 8080
	templateHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>重复文件管理器</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; padding: 20px; }
        .container { max-width: 1400px; margin: 0 auto; }
        h1 { color: #2c3e50; margin-bottom: 20px; text-align: center; }
        .stats { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; border-radius: 10px; margin-bottom: 20px; display: flex; justify-content: space-around; flex-wrap: wrap; }
        .stat-item { text-align: center; padding: 10px; }
        .stat-value { font-size: 28px; font-weight: bold; }
        .stat-label { font-size: 14px; opacity: 0.9; }
        .toolbar { background: white; padding: 15px; border-radius: 10px; margin-bottom: 20px; display: flex; gap: 10px; flex-wrap: wrap; align-items: center; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .btn { padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; font-size: 14px; transition: all 0.3s; }
        .btn-primary { background: #3498db; color: white; }
        .btn-primary:hover { background: #2980b9; }
        .btn-danger { background: #e74c3c; color: white; }
        .btn-danger:hover { background: #c0392b; }
        .btn-success { background: #27ae60; color: white; }
        .btn-success:hover { background: #229954; }
        .btn-warning { background: #f39c12; color: white; }
        .btn-warning:hover { background: #d68910; }
        .group { background: white; border-radius: 10px; margin-bottom: 15px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .group-header { background: #34495e; color: white; padding: 15px; display: flex; justify-content: space-between; align-items: center; }
        .group-title { font-weight: bold; }
        .group-info { font-size: 14px; opacity: 0.9; }
        .file-list { padding: 0; }
        .file-item { display: flex; align-items: center; padding: 12px 15px; border-bottom: 1px solid #ecf0f1; transition: background 0.2s; }
        .file-item:hover { background: #f8f9fa; }
        .file-item:last-child { border-bottom: none; }
        .file-item.original { background: #e8f8f5; }
        .file-item.original:hover { background: #d5f4e6; }
        .file-checkbox { width: 20px; height: 20px; margin-right: 15px; cursor: pointer; }
        .file-icon { font-size: 24px; margin-right: 15px; }
        .file-info { flex: 1; }
        .file-name { font-weight: 500; color: #2c3e50; word-break: break-all; }
        .file-path { font-size: 12px; color: #7f8c8d; margin-top: 4px; word-break: break-all; }
        .file-size { color: #95a5a6; font-size: 14px; margin-left: 15px; white-space: nowrap; }
        .original-badge { background: #27ae60; color: white; padding: 2px 8px; border-radius: 3px; font-size: 12px; margin-left: 10px; }
        .duplicate-badge { background: #e74c3c; color: white; padding: 2px 8px; border-radius: 3px; font-size: 12px; margin-left: 10px; }
        .empty-state { text-align: center; padding: 60px 20px; color: #7f8c8d; }
        .empty-icon { font-size: 64px; margin-bottom: 20px; }
        .loading { text-align: center; padding: 40px; }
        .spinner { border: 4px solid #f3f3f3; border-top: 4px solid #3498db; border-radius: 50%; width: 40px; height: 40px; animation: spin 1s linear infinite; margin: 0 auto 20px; }
        @keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
        .toast { position: fixed; bottom: 20px; right: 20px; padding: 15px 25px; border-radius: 5px; color: white; font-weight: 500; opacity: 0; transform: translateY(20px); transition: all 0.3s; z-index: 1000; }
        .toast.show { opacity: 1; transform: translateY(0); }
        .toast.success { background: #27ae60; }
        .toast.error { background: #e74c3c; }
        .select-all-group { margin-right: 10px; }
        #searchBox { padding: 10px 15px; border: 1px solid #ddd; border-radius: 5px; font-size: 14px; width: 250px; }
        #searchBox:focus { outline: none; border-color: #3498db; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔍 重复文件管理器</h1>
        
        <div class="stats">
            <div class="stat-item">
                <div class="stat-value" id="totalFiles">{{.TotalFiles}}</div>
                <div class="stat-label">总文件数</div>
            </div>
            <div class="stat-item">
                <div class="stat-value" id="duplicateGroups">{{.DuplicateGroups}}</div>
                <div class="stat-label">重复文件组</div>
            </div>
            <div class="stat-item">
                <div class="stat-value" id="duplicateCount">{{.DuplicateCount}}</div>
                <div class="stat-label">重复文件数</div>
            </div>
            <div class="stat-item">
                <div class="stat-value" id="savedSpace">{{.SavedSpace}}</div>
                <div class="stat-label">可节省空间</div>
            </div>
        </div>

        <div class="toolbar">
            <input type="text" id="searchBox" placeholder="搜索文件名..." onkeyup="filterFiles()">
            <button class="btn btn-warning" onclick="selectAllDuplicates()">全选所有副本</button>
            <button class="btn btn-primary" onclick="deselectAll()">取消全选</button>
            <button class="btn btn-danger" onclick="deleteSelected()">删除选中文件</button>
            <button class="btn btn-success" onclick="exportCSV()">导出CSV</button>
            <span id="selectedCount" style="margin-left: auto; color: #7f8c8d;"></span>
        </div>

        {{if eq (len .Results) 0}}
        <div class="empty-state">
            <div class="empty-icon">✨</div>
            <h2>没有发现重复文件</h2>
            <p>您的目录很干净，没有重复的文件</p>
        </div>
        {{else}}
        {{range $index, $result := .Results}}
        <div class="group" data-md5="{{$result.MD5}}">
            <div class="group-header">
                <div>
                    <span class="group-title">📦 重复组 #{{add $index 1}}</span>
                    <span class="group-info" style="margin-left: 15px;">MD5: {{$result.MD5}}</span>
                </div>
                <div>
                    <span class="group-info">{{$result.Count}} 个文件 · {{$result.SizeFormatted}}</span>
                    <button class="btn btn-danger select-all-group" onclick="mergeGroupDuplicates({{$index}}, {{$result.Count}})" style="margin-left: 10px; padding: 5px 10px;">合并副本</button>
                    <button class="btn btn-warning select-all-group" onclick="selectGroupDuplicates({{$index}})" style="margin-left: 5px; padding: 5px 10px;">全选副本</button>
                </div>
            </div>
            <div class="file-list">
                {{range $fileIndex, $file := $result.Files}}
                <div class="file-item {{if eq $fileIndex 0}}original{{end}}" data-group="{{$index}}" data-file="{{$fileIndex}}">
                    <input type="checkbox" class="file-checkbox" data-filepath="{{$file.Path}}" data-group="{{$index}}" data-index="{{$fileIndex}}" onchange="updateSelectedCount()">
                    <div class="file-icon">{{if eq $fileIndex 0}}📁{{else}}📄{{end}}</div>
                    <div class="file-info">
                        <div class="file-name">{{$file.Name}} {{if eq $fileIndex 0}}<span class="original-badge">原文件</span>{{else}}<span class="duplicate-badge">副本</span>{{end}}</div>
                        <div class="file-path">{{$file.Path}}</div>
                    </div>
                    <div class="file-size">{{$file.SizeFormatted}}</div>
                </div>
                {{end}}
            </div>
        </div>
        {{end}}
        {{end}}
    </div>

    <div id="toast" class="toast"></div>

    <script>
        function showToast(message, type) {
            const toast = document.getElementById('toast');
            toast.textContent = message;
            toast.className = 'toast ' + type + ' show';
            setTimeout(() => toast.classList.remove('show'), 3000);
        }

        function updateSelectedCount() {
            const checked = document.querySelectorAll('.file-checkbox:checked').length;
            document.getElementById('selectedCount').textContent = checked > 0 ? '已选择 ' + checked + ' 个文件' : '';
        }

        function selectAllDuplicates() {
            document.querySelectorAll('.file-checkbox').forEach((cb, index) => {
                const groupIndex = parseInt(cb.dataset.group);
                const fileIndex = parseInt(cb.dataset.index);
                if (fileIndex > 0) cb.checked = true;
            });
            updateSelectedCount();
            showToast('已选择所有副本文件', 'success');
        }

        function selectGroupDuplicates(groupIndex) {
            document.querySelectorAll('.file-checkbox[data-group="' + groupIndex + '"]').forEach((cb, index) => {
                const fileIndex = parseInt(cb.dataset.index);
                if (fileIndex > 0) cb.checked = true;
            });
            updateSelectedCount();
        }

        function mergeGroupDuplicates(groupIndex, totalCount) {
            const duplicateCount = totalCount - 1;
            if (duplicateCount <= 0) {
                showToast('该组没有副本文件', 'error');
                return;
            }
            
            if (!confirm('确定要合并此重复组吗？\n\n将删除 ' + duplicateCount + ' 个副本文件，仅保留原文件。\n\n此操作不可恢复！')) {
                return;
            }

            const files = [];
            document.querySelectorAll('.file-checkbox[data-group="' + groupIndex + '"]').forEach(cb => {
                const fileIndex = parseInt(cb.dataset.index);
                if (fileIndex > 0) {
                    files.push(cb.dataset.filepath);
                }
            });

            if (files.length === 0) {
                showToast('没有可删除的副本文件', 'error');
                return;
            }

            fetch('/delete', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({files: files})
            })
            .then(r => r.json())
            .then(data => {
                if (data.success) {
                    showToast('成功合并，删除了 ' + data.deleted + ' 个副本文件', 'success');
                    if (data.failed > 0) {
                        setTimeout(() => showToast(data.failed + ' 个文件删除失败', 'error'), 3500);
                    }
                    setTimeout(() => location.reload(), 1500);
                } else {
                    showToast('合并失败: ' + data.error, 'error');
                }
            })
            .catch(err => showToast('请求失败: ' + err, 'error'));
        }

        function deselectAll() {
            document.querySelectorAll('.file-checkbox').forEach(cb => cb.checked = false);
            updateSelectedCount();
            showToast('已取消所有选择', 'success');
        }

        function deleteSelected() {
            const files = [];
            const checkedBoxes = document.querySelectorAll('.file-checkbox:checked');
            checkedBoxes.forEach(cb => {
                files.push(cb.dataset.filepath);
            });
            
            if (files.length === 0) {
                showToast('请先选择要删除的文件', 'error');
                return;
            }

            // 检查是否有重复组的所有文件都被选中
            const groups = {};
            checkedBoxes.forEach(cb => {
                const group = cb.dataset.group;
                if (!groups[group]) groups[group] = { total: 0, checked: 0 };
            });
            
            // 统计每个组的总数和选中数
            document.querySelectorAll('.file-checkbox').forEach(cb => {
                const group = cb.dataset.group;
                if (groups[group]) {
                    groups[group].total++;
                    if (cb.checked) groups[group].checked++;
                }
            });

            // 检查是否有组被全选
            for (const group in groups) {
                if (groups[group].checked >= groups[group].total) {
                    showToast('不能删除重复组中的所有文件，请至少保留一个文件！', 'error');
                    return;
                }
            }

            if (!confirm('确定要删除选中的 ' + files.length + ' 个文件吗？\n\n此操作不可恢复！')) {
                return;
            }

            fetch('/delete', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({files: files})
            })
            .then(r => r.json())
            .then(data => {
                if (data.success) {
                    showToast('成功删除 ' + data.deleted + ' 个文件', 'success');
                    if (data.failed > 0) {
                        setTimeout(() => showToast(data.failed + ' 个文件删除失败', 'error'), 3500);
                    }
                    setTimeout(() => location.reload(), 1500);
                } else {
                    showToast('删除失败: ' + data.error, 'error');
                }
            })
            .catch(err => showToast('请求失败: ' + err, 'error'));
        }

        function exportCSV() {
            window.location.href = '/export';
        }

        function filterFiles() {
            const keyword = document.getElementById('searchBox').value.toLowerCase();
            document.querySelectorAll('.group').forEach(group => {
                let hasMatch = false;
                group.querySelectorAll('.file-item').forEach(item => {
                    const name = item.querySelector('.file-name').textContent.toLowerCase();
                    if (name.includes(keyword)) {
                        item.style.display = 'flex';
                        hasMatch = true;
                    } else {
                        item.style.display = keyword ? 'none' : 'flex';
                    }
                });
                group.style.display = hasMatch || !keyword ? 'block' : 'none';
            });
        }

        updateSelectedCount();
    </script>
</body>
</html>`
)

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

// duplicateCmd represents the duplicate command
var duplicateCmd = &cobra.Command{
	Use:   "duplicate",
	Short: "查找重复文件",
	Long:  `查找重复文件，支持并发扫描，结果输出到控制台、CSV文件和Web页面`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		// 检查目录参数
		if path == "" {
			fmt.Printf("请输入要扫描的目录: ")
			_, _ = fmt.Scan(&path)
		}

		_, err := os.Stat(path)
		if err != nil {
			fmt.Println(color.RedString("错误: %v", err))
			return
		}

		fmt.Println(color.CyanString("========================================"))
		fmt.Println(color.CyanString("开始扫描目录: %s", path))
		fmt.Println(color.CyanString("========================================"))

		// 使用 channel 进行并发处理
		fileChan := make(chan *DuplicateFile, consumer)
		var scanWg sync.WaitGroup
		scanWg.Add(1)

		// 文件计数器
		var fileCount int64
		var countMutex sync.Mutex

		// 启动扫描 goroutine
		go func() {
			walkFiles(path, &scanWg, fileChan, &fileCount, &countMutex)
		}()

		// 收集结果
		fileList := make(map[string][]*DuplicateFile)
		var collectWg sync.WaitGroup
		collectWg.Add(1)

		go func() {
			defer collectWg.Done()
			for f := range fileChan {
				fileList[f.MD5] = append(fileList[f.MD5], f)
			}
		}()

		// 等待扫描完成
		scanWg.Wait()
		close(fileChan)
		collectWg.Wait()

		// 处理结果
		results := processResults(fileList)
		totalFilesScanned = fileCount
		scanResults = results

		// 输出到控制台
		printToConsole(results, fileCount, start)

		// 保存到 CSV 文件
		saveToCSV(results, path)

		// 启动 Web 服务器
		startWebServer()
	},
}

// walkFiles 递归遍历目录
func walkFiles(dir string, wg *sync.WaitGroup, ch chan *DuplicateFile, fileCount *int64, countMutex *sync.Mutex) {
	defer wg.Done()

	readDir, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var nextWg sync.WaitGroup
	for _, file := range readDir {
		fileName := file.Name()
		// 过滤隐藏文件
		if fileName[0] == '.' {
			continue
		}

		fullName := filepath.Join(dir, fileName)

		if file.IsDir() {
			nextWg.Add(1)
			go walkFiles(fullName, &nextWg, ch, fileCount, countMutex)
		} else {
			// 获取文件信息
			info, err := file.Info()
			if err != nil {
				continue
			}

			// 计数
			countMutex.Lock()
			*fileCount++
			countMutex.Unlock()

			// 发送到 channel
			ch <- &DuplicateFile{
				FullFileName: fullName,
				MD5:          help.GetFileMd5(fullName),
				Size:         info.Size(),
			}
		}
	}

	nextWg.Wait()
}

// processResults 处理扫描结果，筛选重复文件
func processResults(fileList map[string][]*DuplicateFile) []DuplicateResult {
	var results []DuplicateResult

	for md5, files := range fileList {
		if len(files) > 1 {
			var fileInfos []FileInfo
			var size int64
			for _, f := range files {
				fileInfos = append(fileInfos, FileInfo{
					Name:          filepath.Base(f.FullFileName),
					Path:          f.FullFileName,
					Size:          f.Size,
					SizeFormatted: help.FormatFileSize1000(f.Size),
				})
				size = f.Size
			}

			results = append(results, DuplicateResult{
				MD5:           md5,
				Files:         fileInfos,
				Size:          size,
				SizeFormatted: help.FormatFileSize1000(size),
				Count:         len(files),
				Saved:         size * int64(len(files)-1),
			})
		}
	}

	return results
}

// printToConsole 输出到控制台
func printToConsole(results []DuplicateResult, totalFiles int64, start time.Time) {
	fmt.Println(color.CyanString("\n========================================"))
	fmt.Println(color.CyanString("扫描结果"))
	fmt.Println(color.CyanString("========================================"))

	fmt.Printf("总文件数: %s\n", color.YellowString("%d", totalFiles))
	fmt.Printf("重复文件组: %s\n", color.YellowString("%d", len(results)))

	if len(results) == 0 {
		fmt.Println(color.GreenString("\n✓ 没有发现重复文件"))
		return
	}

	// 计算总可节省空间
	var totalSaved int64
	var duplicateCount int
	for _, r := range results {
		totalSaved += r.Saved
		duplicateCount += r.Count - 1
	}
	fmt.Printf("重复文件数: %s\n", color.YellowString("%d", duplicateCount))
	fmt.Printf("可节省空间: %s\n", color.RedString(help.FormatFileSize1000(totalSaved)))

	fmt.Println(color.CyanString("\n----------------------------------------"))
	fmt.Println(color.CyanString("重复文件详情:"))
	fmt.Println(color.CyanString("----------------------------------------"))

	for i, r := range results {
		fmt.Printf("\n%s %s (共 %d 个副本，大小: %s)\n",
			color.BlueString("[%d]", i+1),
			color.RedString("MD5: %s...", r.MD5[:16]),
			r.Count,
			color.YellowString(r.SizeFormatted))

		for j, file := range r.Files {
			if j == 0 {
				fmt.Printf("  %s %s\n", color.GreenString("★ 原文件:"), file.Path)
			} else {
				fmt.Printf("  %s %s\n", color.RedString("  副本:"), file.Path)
			}
		}
	}
}

// saveToCSV 保存结果到 CSV 文件
func saveToCSV(results []DuplicateResult, scanPath string) {
	if len(results) == 0 {
		return
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("duplicate_files_%s.csv", timestamp)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(color.RedString("创建CSV文件失败: %v", err))
		return
	}
	defer file.Close()

	// 写入 UTF-8 BOM，确保 Excel 正确识别编码
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	header := []string{"序号", "MD5", "文件大小", "副本数量", "可节省空间", "文件路径"}
	if err := writer.Write(header); err != nil {
		fmt.Println(color.RedString("写入CSV表头失败: %v", err))
		return
	}

	// 写入数据
	for i, r := range results {
		for j, fileInfo := range r.Files {
			row := []string{
				fmt.Sprintf("%d", i+1),
				r.MD5,
				r.SizeFormatted,
				fmt.Sprintf("%d", r.Count),
				help.FormatFileSize1000(r.Saved),
				fileInfo.Path,
			}
			if j == 0 {
				row[0] = fmt.Sprintf("%d (原文件)", i+1)
			} else {
				row[0] = fmt.Sprintf("%d (副本%d)", i+1, j)
			}
			if err := writer.Write(row); err != nil {
				fmt.Println(color.RedString("写入CSV数据失败: %v", err))
				return
			}
		}
	}

	fmt.Printf("\n%s %s\n", color.GreenString("✓ 结果已保存到:"), color.YellowString(filename))
}

// startWebServer 启动 Web 服务器
func startWebServer() {
	// 计算统计数据
	var totalSaved int64
	var duplicateCount int
	for _, r := range scanResults {
		totalSaved += r.Saved
		duplicateCount += r.Count - 1
	}

	// 主页面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index").Funcs(template.FuncMap{
			"add": func(a, b int) int { return a + b },
			"len": func(s []DuplicateResult) int { return len(s) },
		}).Parse(templateHTML)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := PageData{
			TotalFiles:      totalFilesScanned,
			DuplicateGroups: len(scanResults),
			DuplicateCount:  duplicateCount,
			SavedSpace:      help.FormatFileSize1000(totalSaved),
			Results:         scanResults,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
	})

	// 删除文件 API
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Files []string `json:"files"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var deleted, failed int
		for _, file := range req.Files {
			if err := os.Remove(file); err != nil {
				failed++
				fmt.Printf("删除失败: %s - %v\n", file, err)
			} else {
				deleted++
				fmt.Printf("已删除: %s\n", file)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"deleted": deleted,
			"failed":  failed,
		})
	})

	// 导出 CSV
	http.HandleFunc("/export", func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("duplicate_files_%s.csv", timestamp)

		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

		// UTF-8 BOM
		w.Write([]byte("\xEF\xBB\xBF"))

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// 写入表头
		writer.Write([]string{"序号", "MD5", "文件大小", "副本数量", "可节省空间", "文件路径"})

		// 写入数据
		for i, r := range scanResults {
			for j, fileInfo := range r.Files {
				row := []string{
					fmt.Sprintf("%d", i+1),
					r.MD5,
					r.SizeFormatted,
					fmt.Sprintf("%d", r.Count),
					help.FormatFileSize1000(r.Saved),
					fileInfo.Path,
				}
				if j == 0 {
					row[0] = fmt.Sprintf("%d (原文件)", i+1)
				} else {
					row[0] = fmt.Sprintf("%d (副本%d)", i+1, j)
				}
				writer.Write(row)
			}
		}
	})

	fmt.Println(color.CyanString("\n========================================"))
	fmt.Println(color.GreenString("✓ Web 服务器启动成功"))
	fmt.Println(color.YellowString("请在浏览器中打开: http://localhost:%d", webPort))
	fmt.Println(color.CyanString("========================================"))
	fmt.Println(color.HiBlackString("按 Ctrl+C 停止服务器"))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", webPort), nil); err != nil {
		fmt.Println(color.RedString("Web服务器启动失败: %v", err))
	}
}

func init() {
	duplicateCmd.Flags().StringVarP(&path, "dir", "d", "", "扫描的目录")
	rootCmd.AddCommand(duplicateCmd)
}
