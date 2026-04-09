package main

import (
	"VPMBuilder/cmd"
	"VPMBuilder/internal/log"

	"go.uber.org/zap"
)

func main() {
	defer log.Close()
	err := cmd.Execute()
	if err != nil {
		log.ErrorE("Failed to execute command", zap.Error(err))
	}
}
