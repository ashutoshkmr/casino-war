package client

import (
	"errors"
	"io"
	"strconv"

	"github.com/manifoldco/promptui"
)

type TieBreakChoice int

const (
	Surrender TieBreakChoice = iota
	Gotowar
)

type PromptOptions int

const (
	Draw PromptOptions = iota
	Quit
)



func ChipsPrompt(writer io.WriteCloser, reader io.ReadCloser) (int, error) {
	//validate input
	validate := func(input string) error {
		val, err := strconv.ParseUint(input, 10, 32)
		if err != nil {
			return errors.New("invalid chips Amount")
		}

		if val < 10 {
			return errors.New("you have to bet a minimun of 10 chips")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter initial chips amount",
		Validate: validate,
		Stdin:    reader,
		Stdout:   writer,
	}
	result, err := prompt.Run()

	if err != nil {
		return 0, err
	}
	resAsInt, _ := strconv.ParseUint(result, 10, 32)
	return int(resAsInt), nil
}

func BetPrompt(writer io.WriteCloser, reader io.ReadCloser) (int, error) {
	//validate input
	validate := func(input string) error {
		val, err := strconv.ParseUint(input, 10, 32)
		if err != nil {
			return errors.New("invalid chips Amount")
		}

		if val < 10 {
			return errors.New("you have to bet a minimun of 10 chips")
		}

		if val > 500 {
			return errors.New("you cannot bet more than 500 chips")
		}

		if val % 2 != 0 {
			// bet amount should be even to handle surrender smoothly
			return errors.New("bet amount should be even")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Amount of chips to bet",
		Validate: validate,
		Stdin:    reader,
		Stdout:   writer,
	}
	result, err := prompt.Run()

	if err != nil {
		return 0, err
	}
	resAsInt, _ := strconv.ParseUint(result, 10, 32)
	return int(resAsInt), nil
}

func InitTiePrompt(writer io.WriteCloser, reader io.ReadCloser) (TieBreakChoice, error) {
	prompt := promptui.Select{
		Label:  "It's a tie what do you want to do ? ",
		Items:  []string{"Surrender", "Go to War"},
		Stdin:  reader,
		Stdout: writer,
	}
	_, result, _ := prompt.Run()
	if result == "Surrender" {
		return Surrender, nil
	} else if result == "Go to War" {
		return Gotowar, nil
	}
	return Surrender, nil
}

func PromptToDraw(writer io.WriteCloser, reader io.ReadCloser) (PromptOptions, error) {
	prompt := promptui.Select{
		Label:  "Draw or Quit ?",
		Items:  []string{"Draw", "Quit"},
		Stdin:  reader,
		Stdout: writer,
	}
	_, result, _ := prompt.Run()
	if result == "Draw" {
		return Draw, nil
	} else if result == "Quit" {
		return Quit, nil
	}
	return Quit, nil
}

