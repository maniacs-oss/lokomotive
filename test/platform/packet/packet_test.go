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

// +build packet_fluo
// +build poste2e

package packet_test

import (
	"context"
	"os"
	"testing"

	"github.com/packethost/packngo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	testutil "github.com/kinvolk/lokomotive/test/components/util"
)

func TestBGPDisabled(t *testing.T) {
	client := testutil.CreateKubeClient(t)

	cl, err := packngo.NewClient()
	if err != nil {
		t.Fatalf("Creating Packet API client: %v", err)
	}

	projectID := os.Getenv("PACKET_PROJECT_ID")
	if projectID == "" {
		t.Fatalf("Packet project ID can't be empty. Is %q environment variable set?", "PACKET_PROJECT_ID")
	}

	nodeLabel := "lokomotive.alpha.kinvolk.io/bgp-enabled=false"

	// Select a node from the general worker pool.
	nodesList, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: nodeLabel,
	})
	if err != nil {
		t.Fatalf("Listing nodes with label %q: %v", nodeLabel, err)
	}

	nodes := nodesList.Items
	if len(nodes) < 1 {
		t.Fatalf("Wanted one or more nodes with label %q, found none.", nodeLabel)
	}

	hostname := nodes[0].Name

	devices, _, err := cl.Devices.List(projectID, nil)
	if err != nil {
		t.Fatalf("Listing devices in project %q: %v", projectID, err)
	}

	deviceID := ""

	for _, device := range devices {
		if device.Hostname == hostname {
			deviceID = device.ID

			break
		}
	}

	if deviceID == "" {
		t.Fatalf("No Packet device found with hostname %q", hostname)
	}

	sessions, _, err := cl.Devices.ListBGPSessions(deviceID, nil)
	if err != nil {
		t.Fatalf("Getting BGP sessions for device %q: %v", deviceID, err)
	}

	if len(sessions) != 0 {
		t.Fatalf("Worker pool with BGP disabled should not have any BGP sessions")
	}
}
