package discordkit

import "github.com/bwmarrin/discordgo"

// Row agrupa componentes interativos.
func Row(
	components ...discordgo.MessageComponent,
) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: components,
	}
}

// Section cria uma seção de texto com um acessório.
//
// O accessory normalmente é um Button ou Thumbnail.
func Section(
	accessory discordgo.MessageComponent,
	components ...discordgo.MessageComponent,
) discordgo.Section {
	return discordgo.Section{
		Components: components,
		Accessory:  accessory,
	}
}

type ContainerOption func(*discordgo.Container)

// WithColor define a barra lateral do Container.
func WithColor(color int) ContainerOption {
	return func(container *discordgo.Container) {
		container.AccentColor = &color
	}
}

// WithSpoiler torna todo o Container um spoiler.
func WithSpoiler() ContainerOption {
	return func(container *discordgo.Container) {
		container.Spoiler = true
	}
}

// Container cria um container V2.
func Container(
	components []discordgo.MessageComponent,
	options ...ContainerOption,
) discordgo.Container {
	container := discordgo.Container{
		Components: components,
	}

	for _, option := range options {
		option(&container)
	}

	return container
}

// Separator cria um divisor padrão.
func Separator() discordgo.Separator {
	return discordgo.Separator{}
}

// LargeSeparator cria um divisor com espaçamento maior.
func LargeSeparator() discordgo.Separator {
	return discordgo.Separator{
		Spacing: new(discordgo.SeparatorSpacingSizeLarge),
	}
}

// Spacer cria espaçamento vertical sem linha divisória.
func Spacer() discordgo.Separator {
	return discordgo.Separator{
		Divider: new(false),
	}
}

// LargeSpacer cria espaçamento vertical maior sem linha divisória.
func LargeSpacer() discordgo.Separator {
	return discordgo.Separator{
		Divider: new(false),
		Spacing: new(discordgo.SeparatorSpacingSizeLarge),
	}
}
