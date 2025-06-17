package jenkins

import (
	"context"
	"fmt"
)

// JobClient is a client for the jobs API.
type JobClient struct {
	client *Client
}

// Root returns a root API response.
func (c *JobClient) Root(ctx context.Context) (Hudson, error) {
	result := Hudson{}
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s/api/json?depth=1", c.client.endpoint), nil)

	if err != nil {
		return result, err
	}

	if _, err := c.client.Do(req, &result); err != nil {
		return result, err
	}

	return result, nil
}

// Build returns a specific build.
func (c *JobClient) Build(ctx context.Context, build *BuildNumber) (Build, error) {
	result := Build{}
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%sapi/json", build.URL), nil)

	if err != nil {
		return result, err
	}

	if _, err := c.client.Do(req, &result); err != nil {
		return result, err
	}

	return result, nil
}

// All returns all available jobs.
func (c *JobClient) All(ctx context.Context) ([]Job, error) {
	hudson, err := c.Root(ctx)

	if err != nil {
		return []Job{}, err
	}

	jobs, err := c.recursiveFolders(ctx, hudson.Folders)

	if err != nil {
		return []Job{}, err
	}

	return jobs, nil
}

func (c *JobClient) recursiveFolders(ctx context.Context, folders []Folder) ([]Job, error) {
	result := make([]Job, 0)

	for _, folder := range folders {
		switch class := folder.Class; class {
		case "com.cloudbees.hudson.plugins.folder.Folder":
			req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s/api/json?depth=1", folder.URL), nil)

			if err != nil {
				return nil, err
			}

			nextFolder := Folder{}

			if _, err := c.client.Do(req, &nextFolder); err != nil {
				return result, err
			}

			nextResult, err := c.recursiveFolders(ctx, nextFolder.Folders)

			if err != nil {
				return result, err
			}

			result = append(result, nextResult...)
		default:
			req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("%s/api/json", folder.URL), nil)

			if err != nil {
				return nil, err
			}

			job := Job{}

			if _, err := c.client.Do(req, &job); err != nil {
				return result, err
			}

			result = append(result, job)
		}
	}

	return result, nil
}
