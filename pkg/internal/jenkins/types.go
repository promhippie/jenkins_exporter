package jenkins

// Build defines the response from specific builds.
type Build struct {
	Timestamp int64 `json:"timestamp"`
	Duration  int64 `json:"duration"`
}

// Folder is a simple type used for folder listings.
type Folder struct {
	Class   string   `json:"_class"`
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Folders []Folder `json:"jobs"`
}

// Hudson defines the root type returned by the API.
type Hudson struct {
	Mode         string   `json:"mode"`
	NumExecutors int      `json:"numExecutors"`
	Folders      []Folder `json:"jobs"`
}

// BuildNumber defines a type for build numbers.
type BuildNumber struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
}

// Job defines the response from specific jobs.
type Job struct {
	Class                 string       `json:"_class"`
	Name                  string       `json:"displayName"`
	Path                  string       `json:"fullName"`
	URL                   string       `json:"url"`
	Disabled              bool         `json:"disabled"`
	Buildable             bool         `json:"buildable"`
	Color                 string       `json:"color"`
	LastBuild             *BuildNumber `json:"lastBuild"`
	LastCompletedBuild    *BuildNumber `json:"lastCompletedBuild"`
	LastFailedBuild       *BuildNumber `json:"lastFailedBuild"`
	LastStableBuild       *BuildNumber `json:"lastStableBuild"`
	LastSuccessfulBuild   *BuildNumber `json:"lastSuccessfulBuild"`
	LastUnstableBuild     *BuildNumber `json:"lastUnstableBuild"`
	LastUnsuccessfulBuild *BuildNumber `json:"lastUnsuccessfulBuild"`
	NextBuildNumber       int          `json:"nextBuildNumber"`
}
