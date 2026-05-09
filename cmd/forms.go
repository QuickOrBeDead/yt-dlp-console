package cmd

import (
	"fmt"

	"charm.land/huh/v2"
)

var defaultForms FormProvider = RealFormProvider{}

type FormProvider interface {
	Input(title string, validate func(string) error) (string, error)
	InputPassword(title string) (string, error)
	Select(title string, options []string) (string, error)
	Confirm(title, description string) (bool, error)
}

type RealFormProvider struct{}

func (RealFormProvider) Input(title string, validate func(string) error) (string, error) {
	var val string
	f := huh.NewInput().Title(title).Value(&val)
	if validate != nil {
		f = f.Validate(validate)
	}
	err := runHuh(f)
	return val, err
}

func (RealFormProvider) InputPassword(title string) (string, error) {
	var val string
	err := runHuh(huh.NewInput().
		Title(title).
		EchoMode(huh.EchoModePassword).
		Value(&val))
	return val, err
}

func (RealFormProvider) Select(title string, options []string) (string, error) {
	var val string
	err := runHuh(huh.NewSelect[string]().
		Title(title).
		Options(huh.NewOptions(options...)...).
		Value(&val))
	return val, err
}

func (RealFormProvider) Confirm(title, description string) (bool, error) {
	var val bool
	err := huh.NewConfirm().
		Title(title).
		Description(description).
		Affirmative("Yes").
		Negative("No").
		Value(&val).
		Run()
	return val, err
}

func runHuh(f interface{ Run() error }) error {
	err := f.Run()
	fmt.Print("\r\x1b[K")
	return err
}
