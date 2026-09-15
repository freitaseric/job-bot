package discordkit

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ComponentKind uint8

const (
	ComponentButton ComponentKind = iota + 1
	ComponentStringSelect
	ComponentUserSelect
	ComponentRoleSelect
	ComponentMentionableSelect
	ComponentChannelSelect
	ComponentModal
)

func componentKind(
	interaction *discordgo.InteractionCreate,
) (ComponentKind, error) {
	switch interaction.Type {
	case discordgo.InteractionMessageComponent:
		data := interaction.MessageComponentData()

		switch data.ComponentType {
		case discordgo.ButtonComponent:
			return ComponentButton, nil

		case discordgo.SelectMenuComponent:
			return ComponentStringSelect, nil

		case discordgo.UserSelectMenuComponent:
			return ComponentUserSelect, nil

		case discordgo.RoleSelectMenuComponent:
			return ComponentRoleSelect, nil

		case discordgo.MentionableSelectMenuComponent:
			return ComponentMentionableSelect, nil

		case discordgo.ChannelSelectMenuComponent:
			return ComponentChannelSelect, nil

		default:
			return 0, fmt.Errorf(
				"unsupported message component type: %d",
				data.ComponentType,
			)
		}

	case discordgo.InteractionModalSubmit:
		return ComponentModal, nil

	default:
		return 0, fmt.Errorf(
			"unsupported interaction type: %d",
			interaction.Type,
		)
	}
}

func (k ComponentKind) String() string {
	switch k {
	case ComponentButton:
		return "button"

	case ComponentStringSelect:
		return "string-select"

	case ComponentUserSelect:
		return "user-select"

	case ComponentRoleSelect:
		return "role-select"

	case ComponentMentionableSelect:
		return "mentionable-select"

	case ComponentChannelSelect:
		return "channel-select"

	case ComponentModal:
		return "modal"

	default:
		return "unknown"
	}
}
