package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pterm/pterm"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func viewSelectEntryType() models.EntryType {
	mapEntryTypes := map[string]models.EntryType{
		models.CredentialsType.String(): models.CredentialsType,
		models.TextType.String():        models.TextType,
	}

	keys := make([]string, 0, len(mapEntryTypes))
	for key := range mapEntryTypes {
		keys = append(keys, key)
	}

	selectedType, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Select secret type").
		WithOptions(keys).Show()

	return mapEntryTypes[selectedType]
}

func viewCreateEntry(t models.EntryType) (models.Entry, error) {
	var entry models.Entry

	entryName, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Name").
		WithMultiLine(false).Show()
	switch t {
	case models.CredentialsType:
		var username, password string
		username, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Username").
			WithMultiLine(false).Show()
		password, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Password").
			WithMultiLine(false).WithMask("*").Show()
		entry = models.NewCredentials(entryName, username, password)
	case models.TextType:
		var value string
		value, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Value").
			WithMultiLine(true).Show()
		entry = models.NewText(entryName, value)
	default:
		return nil, models.ErrUnknownSecretType
	}

	return entry, nil
}

func viewEntryOutput(entry models.Entry, jsonOutput bool) error {
	var (
		entryName   string
		outputItems []string
	)

	data, err := entry.Bytes()
	if err != nil {
		return err
	}

	if jsonOutput {
		var prettyJSON bytes.Buffer

		if err = json.Indent(&prettyJSON, data, "", "    "); err != nil {
			return err
		}

		pterm.Println(prettyJSON.String())
		return nil
	}

	switch entry.Type() {
	case models.CredentialsType:
		var secret models.Credentials

		err = json.Unmarshal(data, &secret)
		if err != nil {
			return err
		}

		entryName = secret.EntryName

		outputItems = []string{
			fmt.Sprintf("ID: %s", secret.EntryID),
			fmt.Sprintf("Name: %s", secret.EntryName),
			fmt.Sprintf("Username: %s", secret.Username),
			fmt.Sprintf("Password: %s", secret.Password),
		}
	case models.TextType:
		var secret models.Text

		err = json.Unmarshal(data, &secret)
		if err != nil {
			return err
		}

		entryName = secret.EntryName

		outputItems = []string{
			fmt.Sprintf("ID: %s", secret.EntryID),
			fmt.Sprintf("Name: %s", secret.EntryName),
			fmt.Sprintf("Value: %s", secret.Value),
		}
	default:
		return models.ErrUnknownSecretType
	}

	panel := pterm.DefaultBox.WithTitle(entryName).Sprint(strings.Join(outputItems, "\n"))
	return pterm.DefaultPanel.WithPanels([][]pterm.Panel{{{Data: panel}}}).Render()
}

func viewUpdateEntry(entry models.Entry) (models.Entry, error) {
	var newEntry models.Entry

	data, err := entry.Bytes()
	if err != nil {
		return nil, err
	}

	switch entry.Type() {
	case models.CredentialsType:
		var secret models.Credentials

		err = json.Unmarshal(data, &secret)
		if err != nil {
			return nil, err
		}
		name, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue(secret.EntryName).WithDefaultText("Name").
			WithMultiLine(false).Show()
		username, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue(secret.Username).WithDefaultText("Username").
			WithMultiLine(false).Show()
		password, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue(secret.Password).WithDefaultText("Password").
			WithMultiLine(false).Show()

		secret.EntryName = name
		secret.Username = username
		secret.Password = password
		newEntry = &secret
	case models.TextType:
		var secret models.Text

		err = json.Unmarshal(data, &secret)
		if err != nil {
			return nil, err
		}

		name, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue(secret.EntryName).WithDefaultText("Name").
			WithMultiLine(false).Show()
		value, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue(secret.Value).WithDefaultText("Value").
			WithMultiLine(true).Show()

		secret.EntryName = name
		secret.Value = value
		newEntry = &secret

	default:
		return nil, models.ErrUnknownSecretType
	}

	return newEntry, nil
}

func viewListEntry(entries []models.Entry) error {
	entriesTableData := pterm.TableData{
		{"№", "ID", "Name"},
	}

	for idx, entry := range entries {
		data, err := entry.Bytes()
		if err != nil {
			return err
		}

		entryData := []string{fmt.Sprintf("%d", idx+1)}

		switch entry.Type() {
		case models.CredentialsType:
			var secret models.Credentials
			err = json.Unmarshal(data, &secret)
			if err != nil {
				return err
			}

			entryData = append(entryData, secret.ID().String(), secret.EntryName)

		case models.TextType:
			var secret models.Text
			err = json.Unmarshal(data, &secret)
			if err != nil {
				return err
			}

			entryData = append(entryData, secret.ID().String(), secret.EntryName)

		default:
			return models.ErrUnknownSecretType
		}

		entriesTableData = append(entriesTableData, entryData)
	}

	return pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(entriesTableData).Render()
}

func viewConflictResolve(entries []models.Entry) (int, error) {
	if len(entries) != 2 {
		return 0, errors.New("invalid conflict entries count")
	}

	options := []string{
		"local", "remote",
	}

	panels := make([][]pterm.Panel, len(entries))

	for idx, name := range options {
		var outputItems []string

		data, err := entries[idx].Bytes()
		if err != nil {
			return 0, err
		}

		switch entries[idx].Type() {
		case models.CredentialsType:
			var secret models.Credentials

			err = json.Unmarshal(data, &secret)
			if err != nil {
				return 0, err
			}

			outputItems = []string{
				fmt.Sprintf("ID: %s", secret.EntryID),
				fmt.Sprintf("Name: %s", secret.EntryName),
				fmt.Sprintf("Username: %s", secret.Username),
				fmt.Sprintf("Password: %s", secret.Password),
			}

			panel := pterm.DefaultBox.WithTitle(name).
				Sprint(strings.Join(outputItems, "\n"))
			panels[idx] = []pterm.Panel{{Data: panel}}

		case models.TextType:
			var secret models.Text

			err = json.Unmarshal(data, &secret)
			if err != nil {
				return 0, err
			}

			outputItems = []string{
				fmt.Sprintf("ID: %s", secret.EntryID),
				fmt.Sprintf("Name: %s", secret.EntryName),
				fmt.Sprintf("Value: %s", secret.Value),
			}

			panel := pterm.DefaultBox.WithTitle(name).
				Sprint(strings.Join(outputItems, "\n"))
			panels[idx] = []pterm.Panel{{Data: panel}}

		default:
			return 0, models.ErrUnknownSecretType
		}
	}

	pterm.Warning.Println("Need resolve sync conflict")

	err := pterm.DefaultPanel.WithPanels(panels).Render()
	if err != nil {
		return 0, err
	}
	selectedOption, err := pterm.DefaultInteractiveSelect.WithDefaultText("Select secret to save").
		WithOptions(options).Show()
	if err != nil {
		return 0, err
	}

	for idx, option := range options {
		if option == selectedOption {
			return idx, nil
		}
	}

	return 0, models.ErrNotFound
}
