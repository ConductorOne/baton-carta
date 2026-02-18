package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	AccessToken = field.StringField(
		"token",
		field.WithDescription("The Carta personal access token used to connect to the Carta API."),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)
	BaseURLField = field.StringField(
		"base-url",
		field.WithDescription("Override the Carta API URL (for testing)"),
		field.WithHidden(true),
		field.WithExportTarget(field.ExportTargetCLIOnly),
	)
)

//go:generate go run ./gen
var Config = field.NewConfiguration([]field.SchemaField{
	AccessToken,
	BaseURLField,
})
