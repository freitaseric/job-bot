package discordkit

import "github.com/bwmarrin/discordgo"

// Text cria um Text Display com suporte a Markdown.
func Text(content string) discordgo.TextDisplay {
	return discordgo.TextDisplay{
		Content: content,
	}
}

type ThumbnailOption func(*discordgo.Thumbnail)

// ThumbnailDescription define o texto alternativo da imagem.
func ThumbnailDescription(description string) ThumbnailOption {
	return func(thumbnail *discordgo.Thumbnail) {
		thumbnail.Description = &description
	}
}

// ThumbnailSpoiler transforma a imagem em spoiler.
func ThumbnailSpoiler() ThumbnailOption {
	return func(thumbnail *discordgo.Thumbnail) {
		thumbnail.Spoiler = true
	}
}

// Thumbnail cria uma thumbnail a partir de uma URL.
func Thumbnail(
	url string,
	options ...ThumbnailOption,
) discordgo.Thumbnail {
	thumbnail := discordgo.Thumbnail{
		Media: discordgo.UnfurledMediaItem{
			URL: url,
		},
	}

	for _, option := range options {
		option(&thumbnail)
	}

	return thumbnail
}

type MediaOption func(*discordgo.MediaGalleryItem)

func MediaDescription(
	description string,
) MediaOption {
	return func(item *discordgo.MediaGalleryItem) {
		item.Description = &description
	}
}

func MediaSpoiler() MediaOption {
	return func(item *discordgo.MediaGalleryItem) {
		item.Spoiler = true
	}
}

func Media(
	url string,
	options ...MediaOption,
) discordgo.MediaGalleryItem {
	item := discordgo.MediaGalleryItem{
		Media: discordgo.UnfurledMediaItem{
			URL: url,
		},
	}

	for _, option := range options {
		option(&item)
	}

	return item
}

func Gallery(
	items ...discordgo.MediaGalleryItem,
) discordgo.MediaGallery {
	return discordgo.MediaGallery{
		Items: items,
	}
}

type FileComponentOption func(*discordgo.FileComponent)

func FileSpoiler() FileComponentOption {
	return func(file *discordgo.FileComponent) {
		file.Spoiler = true
	}
}

func FileComponent(
	filename string,
	options ...FileComponentOption,
) discordgo.FileComponent {
	file := discordgo.FileComponent{
		File: discordgo.UnfurledMediaItem{
			URL: "attachment://" + filename,
		},
	}

	for _, option := range options {
		option(&file)
	}

	return file
}
