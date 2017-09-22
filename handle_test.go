// Copyright 2017 Canonical Ltd.
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

package raftmembership_test

import (
	"testing"
	"time"

	"github.com/CanonicalLtd/raft-membership"
	"github.com/CanonicalLtd/raft-test"
	"github.com/hashicorp/raft"
)

func TestHandleChangeRequests_ErrUnknownLeader(t *testing.T) {
	rafts, cleanup := rafttest.Cluster(t, rafttest.FSMs(3))
	defer cleanup()

	request := raftmembership.NewJoinRequest("1.2.3.4")
	handleOneChangeRequest(rafts[0], request)

	err := request.Error(time.Second)

	if err == nil {
		t.Fatal("no error returned despite no leader was elected yet")
	}
	switch err := err.(type) {
	case *raftmembership.ErrUnknownLeader:
	default:
		t.Fatalf("error is not of type ErrUnknownLeader: %v", err)
	}
}

func TestHandleChangeRequests_ErrDifferentLeader(t *testing.T) {
	notify := rafttest.Notify()
	rafts, cleanup := rafttest.Cluster(t, rafttest.FSMs(2), notify)
	defer cleanup()

	i := notify.NextAcquired(time.Second)
	j := rafttest.Other(rafts, i)

	raft := rafts[j]
	rafttest.WaitLeader(t, raft, time.Second)

	request := raftmembership.NewLeaveRequest("1.2.3.4")
	handleOneChangeRequest(raft, request)

	err := request.Error(time.Second)

	if err == nil {
		t.Fatal("no error returned despite request was made to non-leader")
	}
	switch err := err.(type) {
	case *raftmembership.ErrDifferentLeader:
		leader := raft.Leader()
		if err.Leader() != leader {
			t.Errorf("expected leader\n%q\ngot\n%q", leader, err.Leader())
		}
		break
	default:
		t.Fatalf("error is not of type ErrDifferentLeader: %v", err)
		break
	}
}

func TestHandleChangeRequests_KnownPeer(t *testing.T) {
	raft := rafttest.Node(t, rafttest.FSM())
	defer raft.Shutdown()

	request := raftmembership.NewJoinRequest("0")
	handleOneChangeRequest(raft, request)

	// The request is effectively a no-op and returns no error.
	if err := request.Error(time.Second); err != nil {
		t.Error(err)
	}
}

func TestHandleChangeRequests_LeaveRequest(t *testing.T) {
	raft := rafttest.Node(t, rafttest.FSM())
	defer raft.Shutdown()

	request := raftmembership.NewLeaveRequest("0")
	handleOneChangeRequest(raft, request)

	// The request succeeds.
	if err := request.Error(time.Second); err != nil {
		t.Error(err)
	}
}

// Wrapper around HandleChangeRequests that synchronously handles a
// single ChangeRequest.
func handleOneChangeRequest(raft *raft.Raft, request *raftmembership.ChangeRequest) {
	requests := make(chan *raftmembership.ChangeRequest, 1)
	requests <- request
	close(requests)
	raftmembership.HandleChangeRequests(raft, requests)
}
