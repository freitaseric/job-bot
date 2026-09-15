package discordkit

import "github.com/bwmarrin/discordgo"

type ButtonOption func(*discordgo.Button)

// ButtonEmoji adiciona um emoji.
func ButtonEmoji(name string) ButtonOption {
	return func(button *discordgo.Button) {
		button.Emoji = &discordgo.ComponentEmoji{
			Name: name,
		}
	}
}

// ButtonDisabled desabilita o botão.
func ButtonDisabled() ButtonOption {
	return func(button *discordgo.Button) {
		button.Disabled = true
	}
}

// Button cria um botão interativo.
func Button(
	label string,
	customID string,
	style discordgo.ButtonStyle,
	options ...ButtonOption,
) discordgo.Button {
	button := discordgo.Button{
		Label:    label,
		CustomID: customID,
		Style:    style,
	}

	for _, option := range options {
		option(&button)
	}

	return button
}

func LinkButton(
	label string,
	url string,
	options ...ButtonOption,
) discordgo.Button {
	button := discordgo.Button{
		Label: label,
		URL:   url,
		Style: discordgo.LinkButton,
	}

	for _, option := range options {
		option(&button)
	}

	return button
}

func SelectOption(
	label string,
	value string,
	description string,
) discordgo.SelectMenuOption {
	return discordgo.SelectMenuOption{
		Label:       label,
		Value:       value,
		Description: description,
	}
}

type SelectOptionConfig func(*discordgo.SelectMenu)

func SelectPlaceholder(
	placeholder string,
) SelectOptionConfig {
	return func(selectMenu *discordgo.SelectMenu) {
		selectMenu.Placeholder = placeholder
	}
}

func SelectRange(
	minimum int,
	maximum int,
) SelectOptionConfig {
	return func(selectMenu *discordgo.SelectMenu) {
		selectMenu.MinValues = &minimum
		selectMenu.MaxValues = maximum
	}
}

func StringSelect(
	customID string,
	options []discordgo.SelectMenuOption,
	config ...SelectOptionConfig,
) discordgo.SelectMenu {
	menu := discordgo.SelectMenu{
		MenuType: discordgo.StringSelectMenu,
		CustomID: customID,
		Options:  options,
	}

	for _, option := range config {
		option(&menu)
	}

	return menu
}
