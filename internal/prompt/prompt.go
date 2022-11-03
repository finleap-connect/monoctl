package prompt

import (
	"github.com/manifoldco/promptui"
)

// SelectWithAdd prompts the user to select one of the specified options or write its own
func SelectWithAdd(label string, AddLabel string, items []string) (int, string, error) {
	prompt := promptui.SelectWithAdd{
		Label:    label,
		Items:    items,
		AddLabel: AddLabel,
	}

	idx, result, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}

	return idx, result, nil
}

// Confirm prompts the user to confirm that they want to proceed
func Confirm(msg string) bool {
	prompt := promptui.Prompt{
		Label:     msg,
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil || result != "y" {
		return false
	}
	return true
}
