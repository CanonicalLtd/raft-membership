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

package raftmembership

import (
	"fmt"
	"time"

	"github.com/hashicorp/raft"
)

// HandleChangeRequests processes ChangeRequest's received through the
// given channel, using the given raft.Raft instance to add or remove
// peers to the cluster according to the received requests.
func HandleChangeRequests(r *raft.Raft, requests <-chan *ChangeRequest) {
	for request := range requests {

		// If we currently think we're the leader, let's try
		// to handle the request, otherwise let's bail out
		// directly.
		var err error
		if r.State() == raft.Leader {
			raftMethod := raftMethodForRequest(request)
			err = raftMethod(r, request.Peer()).Error()
		} else {
			err = raft.ErrNotLeader
		}

		// Wrap not-leader errors.
		if err == raft.ErrNotLeader {
			if r.Leader() != "" {
				err = &ErrDifferentLeader{leader: string(r.Leader())}
			} else {
				err = &ErrUnknownLeader{}
			}
		}

		// Ignore errors due to a joining peer being already
		// known.
		//if err == raft.ErrKnownPeer {
		//err = nil
		//}

		request.Done(err)
	}
}

// Return the appropriate Raft method to invoke for handling the given
// request.
func raftMethodForRequest(request *ChangeRequest) func(*raft.Raft, string) raft.Future {
	kind := request.Kind()
	switch kind {
	case JoinRequest:
		return func(r *raft.Raft, addr string) raft.Future {
			return r.AddVoter(
				raft.ServerID(addr), raft.ServerAddress(addr), 0, 10*time.Second)
		}
	case LeaveRequest:
		return func(r *raft.Raft, addr string) raft.Future {
			return r.RemoveServer(
				raft.ServerID(addr), 0, 10*time.Second)
		}
	default:
		panic(fmt.Sprintf("invalid change request kind: %d", int(kind)))
	}
}
