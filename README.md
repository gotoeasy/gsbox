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
- [x] 文件格式的装换，支持`3DGS`的`.ply`、`.splat`、`.spx`格式
- [x] 查看`.ply`、`.spx`的文件头信息或`.splat`的高斯点数量

## 关于`.spx`格式的说明
- `.spx`格式灵活可扩充且支持专有数据保护，其开放格式综合参考了`.splat`和`.spz`的编码方式，再增加分块压缩处理支持渐进加载，压缩后的`.spx`不足`.splat`的一半大小。`gsbox`支持这些格式间的相互转换，模型文件的下载和应用已基本无负担
- `.spx`格式说明 https://github.com/reall3d-com/Reall3dViewer/blob/main/SPX_ZH.md
- `.spx`格式模型的渲染查看 https://github.com/reall3d-com/Reall3dViewer

## 命令示例
```shell
Usage:
  gsbox [options]

Options:
  p2s, ply2splat           convert ply to splat
  p2x, ply2spx             convert ply to spx
  s2p, splat2ply           convert splat to ply
  s2x, splat2spx           convert splat to spx
  x2p, spx2ply             convert spx to ply
  x2s, spx2splat           convert spx to splat
  simple-ply               simple mode to write ply
  info <file>              display the model file information
  -i, --input <file>       specify the input file
  -o, --output <file>      specify the output file
  -c, --comment <text>     output ply/spx with the comment
  -v, --version            display version information
  -h, --help               display help information

Examples:
  gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat
  gsbox s2x -i /path/to/input.splat -o /path/to/output.spx -c "your comment"
  gsbox x2p -i /path/to/input.spx -o /path/to/output.ply simple-ply
  gsbox s2p -i /path/to/input.splat -o /path/to/output.ply -c "your comment"
  gsbox info -i /path/to/file.spx


# 把3dgs的ply转成spx并添加自定义说明
gsbox ply2spx -i /path/to/input.ply -o /path/to/output.spx -c "your comment here"

# 查看spx的文件头信息
gsbox info -i /path/to/file.spx
```


## 更新履历、二进制执行文件下载
https://github.com/gotoeasy/gsbox/releases
