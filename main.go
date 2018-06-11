package main

import	(
	"log"
	"os"
	"github.com/urfave/cli"
	"github.com/TDabasinskas/drone-passwordstate/plugin"
)

var (
	version = "0.0.0"
)

func main() {
	app := cli.NewApp()
	app.Action = run
	app.Version = version
	app.Name = "PasswordState Plugin"
	app.Description = "Drone plugin for integration with Click Studios Passwordstate"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api_endpoint",
			Usage:  "the PasswordState API endpoint (e.g. https://passwordstate/api/)",
			EnvVar: "PLUGIN_API_ENDPOINT,API_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "api_key",
			Usage:  "an API key, with access to the specified list (e.g. d437b3c2f586b9eaed8b736f95324cd5)",
			EnvVar: "PLUGIN_API_KEY,API_KEY",
		},
		cli.IntFlag{
			Name:   "password_list_id",
			Usage:  "the ID of the password list (e.g 1235)",
			EnvVar: "PLUGIN_PASSWORD_LIST_ID,PASSWORD_LIST_ID",
		},
		cli.IntFlag{
			Name:   "connection_retries",
			Usage:  "number of retries in case of failed connections",
			Value: 3,
			EnvVar: "PLUGIN_CONNECTION_RETRIES,CONNECTION_RETRIES",
		},
		cli.IntFlag{
			Name:   "connection_timeout",
			Usage:  "number of seconds for the connection to timeout",
			Value: 5,
			EnvVar: "PLUGIN_CONNECTION_TIMEOUT,CONNECTION_TIMEOUT",
		},
		cli.BoolFlag{
			Name:   "skip_tls_verify",
			Usage:  "ignore TLS validation errors",
			EnvVar: "PLUGIN_SKIP_TLS_VERIFY,SKIP_TLS_VERIFY",
		},
		cli.StringFlag{
			Name:   "key_field",
			Value: "UserName",
			Usage:  "the field, which will be used as the key of the secret",
			EnvVar: "PLUGIN_KEY_FIELD,KEY_FIELD",
		},
		cli.StringFlag{
			Name:   "value_field",
			Value: "Password",
			Usage:  "the field, which will be used as the value of the secret",
			EnvVar: "PLUGIN_VALUE_FIELD,VALUE_FIELD",
		},
		cli.StringFlag{
			Name:   "output_path",
			Value: "./secrets.yaml",
			Usage:  "the path of the file to export the secrets to",
			EnvVar: "PLUGIN_OUTPUT_PATH,OUTPUT_PATH",
		},
		cli.StringFlag{
			Name:   "output_format",
			Value: "YAML",
			Usage:  "output format (currently, only YAML is supported)",
			EnvVar: "PLUGIN_OUTPUT_FORMAT,OUTPUT_FORMAT",
		},
		cli.StringFlag{
			Name:   "section_name",
			Value: "secrets",
			Usage:  "add the secrets under the specified key",
			EnvVar: "PLUGIN_SECTION_NAME,SECTION_NAME",
		},
		cli.BoolFlag{
			Name:   "encode_secrets",
			Usage:  "should the secrets be encoded with BASE64",
			EnvVar: "PLUGIN_ENCODE_SECRETS,ENCODE_SECRETS",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "enable debug output",
			EnvVar: "PLUGIN_DEBUG,DEBUG",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := plugin.Plugin{
		Config: plugin.Config{
			ApiEndpoint: c.String("api_endpoint"),
			ApiKey: c.String("api_key"),
			PasswordListId: c.Int("password_list_id"),
			ConnectionRetries: c.Int("connection_retries"),
			ConnectionTimeout: c.Int("connection_timeout"),
			KeyField: c.String("key_field"),
			ValueField: c.String("value_field"),
			SkipTlsVerify: c.Bool("skip_tls_verify"),
			OutputPath: c.String("output_path"),
			OutputFormat: c.String("output_format"),
			SectionName: c.String("section_name"),
			EncodeSecrets: c.Bool("encode_secrets"),
			Debug: c.Bool("debug"),
		},
	}

	return plugin.Exec()
}
