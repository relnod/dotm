package dotm

import "testing"

func TestSanitizeRemote(t *testing.T) {
	for _, test := range []struct {
		remote string
		want   string
	}{
		{
			remote: "https://github.com/user/repo.git",
			want:   "https://github.com/user/repo.git",
		}, {
			remote: "git@github.com:user/repo.git",
			want:   "git@github.com:user/repo.git",
		}, {
			remote: "github.com/user/repo",
			want:   "https://github.com/user/repo.git",
		},
	} {
		got := sanitizeRemote(test.remote)
		if test.want != got {
			t.Fatalf("sanitizeRemote(%s). got %s. want %s", test.remote, got, test.want)
		}
	}
}
