# gsbox
一个关于 3d gaussian splatting 的小工具盒

<br>

<p align="center">
    <a href="https://github.com/gotoeasy/gsbox/releases/latest"><img src="https://img.shields.io/github/release/gotoeasy/gsbox.svg"></a>
    <a href="https://github.com/gotoeasy/gsbox/blob/master/LICENSE"><img src="https://img.shields.io/github/license/gotoeasy/gsbox"></a>
<p>

## 功能
- [x] 把3dgs的ply转成splat
- [x] 把splat转成3dgs的ply
- [x] 查看ply的文件头信息
- [ ] TODO

## 命令示例
```shell
# 把3dgs的ply转成splat
gsbox ply2splat from.ply to.splat

# 把splat转成3dgs的ply
gsbox splat2ply from.splat to.ply

# 查看ply的文件头信息
gsbox info from.ply
```

## 更新履历

### 版本`v1.0.0`

- [x] 把3dgs的ply转成splat
- [x] 把splat转成3dgs的ply
- [x] 查看ply的文件头信息
