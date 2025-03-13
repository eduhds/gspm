package gitservice

type GSReleaseAsset struct {
	Url  string
	Id   int64
	Name string
}

type GSRelease struct {
	TagName string
	Assets  []GSReleaseAsset
}
