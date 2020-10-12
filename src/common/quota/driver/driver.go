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

package driver

import (
	"sync"

	"github.com/goharbor/harbor/src/pkg/types"
)

var (
	driversMu sync.RWMutex
	drivers   = map[string]Driver{}
)

// RefObject type for quota ref object
type RefObject map[string]interface{}

// Driver the driver for quota
type Driver interface {
	// HardLimits returns default resource list
	HardLimits() types.ResourceList
	// Load returns quota ref object by key
	Load(key string) (RefObject, error)
	// Validate validate the hard limits
	Validate(hardLimits types.ResourceList) error
}

// Register register quota driver
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("quota: Register driver is nil")
	}

	drivers[name] = driver
}

// Get returns quota driver by name
func Get(name string) (Driver, bool) {
	driversMu.Lock()
	defer driversMu.Unlock()

	driver, ok := drivers[name]
	return driver, ok
}
