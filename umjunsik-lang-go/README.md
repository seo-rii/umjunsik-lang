# Go interpreter

This directory contains the Go v2 interpreter. The module uses only the Go
standard library and retains its existing module path and Go 1.17 minimum.

## Build and run

From the repository root:

```sh
cd umjunsik-lang-go
go build -trimpath -o /tmp/umjunsik-lang-go .
/tmp/umjunsik-lang-go ../examples/helloworld.umm
```

The executable accepts a source-file path and reads integer input from stdin.

## Source-line, input, and arithmetic behavior

- Every source-line slot is preserved, including the header, consecutive blank
  lines, and a blank line following an empty assignment. LF, CRLF, and `~`
  separators keep the same numbering. Blank lines are no-op AST entries; the
  footer terminates parsing as EOF.
- The one-based `준` target is translated using `target - 2` before the
  evaluator's loop increment. This is paired with preserving source-line slots.
- `식?` reads whitespace-separated integers, including several integers on the
  same input line.
- Integer suffixes count each `.` and `,` once. The previous implementation
  counted the first sign twice and omitted the last sign.

For example, this program prints `4`, not `6`:

```text
어떻게
엄....
식어.,!
이 사람이름이냐ㅋㅋ
```

These are corrections to existing language behavior, not a new language
version. There is no instruction limit or forced program termination.
Other language implementations in the repository are unchanged.

## Tests

```sh
cd umjunsik-lang-go  # from the repository root
go test -count=1 -timeout=120s ./...
go vet ./...
```

The integration package builds the current CLI once into a temporary directory,
executes it as a subprocess, and checks exact stdout, empty stderr, exit status,
and a two-second deadline for each finite program. No external checkout,
separately installed interpreter, environment variable, or build tag is needed.
An intentionally infinite program must still require the caller to terminate it.

An optional absolute `UHMLANG_BINARY` path selects a prebuilt executable. This is
useful for negative controls and for checking the actual race-instrumented CLI,
not just the test driver:

```sh
# Run inside umjunsik-lang-go on a platform with Go race-detector support.
go build -race -o /tmp/umjunsik-lang-go-race .
GORACE=atexit_sleep_ms=0 UHMLANG_BINARY=/tmp/umjunsik-lang-go-race \
  go test -race -count=1 -timeout=120s ./...
```

Coverage includes parser line-slot assertions, blank-line diagnostics, forward,
backward, conditional and variable jump targets, empty assignments, integer
input, program exit, repeated addition with nine blank-line padding patterns,
and all 511 period/comma suffixes of length 0 through 8. The CLI arithmetic
matrix checks 18,396 results across three initial values, four expression
contexts and three separators. Parser tests additionally check the offsets of
all 510 nonempty suffixes after two variable forms and an input expression.
The expected offsets come from sign counts, not a copy of the parsing algorithm.

The CLI regressions are adapted from tests developed for the aonohako judge
runtime integration. They use independently written programs; no submitted
user solution or private judge data is included.
