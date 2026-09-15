package discordkit

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ComponentHandler func(*ComponentContext) error

type ComponentRouter struct {
	routes map[ComponentKind][]componentRoute
}

func NewComponentRouter() *ComponentRouter {
	return &ComponentRouter{
		routes: make(map[ComponentKind][]componentRoute),
	}
}

func (r *ComponentRouter) Handle(
	kind ComponentKind,
	pattern string,
	handler ComponentHandler,
) error {
	route, err := newComponentRoute(
		pattern,
		handler,
	)
	if err != nil {
		return err
	}

	if existing, conflict := r.conflictingRoute(kind, route); conflict {
		return fmt.Errorf(
			"component route conflict for %s: %q conflicts with %q",
			kind,
			route.pattern,
			existing.pattern,
		)
	}

	r.routes[kind] = append(
		r.routes[kind],
		route,
	)

	return nil
}

func (r *ComponentRouter) Button(
	pattern string,
	handler ComponentHandler,
) error {
	return r.Handle(
		ComponentButton,
		pattern,
		handler,
	)
}

func (r *ComponentRouter) StringSelect(
	pattern string,
	handler ComponentHandler,
) error {
	return r.Handle(
		ComponentStringSelect,
		pattern,
		handler,
	)
}

func (r *ComponentRouter) UserSelect(
	pattern string,
	handler ComponentHandler,
) error {
	return r.Handle(
		ComponentUserSelect,
		pattern,
		handler,
	)
}

func (r *ComponentRouter) RoleSelect(
	pattern string,
	handler ComponentHandler,
) error {
	return r.Handle(
		ComponentRoleSelect,
		pattern,
		handler,
	)
}

func (r *ComponentRouter) MentionableSelect(
	pattern string,
	handler ComponentHandler,
) error {
	return r.Handle(
		ComponentMentionableSelect,
		pattern,
		handler,
	)
}

func (r *ComponentRouter) ChannelSelect(
	pattern string,
	handler ComponentHandler,
) error {
	return r.Handle(
		ComponentChannelSelect,
		pattern,
		handler,
	)
}

func (r *ComponentRouter) Modal(
	pattern string,
	handler ComponentHandler,
) error {
	return r.Handle(
		ComponentModal,
		pattern,
		handler,
	)
}

func (r *ComponentRouter) find(
	kind ComponentKind,
	path string,
) (
	*componentRoute,
	map[string]string,
	bool,
) {
	routes := r.routes[kind]

	var bestRoute *componentRoute
	var bestParams map[string]string
	bestStaticCount := -1

	for index := range routes {
		route := &routes[index]

		params, matched := route.match(path)
		if !matched {
			continue
		}

		if route.staticCount <= bestStaticCount {
			continue
		}

		bestRoute = route
		bestParams = params
		bestStaticCount = route.staticCount
	}

	if bestRoute == nil {
		return nil, nil, false
	}

	return bestRoute, bestParams, true
}

func (r *ComponentRouter) conflictingRoute(
	kind ComponentKind,
	route componentRoute,
) (*componentRoute, bool) {
	signature := route.signature()

	for index := range r.routes[kind] {
		existing := &r.routes[kind][index]

		if existing.signature() == signature {
			return existing, true
		}
	}

	return nil, false
}

func (r *ComponentRouter) Dispatch(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
) error {
	kind, err := componentKind(interaction)
	if err != nil {
		return err
	}

	customID, err := componentCustomID(interaction)
	if err != nil {
		return err
	}

	route, params, found := r.find(
		kind,
		customID,
	)
	if !found {
		return fmt.Errorf(
			"no component route found for %q",
			customID,
		)
	}

	ctx := &ComponentContext{
		Session:     session,
		Interaction: interaction,
		params:      params,
	}

	return route.handler(ctx)
}

func componentCustomID(
	interaction *discordgo.InteractionCreate,
) (string, error) {
	switch interaction.Type {
	case discordgo.InteractionMessageComponent:
		return interaction.MessageComponentData().CustomID, nil

	case discordgo.InteractionModalSubmit:
		return interaction.ModalSubmitData().CustomID, nil

	default:
		return "", fmt.Errorf(
			"interaction does not contain a component custom id",
		)
	}
}
