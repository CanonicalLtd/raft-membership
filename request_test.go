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
)

func TestChangeRequestKind_String(t *testing.T) {
	if s := raftmembership.JoinRequest.String(); s != "join" {
		t.Errorf("unexpected string value for join kind: %s", s)
	}
	if s := raftmembership.LeaveRequest.String(); s != "leave" {
		t.Errorf("unexpected string value for leave kind: %s", s)
	}
}

func TestChangeRequest_ErrorTimeout(t *testing.T) {
	request := &raftmembership.ChangeRequest{}
	err := request.Error(time.Microsecond)
	if err == nil {
		t.Fatal("request returned no error")
	}
	if err.Error() != "timeout waiting for membership change" {
		t.Errorf("unexpected error %s", err.Error())
	}
}
