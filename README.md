# gsbox

一个关于 3d gaussian splatting 的小工具盒。<br>
`.ply`和`.splat`之间的转换有nodejs版本也有python版本，都太重了，故由此一举。<br>
<br>
<p align="center">
写一写，加深理解<br>
弄一弄，争取实用<br>
搞搞搞，越来越好
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
# 把3dgs的ply转成splat
gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat

# 把splat转成3dgs的ply
gsbox splat2ply -i /path/to/input.splat -o /path/to/output.ply

# 也支持简化写法，按后缀名自动识别格式
gsbox -i /path/to/input.ply -o /path/to/output.splat

# 查看ply的文件头信息
gsbox info file.ply
```
