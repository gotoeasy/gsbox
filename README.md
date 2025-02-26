# gsbox

一个关于 3d gaussian splatting 的小工具盒。<br>
`.ply`和`.splat`之间的转换有nodejs版本也有python版本，都太重了，故有此一举。<br>
<br>
<p align="center">
写一写，加深理解<br>
弄一弄，争取实用<br>
搞搞搞，越做越好
<p>

<br>

<p align="center">
    <a href="https://github.com/gotoeasy/gsbox/releases/latest"><img src="https://img.shields.io/github/release/gotoeasy/gsbox.svg"></a>
    <a href="https://github.com/gotoeasy/gsbox/blob/master/LICENSE"><img src="https://img.shields.io/github/license/gotoeasy/gsbox"></a>
<p>

## 功能
- [x] 文件格式的装换，支持3dgs的ply、splat
- [x] 查看ply的文件头信息


## 命令示例
```shell
Usage:
  gsbox [options]

Options:
  ply2splat                ply转splat，可省略根据输入输出文件名自动识别
  splat2ply                splat转ply，可省略根据输入输出文件名自动识别
  simple-ply               简单模式，写ply时不输出未使用字段
  info <plyfile>           显示ply文件头信息
  -i, --input <file>       指定输入文件（注意目录分隔符写法，window平台使用反斜杠的话建议用\\）
  -o, --output <file>      指定输出文件（注意目录分隔符写法，window平台使用反斜杠的话建议用\\）
  -c, --comment <text>     在ply文件头中要写入的注释（有空格等特殊字符时注意引号引起来）
  -h, --help               显示帮助信息
  -v, --version            显示版本信息


# 把3dgs的ply转成splat
gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat

# 把splat转成3dgs的ply
gsbox splat2ply -i /path/to/input.splat -o /path/to/output.ply

# 也支持简化写法，按后缀名自动识别格式
gsbox -i /path/to/input.ply -o /path/to/output.splat

# 查看ply的文件头信息
gsbox info file.ply
```

## 试用版（Trial-2025-v2.0.0）
```shell
# The trial version supports conversion to the .bin format,
# which reduces the file size by 37.5% compared to the .splat format.
gsbox p2b -i /path/to/input.ply -o /path/to/output.bin
gsbox s2b -i /path/to/input.splat -o /path/to/output.bin
```
The `.bin` format can be opened with this 3DGS viewer: https://github.com/reall3d-com/Reall3dViewer


## 更新履历、二进制执行文件下载
https://github.com/gotoeasy/gsbox/releases
