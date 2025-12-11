package cmd

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	model           string
	apiKey          string
	saveConfig      bool
	saveModel       bool
	listModels      bool
	informationOnly bool
	youtubeYourls   *[]string
	fileUrls        *[]string
)

func Init() *ConfigFile {
	flag.BoolVar(&saveModel, "save-model", false, "Use the model listed by -m for future queries")
	flag.StringVarP(&model, "model", "m", "anthropic/claude-3-haiku:beta", "Model to use")
	youtubeYourls = flag.StringSlice("youtube-url", nil, "--youtube-url link1 --youtube-url link2")
	fileUrls = flag.StringSlice("file-url", nil, "--file-url link1 --file-url link2")
	flag.BoolVarP(&listModels, "list-models", "l", false, "List models")
	flag.StringVar(&apiKey, "save-key", "", "Straico API key")
	flag.Parse()

	modelFlagModified := flag.Lookup("model").Changed

	configFile := ConfigFile{}
	err := configFile.LoadConfig()
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
	}
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
		err := configFile.SaveConfig()
		if err != nil {
			log.Println("Unable to save config file")
		} else {
			_, _ = os.Stdout.Write([]byte("config saved\n"))
		}
	}

	if listModels {
		informationOnly = true
		models, err := GetModels(configFile.Key)
		if err != nil {
			_, _ = os.Stderr.Write([]byte(err.Error()))
		}

		for _, m := range models {
			outputString := fmt.Sprintf("%s\n\tModel: %s\n\tPricing: %d\n", m.Name, m.Id, m.Pricing)
			_, _ = os.Stdout.Write([]byte(outputString))
		}
	}

	if informationOnly {
		os.Exit(0)
	}

	configFile.Prompt.Model = []string{model}
	configFile.Prompt.YoutubeUrls = *youtubeYourls
	configFile.Prompt.FileUrls = *fileUrls

	return &configFile
}
