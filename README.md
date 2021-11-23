# File Blocker plugin

Contact: @murat.bayan on the [Mattermost Community Server](https://community.mattermost.com/)

This plugin allows Mattermost administrators to block file uploads with specific extensions.

![Screenshot 2020-04-06 at 23 01 56 copy](https://user-images.githubusercontent.com/6051508/78684707-b5704100-78e8-11ea-80d1-5fabd42ff036.png)

## Plugin Parameters

* Allowed extensions - Comma separated list of allowed (whitelisted) extensions. e.g. "doc,docx,png,pdf"
* Require extension - true/false value to allow/reject files without extension
* [Experimental] Validate Mime Content - true/false value to match the MIME extension extracted from the file if the file extension matches. This aims to avoid renaming a forbidden file extension to upload unauthorized files
* Allow guest users to attach files - true/false value to allow/prevent guest accounts to post file attachments (defaults to true)
* Allow Mobile Attachments - true/false value to allow/prevent file attachments from the Mattermost mobile application (defaults to true)