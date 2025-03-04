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
- [x] 文件格式的装换，支持`3DGS`的`.ply`、`.splat`、`.sp20`格式
- [x] 查看`.ply`的文件头信息

## 关于`.sp20`格式的说明
- 字段顺序同`.splat`
- 坐标固定编码为各 24 bits，编码算法参考`.spz`
- 缩放参数固定编码为各 8 bits，编码算法参考`.spz`
- `.sp20`格式每个高斯点固定长 20 bytes，`.splat`则为 32 bytes，有效减少 37.5% 大小。为了能够更方便的进行渐进加载，未采用`.spz`的排列压缩方式进一步减少大小
- 注意：采用`.sp20`格式时肉眼基本识别不出渲染差异，适合绝大多数以减少文件大小为目的的使用场景，但并不是用来替代`.splat`，因为`.sp20`是有损编码方式，因此，也并不建议把`.sp20`转换回`.splat`或`.ply`
- `.sp20`格式可以使用这个渲染器查看 https://github.com/reall3d-com/Reall3dViewer


## 命令示例
```shell
Usage:
  gsbox [options]

Options:
  p2s, ply2splat           convert ply to splat
  p2s20, ply2splat20       convert ply to splat20
  s2p, splat2ply           convert splat to ply
  s2s20, splat2splat20     convert splat to splat20
  simple-ply               simple mode to write ply
  info <plyfile>           display the ply header
  -i, --input <file>       specify the input file
  -o, --output <file>      specify the output file
  -c, --comment <text>     output ply with comment
  -v, --version            display version information
  -h, --help               display help information

Examples:
  gsbox p2s -i /path/to/input.ply -o /path/to/output.splat
  gsbox p2s20 -i /path/to/input.ply -o /path/to/output.sp20
  gsbox s2p -i /path/to/input.splat -o /path/to/output.ply
  gsbox s2s20 -i /path/to/input.splat -o /path/to/output.sp20
  gsbox s2p -i /path/to/input.splat -o /path/to/output.ply simple-ply
  gsbox s2p -i /path/to/input.splat -o /path/to/output.ply -c "your comment"
  gsbox info -i /path/to/file.ply


# 把3dgs的ply转成splat
gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat

# 把3dgs的ply转成splat20
gsbox ply2splat20 -i /path/to/input.ply -o /path/to/output.sp20

# 把splat转成3dgs的ply
gsbox splat2ply -i /path/to/input.splat -o /path/to/output.ply

# 把splat转成3dgs的ply
gsbox splat2splat20 -i /path/to/input.splat -o /path/to/output.sp20

# 查看ply的文件头信息
gsbox info file.ply
```


## 更新履历、二进制执行文件下载
https://github.com/gotoeasy/gsbox/releases
