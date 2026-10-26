# tt-based completion tests

Golden-file tests driving the `example` completions in real headless shells
via [tt](https://github.com/tab-dance/tabdance-test) — replacement for the
VHS GIF tests in `../vhs/` (the docs GIFs there are still rendered by VHS).

## Layout

- `bash.tape`, `zsh.tape`, ... — per-shell setup, `Source`d from every topic
  tape. Inside a `Hide` block (capture paused) it wires the carapace
  integration (`source <(example _carapace <shell>)`), sets an explicit
  minimal prompt and clears the screen, so goldens start from a bare prompt
  on any machine.
- `<topic>/` — one tape per shell (copied from the VHS versions, plus
  `Set Shell` and a quoted `Source` line), next to its `<name>.golden`
  file with the expected terminal content.

Goldens capture carapace behavior for that shell — they are verified, not
hand-written.

## Usage

```sh
go install github.com/tab-dance/tabdance-test/cmd/tt@latest
go build -o /usr/local/bin/example example/   # or any PATH location

# update all goldens in a topic after changing specs/tapes:
tt run docs/src/development/tt/word --update

# verify without updating (CI mode):
tt run docs/src/development/tt/word
```

Running a directory derives the golden path from the tape name
(`<tape>.golden`) and **skips tapes without a golden** — so shells that are
not available locally (or that have no golden yet) are simply skipped.
Commit only goldens that were generated on a machine where the shell exists.

## Determinism notes

- Shell versions matter (menu borders, footers differ between versions):
  goldens are container-canonical, i.e. generate and verify with the same
  shell versions as CI.
- Prompt/banner control lives in the setup tapes; carapace and the shells
  honor `$XDG_CONFIG_HOME`, so run against a clean environment or accept
  the environment as-is since config is typed in.
- Failure diagnostics include the rendered screen and recent PTY traffic
  (run without `--update` to see FAIL plus the diff) when debugging.
