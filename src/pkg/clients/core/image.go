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

package core

import (
	"fmt"

	"github.com/goharbor/harbor/src/common/models"
)

func (c *client) ListAllImages(project, repository string) ([]*models.TagResp, error) {
	url := c.buildURL(fmt.Sprintf("/api/repositories/%s/%s/tags", project, repository))
	var images []*models.TagResp
	if err := c.httpclient.GetAndIteratePagination(url, &images); err != nil {
		return nil, err
	}
	return images, nil
}

func (c *client) DeleteImage(project, repository, tag string) error {
	url := c.buildURL(fmt.Sprintf("/api/repositories/%s/%s/tags/%s", project, repository, tag))
	return c.httpclient.Delete(url)
}

func (c *client) DeleteImageRepository(project, repository string) error {
	url := c.buildURL(fmt.Sprintf("/api/repositories/%s/%s", project, repository))
	return c.httpclient.Delete(url)
}
