# File Block plugin

This plugin allows Mattermost administrators to block file uploads with specific extensions.

## Plugin Parameters

* Allowed extensions - Comma separated list of allowed (whitelisted) extensions. e.g. "doc,docx,png,pdf"
* Require extension - true/false value to allow/reject files without extension
* [Experimental] Validate Mime Content - true/false value to match the MIME extension extracted from the file if the file extension matches. This aims to avoid renaming a forbidden file extension to upload unauthorized files
