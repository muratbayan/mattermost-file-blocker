package main

import (
	"io"
	"strings"
	"sync"

	"github.com/gabriel-vasile/mimetype"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

// FileBlockerPlugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type FileBlockerPlugin struct {
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
func (p *FileBlockerPlugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {
	config := p.getConfiguration()
	p.API.LogInfo("session Id from context", "sessionId", c.SessionId)
	p.API.LogInfo("User Agent from context", "userAgent", c.UserAgent)
	p.API.LogInfo("RequestId from context", "requestId", c.RequestId)
	p.API.LogInfo("Plugin context", "pluginContext", c)

	session, sessionErr := p.API.GetSession(c.SessionId)

	if sessionErr != nil {
		p.API.LogError("Session retrieval error", "error", sessionErr.Error(), "detailedError", sessionErr.DetailedError)
		return nil, "File Blocker plugin - There was an error retrieving the session information"
	}

	if session.IsMobileApp() {
		p.API.LogInfo("The session is a mobile session")
	} else {
		p.API.LogInfo("The session is not a mobile session")
	}

	extensions := strings.Split(config.AllowedExtensions, ",")

	if config.ExtensionIsRequired && info.Extension == "" {
		p.API.LogInfo("File attachments without extensions are not allowed", "filename", info.Name, "user", info.CreatorId, "extension", info.Extension)
		return nil, "File Blocker plugin - File attachments without extensions are not allowed"
	}

	found := stringSliceContains(extensions, info.Extension)

	if !found {
		p.API.LogInfo("Unsupported file attachment extension", "filename", info.Name, "user", info.CreatorId, "extension", info.Extension, "allowedExtensions", strings.Join(extensions, ", "))
		return nil, "File Blocker plugin - This file attachment extension is not allowed"
	}

	if config.CheckMimeType {
		mimeTypeResult, mimeErr := mimetype.DetectReader(file)

		if mimeErr != nil {
			p.API.LogError("MIME Type detection error", "filename", info.Name, "user", info.CreatorId)
			return nil, "File Blocker plugin - An error occurred during the verification of the file attachment - Please contact your administrator"
		}

		p.API.LogDebug("MIME Output", "mimeTypeResult", mimeTypeResult.String())
		p.API.LogDebug("MIME Extension", "mimeTypeResult", mimeTypeResult.Extension())

		mimeExtension := strings.Trim(mimeTypeResult.Extension(), ".")
		mimeFound := stringSliceContains(extensions, mimeExtension)

		// Should we simply fail whenever the extension does not match the mime extension?
		if !mimeFound {
			return nil, "File Blocker plugin - Extension does not match MIME extension and this MIME extension is not whitelisted"
		}
	}

	return info, ""
}
