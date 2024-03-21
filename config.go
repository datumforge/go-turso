package turso

// Config is the configuration for the turso client
type Config struct {
	// Token is the token used to authenticate with the turso API
	Token string `json:"token" koanf:"token" jsonschema:"required"`
	// BaseURL is the base URL for the turso API
	BaseURL string `json:"baseUrl" koanf:"baseUrl" jsonschema:"required" default:"https://api.turso.tech"`
	// OrgName is the name of the organization to use for the turso API
	OrgName string `json:"orgName" koanf:"orgName" jsonschema:"required"`
}
