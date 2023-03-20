package types

import (
	"encoding/json"
)

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

// BuildNumber defines a custom type for build numbers.
type BuildNumber int

// UnmarshalJSON implements the json unmarshal interface.
func (b *BuildNumber) UnmarshalJSON(data []byte) error {
	type details struct {
		Number int `json:"number"`
	}

	var (
		v = details{}
	)

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.Number != 0 {
		*b = BuildNumber(v.Number)
	}

	return nil
}
