# Copilot Instructions

## Build and test commands

- `go build ./cmd/tiffcheck` builds the `tiffcheck` CLI.
- `go test ./...` is full-project test and compile command.
- `go test ./... -run TestName` runs one named test when `_test.go` files exist.
- No dedicated lint command is configured in this repo.

At current HEAD, use `go build ./cmd/tiffcheck` to build the CLI binary from source layout rooted under `cmd/`.

## High-level architecture

- Repo is a small Go CLI with entrypoint at `cmd/tiffcheck/main.go` and core module logic in `pkg/tiffcheck`.
- Program flow is: `main` parses flags (`--help`, `--version`) and validates args, opens file input, then calls `pkg/tiffcheck.Check(reader, opts...)` and prints report lines.
- Tool is structural analyzer, not image decoder. It reads TIFF tags only and never loads pixel data.
- Tiling check is based on first IFD containing tag `322` (`TileWidth`). Presence means tiled; absence means stripped.
- GeoTIFF detection is lightweight. If first IFD contains tag `33550` or `34735`, CLI prints GeoTIFF metadata notice.
- There is currently no `CONTEXT.md` or `docs/adr/`; design intent must be inferred from current CLI behavior.

## Key conventions

- Keep UX CLI-first and stdout-only unless request explicitly changes interface. Current contract is `tiffcheck <filename>`, `tiffcheck --help`, and `tiffcheck --version` plus human-readable report lines.
- Keep flag parsing in `cmd/tiffcheck`; `pkg/tiffcheck` should return data/errors and avoid CLI output concerns.
- Keep filesystem opening (`os.Open`) in `cmd/tiffcheck` adapters; `pkg/tiffcheck` should consume `io`-style readers (`io.ReaderAt` + `io.ReadSeeker`) instead of paths.
- Current logic only inspects first IFD. If you add multi-page support, treat that as behavior change, not refactor.
- TIFF tag IDs are centralized in `pkg/tiffcheck/tags.go`; prefer named constants/predicates over new inline numeric IDs.
- Repo already carries local Copilot skill docs under `.github\skills\`. When adding tests or refactoring structure, align with those docs: test behavior through public interfaces and prefer deeper modules over thin pass-through wrappers.

## Architecture guidance for changes

- Favor a deep module behind the CLI interface: keep `main` as wiring, move TIFF analysis implementation behind a small interface with clear error modes and invariants.
- Keep a seam around `github.com/google/tiff`: library-specific parsing and field access should live in an adapter, so call sites depend on project-level interfaces, not third-party types.
- Use the deletion test before adding a new module: if deleting the module only removes indirection, do not add it; if deleting it would spread complexity across callers, keep it.
- The interface is the test surface: prioritize tests that exercise behavior through module interfaces (for example, `run(args, out)` and analyzer interfaces), not tests coupled to implementation details.
- One adapter means a hypothetical seam; two adapters mean a real seam. Do not introduce an interface seam unless there is at least a production adapter and a test adapter (or another real variant).
- When tag checks expand beyond a couple of call sites, centralize them in a tag vocabulary module (named constants/predicates) to improve leverage and locality.
