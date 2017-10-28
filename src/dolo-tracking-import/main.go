package main

import (
	"dolo-tracking-import/context"
	"dolo-tracking-import/logger"
	"errors"
	"fmt"
	"os"
)

func newConfiguration() (*context.Configuration, error) {
	return &context.Configuration{}, errors.New("FIXME read env vars")
}

func buildContext() (*context.App, error) {
	var (
		err    error
		config *context.Configuration
		ctx    *context.App
	)

	if config, err = newConfiguration(); err != nil {
		return ctx, err
	}

	ctx = &context.App{
		Config: *config,
	}

	return ctx, nil
}

func main() {
	var (
		err error
		ctx *context.App
	)

	// Build the app context
	if ctx, err = buildContext(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println("TODO: shits", ctx.Config.Hubspot.APIKey)
}
