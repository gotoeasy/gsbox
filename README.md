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
- [x] Conversion between file formats, supporting `.ply`, `.splat`, `.spx`, and `.spz` formats for 3DGS.
- [x] Viewing file header information for `.ply`, `.spx`, and `.spz` files, or the number of Gaussian points in `.splat` files.
- [x] Supports data transformation in RST (Rotation, Scale, Translation) order.
- [x] Supports merging multiple model files into one.


## `.spz`
- The `.spz` format is an open 3DGS model format. Its encoding algorithm is highly commendable, and combined with gzip compression, it can significantly reduce the size of model files without any noticeable loss in visual quality.
- The official open-source repository for the `.spz` format is available at https://github.com/nianticlabs/spz. This format is about 10x smaller than the equivalent PLY format and is offered as open source by Niantic Labs. More details can be found at https://scaniverse.com/spz
- For rendering and viewing `.spz` format models, you can refer to https://github.com/mkkellogg/GaussianSplats3D

## `.spx`
- The `.spx` format is flexible, expandable, and supports proprietary data protection. It incorporates encoding methods from both `.splat` and `.spz` formats and adds block compression processing. It supports progressive loading and is suitable for large file models.
- For detailed information about the `.spx` format, please refer to https://github.com/reall3d-com/Reall3dViewer/blob/main/SPX_EN.md
- To render and view models in the `.spx` format, you can use https://github.com/reall3d-com/Reall3dViewer. This viewer is built on Three.js and supports features such as marking, measurements, and text watermarks.

## Usage
```shell
Usage:
  gsbox [options]

Options:
  p2s, ply2splat            convert ply to splat
  p2x, ply2spx              convert ply to spx
  p2z, ply2spz              convert ply to spz
  p2p, ply2ply              convert ply to ply
  s2p, splat2ply            convert splat to ply
  s2x, splat2spx            convert splat to spx
  s2z, splat2spz            convert splat to spz
  s2s, splat2splat          convert splat to splat
  x2p, spx2ply              convert spx to ply
  x2s, spx2splat            convert spx to splat
  x2z, spx2spz              convert spx to spz
  x2x, spx2spx              convert spx to spx
  z2p, spz2ply              convert spz to ply
  z2s, spz2splat            convert spz to splat
  z2x, spz2spx              convert spz to spx
  z2z, spz2spz              convert spz to spz
  join                      join the input model files into a single output file
  info <file>               display the model file information
  -i,  --input <file>       specify the input file
  -o,  --output <file>      specify the output file
  -c,  --comment <text>     specify the comment for ply/spx output
  -bs, --block-size <num>   specify the block size for spx output (default 20480)
  -sh, --shDegree <num>     specify the SH degree for ply/spx/spz output
  -f1, --flag1 <num>        specify the header flag1 for spx output
  -f2, --flag2 <num>        specify the header flag2 for spx output
  -f3, --flag3 <num>        specify the header flag3 for spx output
  -rx, --rotateX <num>      specify the rotation angle in degrees about the x-axis for transform
  -ry, --rotateY <num>      specify the rotation angle in degrees about the y-axis for transform
  -rz, --rotateZ <num>      specify the rotation angle in degrees about the z-axis for transform
  -s,  --scale <num>        specify a uniform scaling factor(0.01~100.0) for transform
  -tx, --translateX <num>   specify the translation value about the x-axis for transform
  -ty, --translateY <num>   specify the translation value about the y-axis for transform
  -tz, --translateZ <num>   specify the translation value about the z-axis for transform
  -v,  --version            display version information
  -h,  --help               display help information

Examples:
  gsbox ply2splat -i /path/to/input.ply -o /path/to/output.splat
  gsbox s2x -i /path/to/input.splat -o /path/to/output.spx -c "your comment" -bs 10240
  gsbox x2z -i /path/to/input.spx -o /path/to/output.spz -sh 0 -rz 90 -s 0.9 -tx 0.1
  gsbox z2p -i /path/to/input.spz -o /path/to/output.ply -c "your comment"
  gsbox join -i a.ply -i b.splat -i c.spx -i d.spz -o output.spx
  gsbox info -i /path/to/file.spx


# Convert the ply to spx without saving SH coefficients and add custom comments.
gsbox p2x -i /path/to/input.ply -o /path/to/output.spx -c "your comment here" -sh 0

# Inspect the header information of the spx file
gsbox info -i /path/to/file.spx
```


## Update History & binary files
https://github.com/gotoeasy/gsbox/releases
