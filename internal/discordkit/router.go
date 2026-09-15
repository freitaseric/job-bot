package discordkit

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Handler func(*Context)

type Router struct {
	routes map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Handler),
	}
}

func (r *Router) Command(
	name string,
	handler Handler,
) {
	r.routes[name] = handler
}

func (r *Router) Subcommand(
	command string,
	subcommand string,
	handler Handler,
) {
	key := strings.Join([]string{command, subcommand}, ".")

	r.routes[key] = handler
}

func (r *Router) GroupSubcommand(
	command string,
	group string,
	subcommand string,
	handler Handler,
) {
	key := strings.Join(
		[]string{
			command,
			group,
			subcommand,
		},
		".",
	)

	r.routes[key] = handler
}

func (r *Router) Handle(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	path, options := resolveCommandPath(interaction)

	handler, ok := r.routes[path]
	if !ok {
		return
	}

	handler(&Context{
		Session:     session,
		Interaction: interaction,
		Options:     options,
	})
}

func resolveCommandPath(interaction *discordgo.InteractionCreate) (string, []*discordgo.ApplicationCommandInteractionDataOption) {
	data := interaction.ApplicationCommandData()

	path := data.Name
	options := data.Options

	if len(options) == 0 {
		return path, options
	}

	first := options[0]

	switch first.Type {
	case discordgo.ApplicationCommandOptionSubCommand:
		path += "." + first.Name
		options = first.Options

	case discordgo.ApplicationCommandOptionSubCommandGroup:
		path += "." + first.Name

		if len(first.Options) > 0 {
			sub := first.Options[0]

			path += "." + sub.Name
			options = sub.Options
		}
	}

	return path, options
}
