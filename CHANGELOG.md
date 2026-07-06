# Changelog

## v0.3.0 (2026-07-06)

Complete rework: fakecloud is now a Terraform learning playground built
around one primitive — a tic-tac-toe board.

BREAKING CHANGES:

- Removed `fakecloud_virtual_machine` resource and the
  `fakecloud_virtual_machine`/`fakecloud_virtual_machines` data sources.
- Provider now targets the hosted fakecloud by default and takes a new
  `sandbox` attribute (`FAKECLOUD_SANDBOX` env var) identifying your
  playground; `endpoint` (`FAKECLOUD_ENDPOINT`) still overrides the host.

FEATURES:

- New resource `fakecloud_tictactoe_board`: `freeplay` or server-refereed
  `duel` mode; computed `cells`, `next_player`, `winner`, `nameplate_text`.
- New resource `fakecloud_tictactoe_move`: create plays a mark, destroy
  takes it back, importable by id; duel boards reject out-of-turn applies
  with the referee's reason.
- New resource `fakecloud_nameplate`: one plaque per board, `text` updates
  in place.
- New data source `fakecloud_tictactoe_board` for reading boards you don't
  manage (e.g. your duel opponent's).

## v0.2.7 and earlier

The original virtual-machine-based provider.
