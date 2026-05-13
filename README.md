# tiffcheck

`tiffcheck` is small Go CLI for TIFF structural checks.
It does **not** decode pixel data.

It reads first IFD and reports:

- tiled vs stripped layout (tile detection via tag `322` / `TileWidth`)
- GeoTIFF metadata presence (tags `33550` / `34735`)

## Build

```bash
go build ./cmd/tiffcheck
```

Binary name: `tiffcheck` (`tiffcheck.exe` on Windows).

### Release

```powershell
git tag "0.0.$(Get-Date -Format 'yyyyMMddHHmmssffff')"
git push --tags
```

```shell
git tag "0.0.$(date -u +'%Y%m%d%H%M%S%4N')"
git push --tags
```

## Usage

```bash
tiffcheck [--help] [--version] <filename>
```

Examples:

```bash
tiffcheck .\scratch\exampletiffs\shapes_lzw_tiled.tif
tiffcheck .\scratch\exampletiffs\shapes_lzw.tif
```

## Output

When tiled:

```text
--- TIFF Structural Analysis ---
Status:  TILED (Indexed for partial loading)
Tile Dim: <value> pixels wide
```

When stripped:

```text
--- TIFF Structural Analysis ---
Status:  STRIPPED (Standard linear loading)
Note:    Not optimized for partial random access.
```

If GeoTIFF tags detected:

```text
Metadata: GeoTIFF tags detected
```

## Exit codes

- `0`: success, help, or version output
- `1`: file open/check failure
- `2`: argument/flag usage error

## Development

```bash
go test ./...
```

No dedicated lint command configured.
