// Copyright 2020 The Lokomotive Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build aws
// +build e2e

package aws

import (
	"testing"

	testingress "github.com/kinvolk/lokomotive/test/ingress"
)

func TestAWSIngress(t *testing.T) {
	testCases := []struct {
		Component string
		Namespace string
		Ingress   string
		Subpath   string
	}{
		{
			Component: "httpbin",
			Namespace: "httpbin",
			Ingress:   "httpbin",
			Subpath:   "get",
		},
		{
			Component: "prometheus-operator",
			Namespace: "monitoring",
			Ingress:   "prometheus-operator-kube-p-prometheus",
			Subpath:   "graph",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Ingress, func(t *testing.T) {
			t.Parallel()
			testingress.AccessIngress(t, tc.Namespace, tc.Ingress, tc.Subpath)
		})
	}
}
