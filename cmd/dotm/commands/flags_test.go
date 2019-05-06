package commands

import "testing"

func TestSanitizePath(t *testing.T) {
	for _, test := range []struct {
		path    string
		profile string
		want    string
	}{
		{
			path:    "$HOME/.config/dotm/profiles/myprofile/",
			profile: "",
			want:    "$HOME/.config/dotm/profiles/myprofile/",
		}, {
			path:    "$HOME/.config/dotm/profiles/<PROFILE>/",
			profile: "myprofile",
			want:    "$HOME/.config/dotm/profiles/myprofile/",
		},
	} {
		got := sanitizePath(test.path, test.profile)
		if test.want != got {
			t.Fatalf("sanitizePath(%s, %s). got %s. want %s", test.path, test.profile, got, test.want)
		}
	}
}
