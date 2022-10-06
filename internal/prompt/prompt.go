package prompt

import (
	"github.com/manifoldco/promptui"
)

type Prompt struct {
}

func NewPrompt() *Prompt {
	return &Prompt{}
}

// SelectWithAdd prompt the user to select one of the specified options or write its own
func (p *Prompt) SelectWithAdd(label string, AddLabel string, items []string) (int, string, error) {
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
