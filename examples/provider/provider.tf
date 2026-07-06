provider "fakecloud" {
  # Your playground id, from the dashboard's "Connect Terraform" panel.
  # Falls back to the FAKECLOUD_SANDBOX environment variable.
  sandbox = "abc123xyz0"

  # Defaults to the hosted fakecloud; point it at a local
  # `wrangler dev` instead if you run your own:
  # endpoint = "http://localhost:8787"
}
