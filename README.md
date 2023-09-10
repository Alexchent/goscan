# go-scan

利用corba编写的本地文件扫描程序

功能包括：
- 扫描制定目录的文件包括子目录，如扫描下载目录的所有文件：`./scan start ~/Download`
- 将所有扫描结果导出 `./scan export`
- 将导出的文件导入到程序中
- 从扫描记录中搜索文件。模糊查询 `./scan find {filenmae}`

## 编译
1. 编译成本地运行的程序
```bash
make local
```
等同于`go build -o scan`

2. 编译成linux可运行的程序
```bash
make default
```