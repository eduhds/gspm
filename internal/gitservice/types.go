package gitservice

type GSRelease struct {
	Id      int64
	TagName string
	Assets  []GSReleaseAsset
}

type GSReleaseAsset struct {
	Url  string
	Id   int64
	Name string
}

type ErrorMessage struct {
	Message string `json:"message"`
}

//
// GITHUB
//

type GSGitHubReleaseAuthor struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	Site_admin        bool   `json:"site_admin"`
}

type GSGitHubReleaseAsset struct {
	Url                string                `json:"url"`
	Id                 int64                 `json:"id"`
	NodeId             string                `json:"node_id"`
	Name               string                `json:"name"`
	Label              string                `json:"label"`
	Uploader           GSGitHubReleaseAuthor `json:"uploader"`
	ContentType        string                `json:"content_type"`
	State              string                `json:"state"`
	Size               int                   `json:"size"`
	DownloadCount      int                   `json:"download_count"`
	CreatedAt          string                `json:"created_at"`
	UpdatedAt          string                `json:"updated_at"`
	BrowserDownloadUrl string                `json:"browser_download_url"`
}

type GSGitHubRelease struct {
	Url             string                 `json:"url"`
	AssetsUrl       string                 `json:"assets_url"`
	UploadUrl       string                 `json:"upload_url"`
	HtmlUrl         string                 `json:"html_url"`
	Id              int64                  `json:"id"`
	Author          GSGitHubReleaseAuthor  `json:"author"`
	NodeId          string                 `json:"node_id"`
	TagName         string                 `json:"tag_name"`
	TargetCommitish string                 `json:"target_commitish"`
	Name            string                 `json:"name"`
	Draft           bool                   `json:"draft"`
	Prerelease      bool                   `json:"prerelease"`
	CreatedAt       string                 `json:"created_at"`
	PublishedAt     string                 `json:"published_at"`
	Assets          []GSGitHubReleaseAsset `json:"assets"`
	TarballUrl      string                 `json:"tarball_url"`
	ZipballUrl      string                 `json:"zipball_url"`
	Body            string                 `json:"body"`
}

//
// BITBUCKET
//

type BitbucketResponse struct {
	Values  []BitbucketDownload `json:"values"`
	Pagelen int                 `json:"pagelen"`
	Size    int                 `json:"size"`
	Page    int                 `json:"page"`
}

type BitbucketDownload struct {
	Name      string         `json:"name"`
	Size      int64          `json:"size"`
	CreatedOn string         `json:"created_on"`
	User      BitbucketUser  `json:"user"`
	Downloads int            `json:"downloads"`
	Links     BitbucketLinks `json:"links"`
	Type      string         `json:"type"`
}

type BitbucketUser struct {
	DisplayName string         `json:"display_name"`
	Links       BitbucketLinks `json:"links"`
	Type        string         `json:"type"`
	Uuid        string         `json:"uuid"`
	AccountId   string         `json:"account_id"`
	Nickname    string         `json:"nickname"`
}

type BitbucketLinks struct {
	Self   BitbucketLink `json:"self"`
	Avatar BitbucketLink `json:"avatar"`
	Html   BitbucketLink `json:"html"`
}

type BitbucketLink struct {
	Href string `json:"href"`
}
