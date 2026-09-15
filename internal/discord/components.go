package discord

import (
	"fmt"
	"log/slog"

	"freitaseric.com/job-bot/internal/discordkit"
)

func (b *Bot) registerComponent() error {
	if err := b.componentRouter.Button(
		"/jobs/:id/save",
		b.handleSaveJob,
	); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleSaveJob(
	ctx *discordkit.ComponentContext,
) error {
	jobID, ok := ctx.Param("id")
	if !ok {
		return fmt.Errorf(
			"job id parameter not found",
		)
	}

	slog.Info(
		"saving job",
		"job_id", jobID,
		"user_id", ctx.UserID(),
	)

	return nil
}
