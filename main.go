package main

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/getsynq/connections-looker/internal"
	"github.com/getsynq/connections-looker/model"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	looker "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
	"time"
)

var LookerUrl string
var LookerClientId string
var LookerClientSecret string

var rootCmd = &cobra.Command{
	Use:   "connections-looker",
	Short: "Small utility to collect Looker information which is only available with Admin permissions",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&LookerUrl, "url", "", "Full URL of the Looker instance")
	rootCmd.PersistentFlags().StringVar(&LookerClientId, "client_id", "", "Client ID")
	rootCmd.PersistentFlags().StringVar(&LookerClientSecret, "client_secret", "", "Client Secret")

	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {

		if LookerUrl == "" {
			survey.AskOne(&survey.Input{
				Message: "Full URL of the Looker instance:",
			}, &LookerUrl, survey.WithValidator(internal.UrlValidator))
		}

		if LookerClientId == "" {
			survey.AskOne(&survey.Input{
				Message: "Client ID:",
			}, &LookerClientId, survey.WithValidator(survey.Required))
		}

		if LookerClientSecret == "" {
			survey.AskOne(&survey.Password{
				Message: "Client Secret:",
			}, &LookerClientSecret, survey.WithValidator(survey.Required))
		}

		if LookerUrl == "" || LookerClientId == "" || LookerClientSecret == "" {
			cmd.Help()
			return errors.New("Not all required parameters provided")
		}
		return nil
	}

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {

		settings := rtl.ApiSettings{
			BaseUrl:      LookerUrl,
			ClientId:     LookerClientId,
			ClientSecret: LookerClientSecret,
		}

		session := rtl.NewAuthSession(settings)

		sdk := looker.NewLookerSDK(session)

		allConnections, err := sdk.AllConnections("", nil)
		if err != nil {
			return errors.Wrap(err, "unable to retrieve connections")
		}

		fmt.Printf("Discovered %d database connections\n", len(allConnections))

		results := []*model.Connection{}
		for _, connection := range allConnections {
			if connection.Name == nil {
				continue
			}
			if connection.Database == nil && connection.Schema == nil && connection.Host == nil {
				fmt.Printf("ERROR: no access to the details of connection %s, this command needs to be run with Admin permissions", *connection.Name)
				continue
			}
			connectionJsonBytes, err := json.Marshal(connection)
			if err != nil {
				fmt.Printf("failed to export connection %s", *connection.Name)
				continue
			}

			exportConnection := &model.Connection{}
			err = json.Unmarshal(connectionJsonBytes, exportConnection)
			if err != nil {
				fmt.Printf("failed to export connection %s", *connection.Name)
				continue
			}
			results = append(results, exportConnection)

		}

		if len(results) == 0 {
			return errors.New("No connections were extracted, please check your permissions")
		}

		jsonBytes, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return errors.Wrap(err, "failed to create json")
		}

		fileName := strings.ReplaceAll(fmt.Sprintf("connections-%s.json", time.Now().UTC().Format(time.RFC3339)), ":", "_")
		err = ioutil.WriteFile(fileName, jsonBytes, 0644)
		if err != nil {
			return errors.Wrapf(err, "failed to write file %s", fileName)
		}

		fmt.Printf("File %s created\n", fileName)

		return nil
	}

}

func main() {

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}
