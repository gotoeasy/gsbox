<p align=center>
<img src="https://gotoeasy.github.io/3dgs/gsbox.png"/>
</p>


# gsbox

A cross-platform command-line tool for 3D Gaussian Splatting, focusing on format conversion and optimization.

<br>

<p align="center">
    <a href="https://repo-sam.inria.fr/fungraph/3d-gaussian-splatting/"><img src="https://img.shields.io/badge/model-3DGS-brightgreen.svg"></a>
    <a href="https://github.com/gotoeasy/gsbox/releases/latest"><img src="https://img.shields.io/github/release/gotoeasy/gsbox.svg"></a>
    <a href="https://github.com/gotoeasy/gsbox/blob/master/LICENSE"><img src="https://img.shields.io/github/license/gotoeasy/gsbox"></a>
<p>

## Features
- [x] Conversion between file formats, supporting `.ply`, `.splat`, `.spx`, and `.spz(v2,v3)` formats for 3DGS.
- [x] Viewing file header information for `.ply`, `.spx`, `.spz`, and `.ksplat` files, or simple information of `.splat`.
- [x] Supports data transformation (Rotation, Scale, Translation).
- [x] Supports merging multiple model files into one.


|       | `.ply`   | `.compressed.ply` | `.splat` | `.spx`   | `.spz`  | `.ksplat` | `.sog` |
|-------|----------|-------------------|----------|----------|---------|-----------|-----------|
| Read  | &#9745; |  &#9745;      | &#9745; | &#9745; | &#9745; | &#9745; | &#9745; |
| Write | &#9745; |               | &#9745; | &#9745; | &#9745; |         |         |
| Ref   |  <a href="https://repo-sam.inria.fr/fungraph/3d-gaussian-splatting/">Link</a> |  <a href="https://github.com/playcanvas/supersplat">Link</a> | <a href="https://github.com/antimatter15/splat">Link</a> | <a href="https://github.com/reall3d-com/Reall3dViewer/blob/main/SPX_EN.md">Link</a> | <a href="https://github.com/nianticlabs/spz">Link</a> | <a href="https://github.com/mkkellogg/GaussianSplats3D">Link</a> | <a href="https://github.com/playcanvas/splat-transform">Link</a> |


## `.spx`
- The `.spx` format is flexible, expandable, and supports proprietary data protection. It incorporates encoding methods from both `.splat` and `.spz` formats and adds block compression processing. It supports progressive loading and is suitable for large file models.
- For detailed information about the `.spx` format, please refer to [SPX Specification](https://github.com/reall3d-com/Reall3dViewer/blob/main/SPX_EN.md)
- To render and view models in the `.spx` format, you can use [Reall3dViewer](https://github.com/reall3d-com/Reall3dViewer). This viewer is built on Three.js and supports features such as marking, measurements, and text watermarks.

## Usage
```shell
Usage:
  gsbox [options]

Options:
  p2s, ply2splat                  convert ply to splat
  p2x, ply2spx                    convert ply to spx
  p2z, ply2spz                    convert ply to spz
  p2p, ply2ply                    convert ply to ply
  s2p, splat2ply                  convert splat to ply
  s2x, splat2spx                  convert splat to spx
  s2z, splat2spz                  convert splat to spz
  s2s, splat2splat                convert splat to splat
  x2p, spx2ply                    convert spx to ply
  x2s, spx2splat                  convert spx to splat
  x2z, spx2spz                    convert spx to spz
  x2x, spx2spx                    convert spx to spx
  z2p, spz2ply                    convert spz to ply
  z2s, spz2splat                  convert spz to splat
  z2x, spz2spx                    convert spz to spx
  z2z, spz2spz                    convert spz to spz
  k2p, ksplat2ply                 convert ksplat to ply
  k2s, ksplat2splat               convert ksplat to splat
  k2x, ksplat2spx                 convert ksplat to spx
  k2z, ksplat2spx                 convert ksplat to spz
  g2p, sog2ply                    convert sog to ply
  g2s, sog2splat                  convert sog to splat
  g2x, sog2spx                    convert sog to spx
  g2z, sog2spx                    convert sog to spz
  ps,  printsplat                 print data to text file like splat format layout
  join                            join the input model files into a single output file
  info <file>                     display the model file information
  -i,  --input <file>             specify the input file
  -o,  --output <file>            specify the output file
  -ct, --compression-type <type>  specify the compression type(0:gzip,1:xz) for spx output, default is gzip
  -c,  --comment <text>           specify the comment for ply/spx output
  -a,  --alpha <num>              specify the minimum alpha(0~255) to filter the output splat data
  -bs, --block-size <num>         specify the block size(64~1048576) for spx output (default is 65536)
  -bf, --block-format <num>       specify the block data format(19~20) for spx output (default is 19)
  -sh, --shDegree <num>           specify the SH degree(0~3) for output
  -f1, --is-inverted <bool>       specify the header flag1(IsInverted) for spx output, default is false
  -rx, --rotateX <num>            specify the rotation angle in degrees about the x-axis for transform
  -ry, --rotateY <num>            specify the rotation angle in degrees about the y-axis for transform
  -rz, --rotateZ <num>            specify the rotation angle in degrees about the z-axis for transform
  -s,  --scale <num>              specify a uniform scaling factor(0.001~1000) for transform
  -tx, --translateX <num>         specify the translation value about the x-axis for transform
  -ty, --translateY <num>         specify the translation value about the y-axis for transform
  -tz, --translateZ <num>         specify the translation value about the z-axis for transform
  -to, --transform-order <RST>    specify the transform order (RST/RTS/SRT/STR/TRS/TSR), default is RST
  -ov, --output-version <num>     specify the output versions for spx(1-2, default is 2) and spz(2-3, default is 2)
  -v,  --version                  display version information
  -h,  --help                     display help information

Examples:
  gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat
  gsbox s2x -i /path/to/input.splat -o /path/to/output.spx -c "your comment" -bs 10240 -ct xz
  gsbox x2z -i /path/to/input.spx -o /path/to/output.spz -sh 0 -rz 90 -s 0.9 -tx 0.1 -to TRS
  gsbox z2p -i /path/to/input.spz -o /path/to/output.ply -c "your comment"
  gsbox k2z -i /path/to/input.ksplat -o /path/to/output.spz -ov 3
  gsbox g2x -i /path/to/input.sog -o /path/to/output.spx
  gsbox g2x -i /path/to/meta.json -o /path/to/output.spx
  gsbox join -i a.ply -i b.splat -i c.spx -i d.spz -i e.ksplat -i f.sog -i meta.json -o output.spx
  gsbox ps -i /path/to/input.spx -o /path/to/output.txt
  gsbox info -i /path/to/file.spx

# Convert the ply to spx without saving SH coefficients and add custom comments.
gsbox p2x -i /path/to/input.ply -o /path/to/output.spx -c "your comment here" -sh 0

# Convert the ply to spz version 3.
gsbox p2z -i /path/to/input.ply -o /path/to/output.spz -ov 3

# Inspect the header information of the spx file
gsbox info -i /path/to/file.spx
```

## Update History & binary files
https://github.com/gotoeasy/gsbox/releases
