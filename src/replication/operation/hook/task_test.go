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

package hook

import (
	"testing"

	"github.com/goharbor/harbor/src/jobservice/job"
	"github.com/goharbor/harbor/src/replication/dao/models"
	"github.com/goharbor/harbor/src/replication/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakedOperationController struct {
	status string
}

func (f *fakedOperationController) StartReplication(*model.Policy, *model.Resource, model.TriggerType) (int64, error) {
	return 0, nil
}
func (f *fakedOperationController) StopReplication(int64) error {
	return nil
}
func (f *fakedOperationController) ListExecutions(...*models.ExecutionQuery) (int64, []*models.Execution, error) {
	return 0, nil, nil
}
func (f *fakedOperationController) GetExecution(int64) (*models.Execution, error) {
	return nil, nil
}
func (f *fakedOperationController) ListTasks(...*models.TaskQuery) (int64, []*models.Task, error) {
	return 0, nil, nil
}
func (f *fakedOperationController) GetTask(int64) (*models.Task, error) {
	return nil, nil
}
func (f *fakedOperationController) UpdateTaskStatus(id int64, status string, statusRevision int64, statusCondition ...string) error {
	f.status = status
	return nil
}
func (f *fakedOperationController) GetTaskLog(int64) ([]byte, error) {
	return nil, nil
}

func TestUpdateTask(t *testing.T) {
	mgr := &fakedOperationController{}
	cases := []struct {
		inputStatus    string
		expectedStatus string
	}{
		{
			inputStatus:    job.PendingStatus.String(),
			expectedStatus: models.TaskStatusPending,
		},
		{
			inputStatus:    job.ScheduledStatus.String(),
			expectedStatus: models.TaskStatusInProgress,
		},
		{
			inputStatus:    job.RunningStatus.String(),
			expectedStatus: models.TaskStatusInProgress,
		},
		{
			inputStatus:    job.StoppedStatus.String(),
			expectedStatus: models.TaskStatusStopped,
		},
		{
			inputStatus:    job.ErrorStatus.String(),
			expectedStatus: models.TaskStatusFailed,
		},
		{
			inputStatus:    job.SuccessStatus.String(),
			expectedStatus: models.TaskStatusSucceed,
		},
	}

	for _, c := range cases {
		err := UpdateTask(mgr, 1, c.inputStatus, 1)
		require.Nil(t, err)
		assert.Equal(t, c.expectedStatus, mgr.status)
	}
}
