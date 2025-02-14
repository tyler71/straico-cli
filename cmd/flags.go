package cmd

import (
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
	Config        ConfigStruct
	model         string
	apiKey        string
	youtubeYourls *[]string
	fileUrls      *[]string
)

func Init() {
	flag.StringVar(&model, "model", "anthropic/claude-3-haiku:beta", "Model to use")
	youtubeYourls = flag.StringSlice("youtube-url", nil, "--youtube-url link1 --youtube-url link2")
	fileUrls = flag.StringSlice("file-url", nil, "--file-url link1 --file-url link2")
	flag.StringVar(&apiKey, "save-key", "", "Straico API key")
	flag.Parse()

	configFile, _ := LoadConfig()
	if apiKey != "" {
		configFile.Key = apiKey
		err := SaveConfig(configFile)
		if err != nil {
			log.Println("Unable to save config file")
		} else {
			os.Stdout.Write([]byte("config saved\n"))
		}
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
