package types

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

// BuildNumber defines a type for build numbers.
type BuildNumber struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
}
