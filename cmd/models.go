package cmd

type Config struct {
	ApiEndpoint string `yaml:"cf_endpoint"`
	OauthToken  string `yaml:"token"`
}
