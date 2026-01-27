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
)

//go:generate go run ./gen
var Config = field.NewConfiguration([]field.SchemaField{
	AccessToken,
})
