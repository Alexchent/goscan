# go-scan

基于 Cobra 编写的本地文件扫描与查找工具，本地存储（SQLite），无需外部服务。

功能包括：
- 扫描指定目录（含子目录）的文件并入库：`scan start -d ~/Download`
- 从扫描记录中模糊搜索文件：`scan find -k 关键字`
- 导出所有扫描记录到文本文件：`scan export`
- 把导出的记录导入回程序：`scan import -f {记录文件}`
- 按文件后缀/关键字清理记录：`scan clean --suffix apk` / `scan clean -k 关键字`
- 查找并展示重复文件（支持 Web 页面查看与删除副本）：`scan duplicate -d 目录`
- 快速查询文件：`sscan 关键字`
- 查看版本：`scan -v`

## 安装

通过 `go install` 一键安装 `scan` 与 `sscan` 到 `~/go/bin`：

```bash
go install github.com/Alexchent/goscan/cmd/scan@latest   # 安装 scan
go install github.com/Alexchent/goscan/cmd/sscan@latest  # 安装 sscan
```

等价地，`make install` 会同时安装两个二进制。

> 安装特定版本可把 `@latest` 换成具体 tag，如 `@v1.0.7`。
> 若上次安装过旧版本，重新 `go install` 即可更新（`scan -v` 会显示实际版本号）。
> 安装后确保 `~/go/bin` 已加入 `PATH`。

## 快速开始

```bash
# 1. 初始化工作区（默认 ~/.scan，生成配置文件并安装 zsh 补全）
scan init

# 2. 扫描某个目录，如下载目录
scan start -d ~/Download

# 3. 查找文件
scan find -k 文件名关键字

# 4. 查看扫描出的重复文件
scan duplicate -d ~/Download
```

工作区可用 `-w` 指定其它目录，或用 `-c` 直接指定配置文件。

## 编译

```bash
make default   # 编译 scan 与 sscan 到当前目录
make linux     # 交叉编译 Linux amd64
make install   # 等价于 go install ...@latest，安装到 ~/go/bin
```

## 常用命令示例

```bash
# 扫描目录（含子目录）
scan start -d ~/Download

# 模糊查找
scan find -k report

# 导出全部记录到文本（原命令名 dump 已更名为 export）
scan export

# 从文件导入记录
scan import -f have_save_file_xxx.txt

# 清理记录：按后缀或关键字
scan clean --suffix apk
scan clean -k 临时文件

# 重复文件检测：结果输出到控制台，有重复时会启动
# 仅在 127.0.0.1:8080 的 Web 页面（可查看/删除副本）
scan duplicate -d ~/Download

# 版本信息
scan -v          # 或 scan --version
```
