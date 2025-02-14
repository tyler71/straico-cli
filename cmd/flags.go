package cmd

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"log"
	"os"
	p "straico-cli.tylery.com/m/v2/prompt"
	"strings"
)

type ConfigStruct struct {
	FlagMessage string
	Prompt      p.Prompt
}

// Config rootCmd represents the base command when called without any subcommands
var (
	Config          ConfigStruct
	model           string
	apiKey          string
	saveConfig      bool
	saveModel       bool
	listModels      bool
	informationOnly bool
	youtubeYourls   *[]string
	fileUrls        *[]string
)

func Init() {
	flag.BoolVar(&saveModel, "save-model", false, "Straico API key")
	flag.StringVarP(&model, "model", "m", "anthropic/claude-3-haiku:beta", "Model to use")
	youtubeYourls = flag.StringSlice("youtube-url", nil, "--youtube-url link1 --youtube-url link2")
	fileUrls = flag.StringSlice("file-url", nil, "--file-url link1 --file-url link2")
	flag.BoolVarP(&listModels, "list-models", "l", false, "List models")
	flag.StringVar(&apiKey, "save-key", "", "Straico API key")
	flag.Parse()

	modelFlagModified := flag.Lookup("model").Changed

	configFile, _ := LoadConfig()
	if saveModel == true && modelFlagModified {
		saveConfig = true
		informationOnly = true
		configFile.Model = model
	}
	if apiKey != "" {
		saveConfig = true
		configFile.Key = apiKey
	}

	if saveConfig {
		err := SaveConfig(configFile)
		if err != nil {
			log.Println("Unable to save config file")
		} else {
			os.Stdout.Write([]byte("config saved\n"))
		}
	}

	if listModels {
		informationOnly = true
		models, err := GetModels(configFile.Key)
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}

		for _, m := range models {
			outputString := fmt.Sprintf("%s\n\tModel: %s\n\tPricing: %d\n", m.Name, m.Id, m.Pricing)
			os.Stdout.Write([]byte(outputString))
		}
	}
	// If information only, we exit after displaying the info
	if informationOnly {
		os.Exit(0)
	}

	Config.FlagMessage = strings.Join(flag.Args(), " ")
	if len(Config.FlagMessage) == 0 {
		if apiKey == "" {
			os.Stderr.Write([]byte("Error: You must provide a message to send.\n"))
		}
		os.Exit(0)
	}
	Config.Prompt.Model = []string{model}
	Config.Prompt.YoutubeUrls = *youtubeYourls
	Config.Prompt.FileUrls = *fileUrls
}
