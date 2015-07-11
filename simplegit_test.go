package simplegit

import (
	"testing"
)

func Test_SimpleGit_detectSimpleGitUser(t *testing.T) {
	s := SimpleGit{}
	user := s.detectGitUser()
	if user == "" {
		t.Fatalf("User empty")
	}
	t.Logf("%s", user)
}
