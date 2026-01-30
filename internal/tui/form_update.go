package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// updateAddRule handles add rule form updates
func (m *Model) updateAddRule(msg tea.Msg) (tea.Model, tea.Cmd) {
	form := m.addRuleForm

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle backspace to delete characters
		if msg.Type == tea.KeyBackspace {
			var cmd tea.Cmd
			currentField := form.FocusedField()
			currentField.input, cmd = currentField.input.Update(msg)
			// Check if product was modified
			if form.focus == 0 {
				form.CheckProductChange()
			}
			return m, cmd
		}

		switch {
		case key.Matches(msg, keys.Enter):
			if form.ShowingOptions() {
				if form.optionFocus >= 0 {
					form.SelectProductOption(form.optionFocus)
				}
				form.ShowProductOptions()
				return m, nil
			}
			if errMsg, valid := form.ValidateCurrentField(); !valid {
				m.lastError = fmt.Errorf("%s", errMsg)
				m.screen = ScreenError
				return m, nil
			}
			if form.focus == len(form.fields)-1 {
				return m.submitAddRule()
			}
			return m, form.NextField()

		case key.Matches(msg, keys.Tab):
			if form.ShowingOptions() {
				form.ShowProductOptions()
				return m, nil
			}
			// Check for auto-population before moving to next field
			if form.focus == 0 {
				form.CheckProductChange()
			}
			return m, form.NextField()

		case key.Matches(msg, keys.Back):
			form.Reset()
			return m, nil

		case key.Matches(msg, keys.Down):
			if form.IsProductField() && form.ShowingOptions() {
				form.optionFocus++
				if form.optionFocus >= len(form.GetProductOptions()) {
					form.optionFocus = 0
				}
				return m, nil
			}

		case key.Matches(msg, keys.Up):
			if form.IsProductField() && form.ShowingOptions() {
				form.optionFocus--
				if form.optionFocus < 0 {
					form.optionFocus = len(form.GetProductOptions()) - 1
				}
				return m, nil
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
			if form.IsProductField() {
				form.ShowProductOptions()
				if form.ShowingOptions() {
					form.optionFocus = 0
				}
				return m, nil
			}
		}

		// Handle character input with validation
		if msg.Type == tea.KeyRunes {
			currentField := form.FocusedField()
			blocked, _ := currentField.HandleInput(msg)
			if blocked {
				return m, nil
			}
			// Check if product was modified
			if form.focus == 0 {
				form.CheckProductChange()
			}
		}
	}

	// Update focused field
	var cmd tea.Cmd
	currentField := form.FocusedField()
	currentField.input, cmd = currentField.input.Update(msg)

	return m, cmd
}
