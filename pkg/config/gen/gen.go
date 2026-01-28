package main

import (
	cfg "github.com/ConductorOne/baton-carta/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("carta", cfg.Config)
}
