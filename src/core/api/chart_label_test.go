// Copyright 2018 Project Harbor Authors
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

package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goharbor/harbor/src/chartserver"
	"github.com/goharbor/harbor/src/common"
	"github.com/goharbor/harbor/src/common/dao"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	resourceLabelAPIPath                = "/api/chartrepo/library/charts/harbor/0.2.0/labels"
	resourceLabelAPIPathWithFakeProject = "/api/chartrepo/not-exist/charts/harbor/0.2.0/labels"
	resourceLabelAPIPathWithFakeChart   = "/api/chartrepo/library/charts/not-exist/0.2.0/labels"
	cProLibraryLabelID                  int64
	mockChartServer                     *httptest.Server
	oldChartController                  *chartserver.Controller
)

func TestToStartMockChartService(t *testing.T) {
	var err error
	mockChartServer, oldChartController, err = mockChartController()
	if err != nil {
		t.Fatalf("failed to start the mock chart service: %v", err)
	}
}

func TestAddToChart(t *testing.T) {
	cSysLevelLabelID, err := dao.AddLabel(&models.Label{
		Name:  "c_sys_level_label",
		Level: common.LabelLevelSystem,
	})
	require.Nil(t, err)
	defer dao.DeleteLabel(cSysLevelLabelID)

	cProTestLabelID, err := dao.AddLabel(&models.Label{
		Name:      "c_pro_test_label",
		Level:     common.LabelLevelUser,
		Scope:     common.LabelScopeProject,
		ProjectID: 100,
	})
	require.Nil(t, err)
	defer dao.DeleteLabel(cProTestLabelID)

	cProLibraryLabelID, err = dao.AddLabel(&models.Label{
		Name:      "c_pro_library_label",
		Level:     common.LabelLevelUser,
		Scope:     common.LabelScopeProject,
		ProjectID: 1,
	})
	require.Nil(t, err)

	cases := []*codeCheckingCase{
		// 401
		{
			request: &testingRequest{
				url:    resourceLabelAPIPath,
				method: http.MethodPost,
			},
			code: http.StatusUnauthorized,
		},
		// 403
		{
			request: &testingRequest{
				url:        resourceLabelAPIPath,
				method:     http.MethodPost,
				credential: projGuest,
			},
			code: http.StatusForbidden,
		},
		// 500 project doesn't exist
		{
			request: &testingRequest{
				url:        resourceLabelAPIPathWithFakeProject,
				method:     http.MethodPost,
				credential: projDeveloper,
			},
			code: http.StatusNotFound,
		},
		// 404 chart doesn't exist
		{
			request: &testingRequest{
				url:        resourceLabelAPIPathWithFakeChart,
				method:     http.MethodPost,
				credential: projDeveloper,
			},
			code: http.StatusNotFound,
		},
		// 400
		{
			request: &testingRequest{
				url:        resourceLabelAPIPath,
				method:     http.MethodPost,
				credential: projDeveloper,
			},
			code: http.StatusBadRequest,
		},
		// 404 label doesn't exist
		{
			request: &testingRequest{
				url:        resourceLabelAPIPath,
				method:     http.MethodPost,
				credential: projDeveloper,
				bodyJSON: struct {
					ID int64
				}{
					ID: 1000,
				},
			},
			code: http.StatusNotFound,
		},
		// 400 system level label
		{
			request: &testingRequest{
				url:        resourceLabelAPIPath,
				method:     http.MethodPost,
				credential: projDeveloper,
				bodyJSON: struct {
					ID int64
				}{
					ID: cSysLevelLabelID,
				},
			},
			code: http.StatusBadRequest,
		},
		// 400 try to add the label of project1 to the image under project2
		{
			request: &testingRequest{
				url:        resourceLabelAPIPath,
				method:     http.MethodPost,
				credential: projDeveloper,
				bodyJSON: struct {
					ID int64
				}{
					ID: cProTestLabelID,
				},
			},
			code: http.StatusBadRequest,
		},
		// 200
		{
			request: &testingRequest{
				url:        resourceLabelAPIPath,
				method:     http.MethodPost,
				credential: projDeveloper,
				bodyJSON: struct {
					ID int64
				}{
					ID: cProLibraryLabelID,
				},
			},
			code: http.StatusOK,
		},
	}
	runCodeCheckingCases(t, cases...)
}

func TestGetOfChart(t *testing.T) {
	labels := []*models.Label{}
	err := handleAndParse(&testingRequest{
		url:        resourceLabelAPIPath,
		method:     http.MethodGet,
		credential: projDeveloper,
	}, &labels)
	require.Nil(t, err)
	require.Equal(t, 1, len(labels))
	assert.Equal(t, cProLibraryLabelID, labels[0].ID)
}

func TestRemoveFromChart(t *testing.T) {
	runCodeCheckingCases(t, &codeCheckingCase{
		request: &testingRequest{
			url:        fmt.Sprintf("%s/%d", resourceLabelAPIPath, cProLibraryLabelID),
			method:     http.MethodDelete,
			credential: projDeveloper,
		},
		code: http.StatusOK,
	})

	labels := []*models.Label{}
	err := handleAndParse(&testingRequest{
		url:        resourceLabelAPIPath,
		method:     http.MethodGet,
		credential: projDeveloper,
	}, &labels)
	require.Nil(t, err)
	require.Equal(t, 0, len(labels))
}

func TestToStopMockChartService(t *testing.T) {
	if mockChartServer != nil {
		mockChartServer.Close()
	}

	if oldChartController != nil {
		chartController = oldChartController
	}

	dao.DeleteLabel(cProLibraryLabelID)
}
