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

package listrepo

import (
	"encoding/json"
	"github.com/goharbor/harbor/src/common/dao"
	"github.com/goharbor/harbor/src/common/utils/log"
	"github.com/goharbor/harbor/src/core/middlewares/util"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
)

const (
	catalogURLPattern = `/v2/_catalog`
)

type listReposHandler struct {
	next http.Handler
}

// New ...
func New(next http.Handler) http.Handler {
	return &listReposHandler{
		next: next,
	}
}

// ServeHTTP ...
func (lrh listReposHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var rec *httptest.ResponseRecorder
	listReposFlag := matchListRepos(req)
	if listReposFlag {
		rec = httptest.NewRecorder()
		lrh.next.ServeHTTP(rec, req)
		if rec.Result().StatusCode != http.StatusOK {
			util.CopyResp(rec, rw)
			return
		}
		var ctlg struct {
			Repositories []string `json:"repositories"`
		}
		decoder := json.NewDecoder(rec.Body)
		if err := decoder.Decode(&ctlg); err != nil {
			log.Errorf("Decode repositories error: %v", err)
			util.CopyResp(rec, rw)
			return
		}
		var entries []string
		for repo := range ctlg.Repositories {
			log.Debugf("the repo in the response %s", ctlg.Repositories[repo])
			exist := dao.RepositoryExists(ctlg.Repositories[repo])
			if exist {
				entries = append(entries, ctlg.Repositories[repo])
			}
		}
		type Repos struct {
			Repositories []string `json:"repositories"`
		}
		resp := &Repos{Repositories: entries}
		respJSON, err := json.Marshal(resp)
		if err != nil {
			log.Errorf("Encode repositories error: %v", err)
			util.CopyResp(rec, rw)
			return
		}

		for k, v := range rec.Header() {
			rw.Header()[k] = v
		}
		clen := len(respJSON)
		rw.Header().Set(http.CanonicalHeaderKey("Content-Length"), strconv.Itoa(clen))
		rw.Write(respJSON)
		return
	}
	lrh.next.ServeHTTP(rw, req)
}

// matchListRepos checks if the request looks like a request to list repositories.
func matchListRepos(req *http.Request) bool {
	if req.Method != http.MethodGet {
		return false
	}
	re := regexp.MustCompile(catalogURLPattern)
	s := re.FindStringSubmatch(req.URL.Path)
	if len(s) == 1 {
		return true
	}
	return false
}
