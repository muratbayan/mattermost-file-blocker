package main

import (
	"io"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// FileBlockPlugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type FileBlockPlugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func stringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// FileWillBeUploaded matches file attachments against known extensions from configuration
func (p *FileBlockPlugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {
	config := p.getConfiguration()
	// message := "This file could not be attached to your post"

	extensions := strings.Split(config.ForbiddenExtensions, ",")

	if config.ExtensionIsRequired && info.Extension == "" {
		p.API.LogWarn("File attachments without extensions are not allowed", "filename", info.Name, "user", info.CreatorId, "extension", info.Extension)
		return nil, "The File Block plugin did not allow you to attach this file."
	}

	found := stringSliceContains(extensions, info.Extension)

	if found {
		postInfo, _ := p.API.GetPost(info.PostId)

		if postInfo != nil {
			p.API.LogWarn("ChannelID", "channelId", postInfo.ChannelId)
		}

		p.API.LogWarn("Unsupported file attachment extension", "filename", info.Name, "user", info.CreatorId, "extension", info.Extension)
		return nil, "The File Block plugin did not allow you to attach this file."
	}

	return info, ""
}
