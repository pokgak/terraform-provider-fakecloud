// terraform-provider-fakecloud is a Terraform provider for fakecloud, a
// pretend cloud built for learning Terraform (and playing tic-tac-toe
// against a friend, one apply at a time).
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/pokgak/terraform-provider-fakecloud/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs
// for the registry/website:
//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name fakecloud

var version = "dev"

func main() {
	debug := flag.Bool("debug", false, "run the provider in debug mode")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/pokgak/fakecloud",
		Debug:   *debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
