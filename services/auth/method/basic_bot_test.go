// Copyright 2026 The Forgejo Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package method

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"forgejo.org/modules/setting"
)

func TestAllowBotTokenManagementBasic(t *testing.T) {
	previous := setting.Service.BasicAuthTokenUser
	setting.Service.BasicAuthTokenUser = "forgeclaw"
	t.Cleanup(func() { setting.Service.BasicAuthTokenUser = previous })

	tests := []struct {
		method, path, user string
		want               bool
	}{
		{http.MethodPost, "/api/v1/users/forgeclaw/tokens", "forgeclaw", true},
		{http.MethodDelete, "/api/v1/users/forgeclaw/tokens/123", "forgeclaw", true},
		{http.MethodGet, "/api/v1/users/forgeclaw/tokens", "forgeclaw", false},
		{http.MethodPost, "/api/v1/users/forgeclaw/tokens", "alice", false},
		{http.MethodPost, "/api/v1/users/alice/tokens", "forgeclaw", false},
		{http.MethodPost, "/api/v1/repos/example/repo", "forgeclaw", false},
		{http.MethodDelete, "/api/v1/users/forgeclaw/tokens", "forgeclaw", false},
	}
	for _, tc := range tests {
		request := httptest.NewRequest(tc.method, tc.path, nil)
		if got := allowBotTokenManagementBasic(request, tc.user); got != tc.want {
			t.Errorf("%s %s as %s: got %t, want %t", tc.method, tc.path, tc.user, got, tc.want)
		}
	}

	setting.Service.BasicAuthTokenUser = ""
	request := httptest.NewRequest(http.MethodPost, "/api/v1/users/forgeclaw/tokens", nil)
	if allowBotTokenManagementBasic(request, "forgeclaw") {
		t.Fatal("disabled bot exception allowed password auth")
	}
}
