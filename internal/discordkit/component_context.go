package discordkit

import "github.com/bwmarrin/discordgo"

type ComponentContext struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate

	params map[string]string
}

func (c *ComponentContext) Param(name string) (string, bool) {
	value, ok := c.params[name]

	return value, ok
}

func (c *ComponentContext) MustParam(
	name string,
) string {
	return c.params[name]
}

func (c *ComponentContext) CustomID() string {
	switch c.Interaction.Type {
	case discordgo.InteractionMessageComponent:
		return c.Interaction.MessageComponentData().CustomID

	case discordgo.InteractionModalSubmit:
		return c.Interaction.ModalSubmitData().CustomID

	default:
		return ""
	}
}

func (c *ComponentContext) UserID() string {
	if c.Interaction.Member != nil &&
		c.Interaction.Member.User != nil {
		return c.Interaction.Member.User.ID
	}

	if c.Interaction.User != nil {
		return c.Interaction.User.ID
	}

	return ""
}

func (c *ComponentContext) Defer() error {
	return deferInteraction(
		c.Session,
		c.Interaction.Interaction,
		false,
	)
}

func (c *ComponentContext) DeferEphemeral() error {
	return deferInteraction(
		c.Session,
		c.Interaction.Interaction,
		true,
	)
}

func (c *ComponentContext) Edit(
	components ...discordgo.MessageComponent,
) error {
	return editInteraction(
		c.Session,
		c.Interaction.Interaction,
		components...,
	)
}

func (c *ComponentContext) Followup(
	components ...discordgo.MessageComponent,
) error {
	return followupInteraction(
		c.Session,
		c.Interaction.Interaction,
		components...,
	)
}
