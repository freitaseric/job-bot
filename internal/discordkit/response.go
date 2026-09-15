package discordkit

import "github.com/bwmarrin/discordgo"

func replyInteraction(
	session *discordgo.Session,
	interaction *discordgo.Interaction,
	ephemeral bool,
	components ...discordgo.MessageComponent,
) error {
	var flags discordgo.MessageFlags

	if ephemeral {
		flags |= discordgo.MessageFlagsEphemeral
	}

	return session.InteractionRespond(
		interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:      flags | discordgo.MessageFlagsIsComponentsV2,
				Components: components,
			},
		},
	)
}

func deferInteraction(
	session *discordgo.Session,
	interaction *discordgo.Interaction,
	ephemeral bool,
) error {
	var flags discordgo.MessageFlags

	if ephemeral {
		flags |= discordgo.MessageFlagsEphemeral
	}

	return session.InteractionRespond(
		interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: flags,
			},
		},
	)
}

func editInteraction(
	session *discordgo.Session,
	interaction *discordgo.Interaction,
	components ...discordgo.MessageComponent,
) error {
	_, err := session.InteractionResponseEdit(
		interaction,
		&discordgo.WebhookEdit{
			Flags:      discordgo.MessageFlagsIsComponentsV2,
			Components: &components,
		},
	)

	return err
}

func followupInteraction(
	session *discordgo.Session,
	interaction *discordgo.Interaction,
	components ...discordgo.MessageComponent,
) error {
	_, err := session.FollowupMessageCreate(
		interaction,
		true,
		&discordgo.WebhookParams{
			Flags:      discordgo.MessageFlagsIsComponentsV2,
			Components: components,
		},
	)

	return err
}
