package gitservice

type GSGitHubReleaseAuthor struct {
	login               string
	id                  int
	node_id             int
	avatar_url          string
	gravatar_id         string
	url                 string
	html_url            string
	followers_url       string
	following_url       string
	gists_url           string
	starred_url         string
	subscriptions_url   string
	organizations_url   string
	repos_url           string
	events_url          string
	received_events_url string
	// type string
	site_admin bool
}

type GSGitHubReleaseAsset struct {
	url     string
	id      int
	node_id string
	name    string
	// label null
	uploader             GSGitHubReleaseAuthor
	content_type         string
	state                string
	size                 int
	download_count       int
	created_at           string
	updated_at           string
	browser_download_url string
}

type GSGitHubRelease struct {
	url              string
	assets_url       string
	upload_url       string
	html_url         string
	id               int
	author           GSGitHubReleaseAuthor
	node_id          string
	tag_name         string
	target_commitish string
	name             string
	draft            bool
	prerelease       bool
	created_at       string
	published_at     string
	assets           []GSGitHubReleaseAsset
	tarball_url      string
	zipball_url      string
	body             string
}
