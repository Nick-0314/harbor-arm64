// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package harbor

import (
	"fmt"
	"strings"

	"github.com/goharbor/harbor/src/common/utils"
	"github.com/goharbor/harbor/src/common/utils/log"
	adp "github.com/goharbor/harbor/src/replication/adapter"
	"github.com/goharbor/harbor/src/replication/model"
	"github.com/goharbor/harbor/src/replication/util"
)

func (a *adapter) FetchImages(filters []*model.Filter) ([]*model.Resource, error) {
	projects, err := a.listCandidateProjects(filters)
	if err != nil {
		return nil, err
	}

	resources := []*model.Resource{}
	for _, project := range projects {
		repositories, err := a.getRepositories(project.ID)
		if err != nil {
			return nil, err
		}
		if len(repositories) == 0 {
			continue
		}
		for _, filter := range filters {
			if err = filter.DoFilter(&repositories); err != nil {
				return nil, err
			}
		}

		var rawResources = make([]*model.Resource, len(repositories))
		runner := utils.NewLimitedConcurrentRunner(adp.MaxConcurrency)
		defer runner.Cancel()

		for i, r := range repositories {
			index := i
			repo := r
			runner.AddTask(func() error {
				vTags, err := a.getTags(repo.Name)
				if err != nil {
					return fmt.Errorf("List tags for repo '%s' error: %v", repo.Name, err)
				}
				if len(vTags) == 0 {
					rawResources[index] = nil
					return nil
				}
				for _, filter := range filters {
					if err = filter.DoFilter(&vTags); err != nil {
						return fmt.Errorf("Filter tags %v error: %v", vTags, err)
					}
				}
				if len(vTags) == 0 {
					rawResources[index] = nil
					return nil
				}
				tags := []string{}
				for _, vTag := range vTags {
					tags = append(tags, vTag.Name)
				}
				rawResources[index] = &model.Resource{
					Type:     model.ResourceTypeImage,
					Registry: a.registry,
					Metadata: &model.ResourceMetadata{
						Repository: &model.Repository{
							Name:     repo.Name,
							Metadata: project.Metadata,
						},
						Vtags: tags,
					},
				}

				return nil
			})
		}
		runner.Wait()

		if runner.IsCancelled() {
			return nil, fmt.Errorf("FetchImages error when collect tags for repos")
		}

		for _, r := range rawResources {
			if r != nil {
				resources = append(resources, r)
			}
		}
	}

	return resources, nil
}

func (a *adapter) listCandidateProjects(filters []*model.Filter) ([]*project, error) {
	pattern := ""
	for _, filter := range filters {
		if filter.Type == model.FilterTypeName {
			pattern = filter.Value.(string)
			break
		}
	}
	projects := []*project{}
	if len(pattern) > 0 {
		substrings := strings.Split(pattern, "/")
		projectPattern := substrings[0]
		names, ok := util.IsSpecificPathComponent(projectPattern)
		if ok {
			for _, name := range names {
				project, err := a.getProject(name)
				if err != nil {
					return nil, err
				}
				if project == nil {
					continue
				}
				projects = append(projects, project)
			}
		}
	}
	if len(projects) > 0 {
		names := []string{}
		for _, project := range projects {
			names = append(names, project.Name)
		}
		log.Debugf("parsed the projects %v from pattern %s", names, pattern)
		return projects, nil
	}
	return a.getProjects("")
}

// override the default implementation from the default image registry
// by calling Harbor API directly
func (a *adapter) DeleteManifest(repository, reference string) error {
	url := fmt.Sprintf("%s/api/repositories/%s/tags/%s", a.url, repository, reference)
	return a.client.Delete(url)
}

func (a *adapter) getTags(repository string) ([]*adp.VTag, error) {
	url := fmt.Sprintf("%s/api/repositories/%s/tags", a.getURL(), repository)
	tags := []*struct {
		Name   string `json:"name"`
		Labels []*struct {
			Name string `json:"name"`
		}
	}{}
	if err := a.client.Get(url, &tags); err != nil {
		return nil, err
	}
	vTags := []*adp.VTag{}
	for _, tag := range tags {
		var labels []string
		for _, label := range tag.Labels {
			labels = append(labels, label.Name)
		}
		vTags = append(vTags, &adp.VTag{
			Name:         tag.Name,
			Labels:       labels,
			ResourceType: string(model.ResourceTypeImage),
		})
	}
	return vTags, nil
}
