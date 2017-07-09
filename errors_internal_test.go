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
	"testing"
)

func TestErrDifferentLeader_Error(t *testing.T) {
	err := &ErrDifferentLeader{leader: "1.2.3.4"}
	const want = "node is not leader, current leader at: 1.2.3.4"
	got := err.Error()
	if got != want {
		t.Errorf("expected error message\n%q\ngot\n%q", want, got)
	}
}

func TestErrDifferentLeader_Leader(t *testing.T) {
	err := &ErrDifferentLeader{leader: "1.2.3.4"}
	const want = "1.2.3.4"
	got := err.Leader()
	if got != want {
		t.Errorf("expected leader address\n%q\ngot\n%q", want, got)
	}
}

func TestErrUnknownLeader_Error(t *testing.T) {
	err := &ErrUnknownLeader{}
	const want = "node is not leader, current leader unknown"
	got := err.Error()
	if got != want {
		t.Errorf("expected error message\n%q\ngot\n%q", want, got)
	}
}
