package main

import (
	"fmt"
	"log"

	api "github.com/ALT-F4-LLC/vorpal/sdk/go/pkg/api/artifact"
	"github.com/ALT-F4-LLC/vorpal/sdk/go/pkg/artifact"
	"github.com/ALT-F4-LLC/vorpal/sdk/go/pkg/artifact/language"
	"github.com/ALT-F4-LLC/vorpal/sdk/go/pkg/config"
)

var SYSTEMS = []api.ArtifactSystem{
	api.ArtifactSystem_X8664_LINUX,
}

func gatewayDevenv(context *config.ConfigContext) (*string, error) {
	gobin, err := artifact.GoBin(context)
	if err != nil {
		return nil, err
	}

	artifacts := []*string{
		gobin,
	}

	contextTarget := context.GetTarget()

	goarch, err := language.GetGOARCH(contextTarget)
	if err != nil {
		return nil, err
	}

	goos, err := language.GetGOOS(contextTarget)
	if err != nil {
		return nil, err
	}

	environments := []string{
		"CGO_ENABLED=0",
		fmt.Sprintf("GOARCH=%s", *goarch),
		fmt.Sprintf("GOOS=%s", *goos),
	}

	return artifact.ScriptDevenv(context, artifacts, environments, "devenv", nil, SYSTEMS)
}

func main() {
	context := config.GetContext()
	contextArtifact := context.GetArtifactName()

	var err error

	switch contextArtifact {
	case "devenv":
		_, err = gatewayDevenv(context)
	default:
		log.Fatalf("unknown artifact %s", contextArtifact)
	}
	if err != nil {
		log.Fatalf("failed to build %s: %v", contextArtifact, err)
	}

	context.Run()
}
