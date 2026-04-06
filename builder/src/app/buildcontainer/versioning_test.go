package buildcontainer

import "testing"

func TestLatestStableVersionTag(t *testing.T) {
	remoteOutput := "" +
		"aaa\trefs/tags/nightly\n" +
		"bbb\trefs/tags/v25.12.29\n" +
		"ccc\trefs/tags/v26.1.4\n" +
		"ddd\trefs/tags/v26.1.22\n" +
		"eee\trefs/tags/tip\n"

	if got := latestStableVersionTag(remoteOutput); got != "v26.1.22" {
		t.Fatalf("latestStableVersionTag() = %q, want %q", got, "v26.1.22")
	}
}

func TestArtifactDirectoryName(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		fallback string
		want     string
	}{
		{name: "strips leading v", version: "v1.3.1", fallback: "20260404-175649", want: "1.3.1"},
		{name: "keeps numeric version", version: "0.15.4", fallback: "20260402-224107", want: "0.15.4"},
		{name: "falls back for nightly", version: "nightly", fallback: "20260402-224420", want: "20260402-224420"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := artifactDirectoryName(test.version, test.fallback); got != test.want {
				t.Fatalf("artifactDirectoryName() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestParseBuildInfo(t *testing.T) {
	content := "package=ghostty\nversion=v1.3.1\ncommit=abc123\noutputs=ghostty\n"

	info, err := parseBuildInfo(content)
	if err != nil {
		t.Fatalf("parseBuildInfo() error = %v", err)
	}
	if info.Version != "v1.3.1" {
		t.Fatalf("parseBuildInfo() version = %q, want %q", info.Version, "v1.3.1")
	}
}
