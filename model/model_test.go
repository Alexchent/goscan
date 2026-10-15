package model

import (
	"testing"
	"time"
)

func newTestScanFile(t *testing.T) *ScanFile {
	return NewScanFile(t.TempDir(), "")
}

func TestNormalizeName(t *testing.T) {
	cases := map[string]string{
		"Ab-C_d/File.TXT": "abcd/file.txt",
		"hello-world":     "helloworld",
		"文件名-导入版_v2":     "文件名导入版v2",
	}
	for in, want := range cases {
		if got := normalizeName(in); got != want {
			t.Errorf("normalizeName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSearchNormalization(t *testing.T) {
	sf := newTestScanFile(t)
	paths := []string{
		"/data/File-Report_v1/TaskA.txt",
		"/data/File_Report-v2/TaskB.txt",
		"/other/Unrelated/doc.md",
	}
	for _, p := range paths {
		sf.AddRecord(p, 1)
	}
	if got := sf.Search("file-report"); len(got) != 2 {
		t.Errorf("Search('file-report') got %d results, want 2: %v", len(got), got)
	}
	if got := sf.Search("filereport"); len(got) != 2 {
		t.Errorf("Search('filereport') got %d results, want 2: %v", len(got), got)
	}
	if got := sf.Search("-"); len(got) != 0 {
		t.Errorf("Search('-') got %d results, want 0: %v", len(got), got)
	}
}

func TestSearchIsCaseInsensitiveSubstring(t *testing.T) {
	sf := newTestScanFile(t)
	sf.AddRecord("/movies/Interstellar.2014.mkv", 1)
	if got := sf.Search("interstellar"); len(got) != 1 {
		t.Errorf("Search('interstellar') got %d, want 1", len(got))
	}
	if got := sf.Search("2014"); len(got) != 1 {
		t.Errorf("Search('2014') got %d, want 1", len(got))
	}
}

func TestAddRecordsDedupAndReturnNew(t *testing.T) {
	sf := newTestScanFile(t)
	first := sf.AddRecords([]*ScanRecord{
		{FullFileName: "/a/1.txt", Size: 1},
		{FullFileName: "/a/2.txt", Size: 2},
	})
	if len(first) != 2 {
		t.Fatalf("first batch added = %d, want 2", len(first))
	}
	second := sf.AddRecords([]*ScanRecord{
		{FullFileName: "/a/1.txt", Size: 1}, // 已存在
		{FullFileName: "/a/3.txt", Size: 3}, // 新增
		{FullFileName: "/a/4.txt", Size: 4}, // 新增
	})
	if len(second) != 2 {
		t.Fatalf("second batch added = %d, want 2 (got %v)", len(second), second)
	}
	var names []string
	for _, r := range second {
		names = append(names, r.FullFileName)
	}
	for _, want := range []string{"/a/3.txt", "/a/4.txt"} {
		found := false
		for _, n := range names {
			if n == want {
				found = true
			}
		}
		if !found {
			t.Errorf("batch should contain newly added %q, got %v", want, names)
		}
	}
	if got := len(sf.AllRecords()); got != 4 {
		t.Errorf("AllRecords count = %d, want 4", got)
	}
}

// TestBackfillNorm 验证存量数据（norm_file_name 为空）能被回填。
func TestBackfillNorm(t *testing.T) {
	sf := newTestScanFile(t)

	// 模拟旧版本写入的行：norm_file_name 为空
	raw := ScanRecord{
		FullFileName: "/x/File-Test_v1.txt",
		NormFileName: "",
		Size:         5,
		Created:      time.Now(),
	}
	if err := sf.db.Create(&raw).Error; err != nil {
		t.Fatalf("seed raw row: %v", err)
	}

	// 触发回填
	sf.backfillNorm()

	if got := sf.Search("filetestv1"); len(got) != 1 {
		t.Errorf("after backfill Search('filetestv1') got %d, want 1 (got %v)", len(got), got)
	}
}
