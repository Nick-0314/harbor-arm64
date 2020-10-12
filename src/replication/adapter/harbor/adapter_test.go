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
	"net/http"
	"strconv"
	"testing"

	"github.com/goharbor/harbor/src/common/utils/test"
	"github.com/goharbor/harbor/src/replication/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfo(t *testing.T) {
	// chart museum enabled
	server := test.NewServer(&test.RequestHandlerMapping{
		Method:  http.MethodGet,
		Pattern: "/api/systeminfo",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			data := `{"with_chartmuseum":true}`
			w.Write([]byte(data))
		},
	})
	registry := &model.Registry{
		URL: server.URL,
	}
	adapter, err := newAdapter(registry)
	require.Nil(t, err)
	info, err := adapter.Info()
	require.Nil(t, err)
	assert.Equal(t, model.RegistryTypeHarbor, info.Type)
	assert.Equal(t, 2, len(info.SupportedResourceFilters))
	assert.Equal(t, 2, len(info.SupportedTriggers))
	assert.Equal(t, 2, len(info.SupportedResourceTypes))
	assert.Equal(t, model.ResourceTypeImage, info.SupportedResourceTypes[0])
	assert.Equal(t, model.ResourceTypeChart, info.SupportedResourceTypes[1])
	server.Close()

	// chart museum disabled
	server = test.NewServer(&test.RequestHandlerMapping{
		Method:  http.MethodGet,
		Pattern: "/api/systeminfo",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			data := `{"with_chartmuseum":false}`
			w.Write([]byte(data))
		},
	})
	registry = &model.Registry{
		URL: server.URL,
	}
	adapter, err = newAdapter(registry)
	require.Nil(t, err)
	info, err = adapter.Info()
	require.Nil(t, err)
	assert.Equal(t, model.RegistryTypeHarbor, info.Type)
	assert.Equal(t, 2, len(info.SupportedResourceFilters))
	assert.Equal(t, 2, len(info.SupportedTriggers))
	assert.Equal(t, 1, len(info.SupportedResourceTypes))
	assert.Equal(t, model.ResourceTypeImage, info.SupportedResourceTypes[0])
	server.Close()
}

func TestPrepareForPush(t *testing.T) {
	server := test.NewServer(&test.RequestHandlerMapping{
		Method:  http.MethodPost,
		Pattern: "/api/projects",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		},
	})
	registry := &model.Registry{
		URL: server.URL,
	}
	adapter, err := newAdapter(registry)
	require.Nil(t, err)
	// nil resource
	err = adapter.PrepareForPush([]*model.Resource{nil})
	require.NotNil(t, err)
	// nil metadata
	err = adapter.PrepareForPush([]*model.Resource{
		{},
	})
	require.NotNil(t, err)
	// nil repository
	err = adapter.PrepareForPush(
		[]*model.Resource{
			{
				Metadata: &model.ResourceMetadata{},
			},
		})
	require.NotNil(t, err)
	// nil repository name
	err = adapter.PrepareForPush(
		[]*model.Resource{
			{
				Metadata: &model.ResourceMetadata{
					Repository: &model.Repository{},
				},
			},
		})
	require.NotNil(t, err)
	// project doesn't exist
	err = adapter.PrepareForPush(
		[]*model.Resource{
			{
				Metadata: &model.ResourceMetadata{
					Repository: &model.Repository{
						Name: "library/hello-world",
					},
				},
			},
		})
	require.Nil(t, err)

	server.Close()

	// project already exists
	server = test.NewServer(&test.RequestHandlerMapping{
		Method:  http.MethodPost,
		Pattern: "/api/projects",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
		},
	})
	registry = &model.Registry{
		URL: server.URL,
	}
	adapter, err = newAdapter(registry)
	require.Nil(t, err)
	err = adapter.PrepareForPush(
		[]*model.Resource{
			{
				Metadata: &model.ResourceMetadata{
					Repository: &model.Repository{
						Name: "library/hello-world",
					},
				},
			},
		})
	require.Nil(t, err)
}

func TestParsePublic(t *testing.T) {
	cases := []struct {
		metadata map[string]interface{}
		result   bool
	}{
		{nil, false},
		{map[string]interface{}{}, false},
		{map[string]interface{}{"public": true}, true},
		{map[string]interface{}{"public": "not_bool"}, false},
		{map[string]interface{}{"public": "true"}, true},
		{map[string]interface{}{"public": struct{}{}}, false},
	}
	for _, c := range cases {
		assert.Equal(t, c.result, parsePublic(c.metadata))
	}
}

func TestMergeMetadata(t *testing.T) {
	cases := []struct {
		m1     map[string]interface{}
		m2     map[string]interface{}
		public bool
	}{
		{
			m1: map[string]interface{}{
				"public": "true",
			},
			m2: map[string]interface{}{
				"public": "true",
			},
			public: true,
		},
		{
			m1: map[string]interface{}{
				"public": "false",
			},
			m2: map[string]interface{}{
				"public": "true",
			},
			public: false,
		},
		{
			m1: map[string]interface{}{
				"public": "false",
			},
			m2: map[string]interface{}{
				"public": "false",
			},
			public: false,
		},
	}
	for _, c := range cases {
		m := mergeMetadata(c.m1, c.m2)
		assert.Equal(t, strconv.FormatBool(c.public), m["public"].(string))
	}
}

func TestAbstractPublicMetadata(t *testing.T) {
	// nil input metadata
	meta := abstractPublicMetadata(nil)
	assert.Nil(t, meta)

	// contains no public metadata
	metadata := map[string]interface{}{
		"other": "test",
	}
	meta = abstractPublicMetadata(metadata)
	assert.Nil(t, meta)

	// contains public metadata
	metadata = map[string]interface{}{
		"other":  "test",
		"public": "true",
	}
	meta = abstractPublicMetadata(metadata)
	require.NotNil(t, meta)
	require.Equal(t, 1, len(meta))
	require.Equal(t, "true", meta["public"].(string))
}
