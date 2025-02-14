/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"os"
	"straico-cli.tylery.com/m/v2/cmd"
)

func main() {
	cmd.Init()
	configFile, _ := cmd.LoadConfig()
	config := cmd.Config
	straicoResponse, err := config.Prompt.Request(configFile.Key, config.FlagMessage)
	if err != nil {
		log.Fatalln(err)
		return
	}
	_, err = os.Stdout.Write([]byte(straicoResponse.Data.Completions[config.Prompt.Model[0]].Completion.Choices[0].Message.Content))
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
	}
}
