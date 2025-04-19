<p align=center>
<img src="https://gotoeasy.github.io/3dgs/gsbox.png"/>
</p>


# gsbox

一个关于 3d gaussian splatting 的小工具盒。<br>
`.ply`和`.splat`等格式之间的转换有nodejs版本也有python版本，都太重了，故有此一举。<br>
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
- [x] 文件格式之间相互装换，支持`3DGS`的`.ply`、`.splat`、`.spx`、`.spz`格式
- [x] 查看`.ply`、`.spx`、`.spz`的文件头信息或`.splat`的高斯点数量

## `.spz`
- `.spz`是个开放的3DGS模型格式，它的编码算法非常值得称赞，再加上gzip压缩，能有效减少模型文件大小
- `.spz`格式官方开源地址 https://github.com/nianticlabs/spz
- `.spz`格式模型的渲染查看 https://github.com/mkkellogg/GaussianSplats3D

## `.spx`
- `.spx`格式灵活可扩充且支持专有数据保护，其开放格式综合参考了`.splat`和`.spz`的编码方式，并增加了分块压缩处理，支持渐进加载，支持大文件模型
- `.spx`格式说明 https://github.com/reall3d-com/Reall3dViewer/blob/main/SPX_ZH.md
- `.spx`格式模型的渲染查看 https://github.com/reall3d-com/Reall3dViewer

## 命令示例
```shell
Usage:
  gsbox [options]

Options:
  p2s, ply2splat           convert ply to splat
  p2x, ply2spx             convert ply to spx
  p2z, ply2spz             convert ply to spz
  s2p, splat2ply           convert splat to ply
  s2x, splat2spx           convert splat to spx
  s2z, splat2spz           convert splat to spz
  x2p, spx2ply             convert spx to ply
  x2s, spx2splat           convert spx to splat
  x2z, spx2spz             convert spx to spz
  z2p, spz2ply             convert spz to ply
  z2s, spz2splat           convert spz to splat
  z2x, spz2spx             convert spz to spx
  info <file>              display the model file information
  -i, --input <file>       specify the input file
  -o, --output <file>      specify the output file
  -c, --comment <text>     output ply/spx with the comment
  -sh, --shDegree <num>    specify the SH degree for ply/spx/spz output
  -v, --version            display version information
  -h, --help               display help information

Examples:
  gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat
  gsbox s2x -i /path/to/input.splat -o /path/to/output.spx -c "your comment"
  gsbox x2z -i /path/to/input.spx -o /path/to/output.spz -sh 0
  gsbox z2p -i /path/to/input.spz -o /path/to/output.ply -c "your comment" -sh 3
  gsbox info -i /path/to/file.spx


# 把3dgs的ply转成spx并添加自定义说明，不保存球谐系数
gsbox p2x -i /path/to/input.ply -o /path/to/output.spx -c "your comment here" -sh 0

# 查看spx的文件头信息
gsbox info -i /path/to/file.spx
```


## Update History & binary files
https://github.com/gotoeasy/gsbox/releases
