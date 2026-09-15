package discordkit

import "github.com/bwmarrin/discordgo"

type Context struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	Options     []*discordgo.ApplicationCommandInteractionDataOption
}

func (c *Context) Reply(components ...discordgo.MessageComponent) error {
	return replyInteraction(
		c.Session,
		c.Interaction.Interaction,
		false,
		components...,
	)
}

func (c *Context) Ephemeral(components ...discordgo.MessageComponent) error {
	return replyInteraction(
		c.Session,
		c.Interaction.Interaction,
		true,
		components...,
	)
}

func (c *Context) Defer() error {
	return deferInteraction(
		c.Session,
		c.Interaction.Interaction,
		false,
	)
}

func (c *Context) DeferEphemeral() error {
	return deferInteraction(
		c.Session,
		c.Interaction.Interaction,
		true,
	)
}

func (c *Context) Edit(
	components ...discordgo.MessageComponent,
) error {
	return editInteraction(
		c.Session,
		c.Interaction.Interaction,
		components...,
	)
}

func (c *Context) Followup(
	components ...discordgo.MessageComponent,
) error {
	return followupInteraction(
		c.Session,
		c.Interaction.Interaction,
		components...,
	)
}

func (c *Context) Option(
	name string,
) (*discordgo.ApplicationCommandInteractionDataOption, bool) {
	for _, option := range c.Options {
		if option.Name == name {
			return option, true
		}
	}

	return nil, false
}

func (c *Context) String(
	name string,
) (string, bool) {
	option, ok := c.Option(name)
	if !ok {
		return "", false
	}

	return option.StringValue(), true
}

func (c *Context) Bool(
	name string,
) (bool, bool) {
	option, ok := c.Option(name)
	if !ok {
		return false, false
	}

	return option.BoolValue(), true
}
