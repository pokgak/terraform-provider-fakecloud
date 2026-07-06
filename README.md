# terraform-provider-fakecloud

The Terraform provider for [fakecloud](https://github.com/pokgak/fakecloud) —
a pretend cloud for **learning Terraform on a tic-tac-toe board**. Every mark
on the board is a resource, and fakecloud's live dashboard shows every
`apply`, `destroy`, and drift event as it happens.

Published on the registry as
[`pokgak/fakecloud`](https://registry.terraform.io/providers/pokgak/fakecloud/latest).

```terraform
terraform {
  required_providers {
    fakecloud = {
      source = "pokgak/fakecloud"
    }
  }
}

provider "fakecloud" {
  # Create a playground on the fakecloud website and paste its id here —
  # the dashboard's "Connect Terraform" panel shows the exact block.
  sandbox = "your-sandbox-id"
}
```

## Resources & data sources

- `fakecloud_tictactoe_board` (resource + data source) — a board; `mode` is
  `freeplay` (default) or `duel` (server-refereed: X starts, turns
  alternate, locks on a win). Computed `cells`, `next_player`, `winner`,
  `nameplate_text`.
- `fakecloud_tictactoe_move` — a mark on a board: create = play,
  destroy = take it back, importable by id.
- `fakecloud_nameplate` — a plaque on a board (one per board); `text`
  updates in place.

The full learning course — six chapters of missions taught on this one
primitive, from resource basics through `count` footguns, drift, modules,
shared state, and the dependency graph — lives in the
[fakecloud repo](https://github.com/pokgak/fakecloud).

## Development

Requirements: [Go](https://golang.org/doc/install) >= 1.24.

```sh
go build -o ~/go/bin/terraform-provider-fakecloud .
```

Use a [dev override](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides)
to point Terraform at your local build, and run the fakecloud server locally
with `wrangler dev` (see the fakecloud repo). Regenerate docs after schema
changes with `go generate ./...`.

## Releasing

Push a semver tag (`git tag v0.3.0 && git push --tags`) and the release
workflow builds, signs, and publishes GoReleaser-style; the Terraform
Registry picks the release up automatically.
