package plugin

import (
	"github.com/go-resty/resty"
	"crypto/tls"
	"time"
	"strings"
	"strconv"
	"reflect"
	"os"
	"fmt"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"github.com/mattn/go-colorable"
	"net/url"
)

type (
	// Plugin configuration
	Config struct {
		ApiEndpoint			string
		ApiKey 				string
		PasswordListId		int
		ConnectionRetries	int
		ConnectionTimeout   int
		SkipTlsVerify  	 	bool
		KeyField			string
		ValueField			string
		EncodeSecrets		bool
		OutputPath			string
		OutputFormat		string
		SectionName			string
		Debug				bool
	}
	// Plugin parameters
	Plugin struct {
		Config Config
	}
	// KV Secret
	Secret struct {
		Key		string
		Value	string
	}
)

// Handles the plugin execution
func (p *Plugin) Exec() error {

	// Initiate the logging
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors:true, FullTimestamp:true})
	logrus.SetOutput(colorable.NewColorableStdout())
	if p.Config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Infoln("Starting the execution.")

	// Validate the parameters:
	_, err := url.ParseRequestURI(p.Config.ApiEndpoint)
	if err != nil {
		logrus.WithError(err).WithField("endpoint", p.Config.ApiEndpoint).Errorln("Provided API endpoint is not valid.")
		return err
	}
	if p.Config.PasswordListId == 0 {
		logrus.WithField("list_id", p.Config.ApiEndpoint).Errorln("Provided list ID is not valid.")
		return nil
	}
	if p.Config.ApiKey == "" {
		logrus.Errorln("API key is mandatory.")
		return nil
	}
	if p.Config.OutputFormat != "YAML" {
		logrus.Errorln("Currently only YAML format is supported.")
		return nil
	}

	// Retrieve the secrets from PasswordState:
	secrets, err := getSecrets(p)
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		logrus.Warnln("Secrets were retrieved from PasswordState, but none of them could be converted to Key-Value pairs. Terminating.")
		return nil;
	}

	// Save the secrets to file:
	if p.Config.OutputFormat == "YAML" {
		outputToYaml(p.Config.OutputPath, p.Config.SectionName, p.Config.EncodeSecrets, secrets)
	}

	logrus.Infoln("Finished the execution.")
	return nil
}

// Saves the secrets to YAML file
func outputToYaml(filename string, section string, encode bool, secrets []Secret) error {
	logrus.WithField("output_path",  filename).Infoln("Writing secrets to the file.")
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		logrus.WithError(err).Errorln("Failed writing secrets to to the file.")
		return err
	}
	f.WriteString(fmt.Sprintf("---\n%s:\n", string(section)))
	for _, secret := range secrets {
		// Encode the secrets if needed:
		var value string
		if encode {
			value = base64.StdEncoding.EncodeToString([]byte(secret.Value))
		} else {
			value = secret.Value
		}

		logrus.WithField("key", secret.Key).WithField("value", "(hidden)").Infoln("Secret saved.")
		f.WriteString(fmt.Sprintf("  %s: '%s'\n", secret.Key, value))
	}

	logrus.WithField("outputPath", filename).WithField("count", len(secrets)).Infoln("Secrets successfully saved to the file.")
	return nil
}

// Retrieves the secrets from PasswordState
func getSecrets(p *Plugin) ([]Secret, error) {
	type (
		// PasswordState JSON response for the passwords
		PasswordList struct {
			PasswordID     int    `json:"PasswordID"`
			Title          string `json:"Title"`
			UserName       string `json:"UserName"`
			Description    string `json:"Description"`
			GenericField1  string `json:"GenericField1"`
			GenericField2  string `json:"GenericField2"`
			GenericField3  string `json:"GenericField3"`
			GenericField4  string `json:"GenericField4"`
			GenericField5  string `json:"GenericField5"`
			GenericField6  string `json:"GenericField6"`
			GenericField7  string `json:"GenericField7"`
			GenericField8  string `json:"GenericField8"`
			GenericField9  string `json:"GenericField9"`
			GenericField10 string `json:"GenericField10"`
			AccountTypeID  int    `json:"AccountTypeID"`
			Notes          string `json:"Notes"`
			URL            string `json:"URL"`
			Password       string `json:"Password"`
			ExpiryDate     string `json:"ExpiryDate"`
			AllowExport    bool   `json:"AllowExport"`
			AccountType    string `json:"AccountType"`
		}
	)

	var (
		url     strings.Builder
		secrets	[]Secret
	)

	url.WriteString(strings.TrimRight(p.Config.ApiEndpoint, "/"))
	url.WriteString("/passwords/{PasswordListID}")

	// Configure the API client:
	client := resty.New()
	client.
		SetRetryCount(p.Config.ConnectionRetries).
		SetTimeout(time.Duration(p.Config.ConnectionTimeout) * time.Second)
	if p.Config.Debug {
		client.SetDebug(true)
	}
	if p.Config.SkipTlsVerify {
		client.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: p.Config.SkipTlsVerify })
	}
	client.
		SetQueryParams(map[string]string{
			"QueryAll": "true",
			"PreventAuditing": "false",
		}).
		SetPathParams(map[string]string{
			"PasswordListID": strconv.Itoa(p.Config.PasswordListId),
		}).
		SetHeaders(map[string]string{
			"APIKey": p.Config.ApiKey,
			"Content-Type": "application/json",
		})

	// Send the request:
	logrus.WithField("endpoint", p.Config.ApiEndpoint).	WithField("list_id", p.Config.PasswordListId).Infoln("Querying PasswordState API.")
	response, err := client.R().
		SetResult([]PasswordList{}).
		Get(url.String())

	if err != nil {
		logrus.WithError(err).Errorln("Failed to retrieved data from PasswordState.")
		return nil,err
	}

	passwords := *response.Result().(*[]PasswordList)
	logrus.WithField("count", len(passwords)).Infoln("Passwords retrieved from PasswordState.")
	logrus.WithField("key_field", p.Config.KeyField).WithField("value_field", p.Config.ValueField).Infoln("Converting retrieved passwords to secrets.")
	for _, password := range passwords {
		key := reflect.Indirect(reflect.ValueOf(password)).FieldByName(p.Config.KeyField).String()
		if key == "" || key == "<invalid Value>" {
			logrus.WithField("password_id", password.PasswordID).WithField("field", p.Config.KeyField).Warnln("Key is empty. Skipping the secret.")
			continue
		}
		value := reflect.Indirect(reflect.ValueOf(password)).FieldByName(p.Config.ValueField).String()
		if value == "" || value == "<invalid Value>" {
			logrus.WithField("password_id", password.PasswordID).WithField("field", p.Config.ValueField).Warnln("Value is empty. Skipping the secret.")
			continue
		}
		secret := Secret{
			Key: key,
			Value: value,
		}
		secrets = append(secrets, secret)
	}

	logrus.WithField("count", len(secrets)).Infoln("Finished processing the secrets.")
	return secrets,nil
}
