package opts

import (
	"runtime/debug"
	"testing"
)

func buildInfo(version string, settings ...string) *debug.BuildInfo {
	info := &debug.BuildInfo{}
	info.Main.Version = version
	for i := 0; i < len(settings); i += 2 {
		info.Settings = append(info.Settings, debug.BuildSetting{
			Key:   settings[i],
			Value: settings[i+1],
		})
	}
	return info
}

func TestBuildInfoVersion(t *testing.T) {
	rev := "e5cbf6f80183e83c1fa10fa17738d0981625b43a"
	//go install <pkg>@v1.2.3
	check(t, buildInfoVersion(buildInfo("v1.2.3")), "v1.2.3")
	//go1.24+ go build on a tagged commit
	check(t, buildInfoVersion(buildInfo("v1.2.3", "vcs.revision", rev, "vcs.modified", "false")), "v1.2.3")
	//...with uncommitted changes
	check(t, buildInfoVersion(buildInfo("v1.2.3+dirty", "vcs.revision", rev, "vcs.modified", "true")), "v1.2.3-src")
	//go1.24+ go build on an untagged commit, pseudo-version is skipped
	check(t, buildInfoVersion(buildInfo("v0.0.0-20260725231850-e5cbf6f80183", "vcs.revision", rev, "vcs.modified", "false")), "e5cbf6f")
	//...with uncommitted changes
	check(t, buildInfoVersion(buildInfo("v0.0.0-20260725231850-e5cbf6f80183+dirty", "vcs.revision", rev, "vcs.modified", "true")), "e5cbf6f-src")
	//go1.18-go1.23 go build, no module version at all
	check(t, buildInfoVersion(buildInfo("(devel)", "vcs.revision", rev, "vcs.modified", "false")), "e5cbf6f")
	//no version control information (go run, go test, -buildvcs=false, ...)
	check(t, buildInfoVersion(buildInfo("(devel)")), "")
	check(t, buildInfoVersion(buildInfo("")), "")
	//short revisions are used as-is
	check(t, buildInfoVersion(buildInfo("(devel)", "vcs.revision", "abc12")), "abc12")
}

func mockBuildInfo(t *testing.T, info *debug.BuildInfo, ok bool) {
	prev := readBuildInfo
	readBuildInfo = func() (*debug.BuildInfo, bool) { return info, ok }
	t.Cleanup(func() { readBuildInfo = prev })
}

func TestVersionDefaultsToBuildInfo(t *testing.T) {
	mockBuildInfo(t, buildInfo("(devel)", "vcs.revision", "e5cbf6f80183e83c1fa10fa17738d0981625b43a"), true)
	type Config struct {
		Foo string
	}
	o, _ := New(&Config{}).Name("myprog").ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, o.Help(), `
  Usage: myprog [options]

  Options:
  --foo, -f
  --version, -v  display version
  --help, -h     display help

  Version:
    e5cbf6f

`)
}

func TestVersionOverridesBuildInfo(t *testing.T) {
	mockBuildInfo(t, buildInfo("v9.9.9"), true)
	type Config struct {
		Foo string
	}
	o, _ := New(&Config{}).Name("myprog").Version("1.2.3").ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, o.Help(), `
  Usage: myprog [options]

  Options:
  --foo, -f
  --version, -v  display version
  --help, -h     display help

  Version:
    1.2.3

`)
}

func TestVersionMissingFromBuildInfo(t *testing.T) {
	mockBuildInfo(t, buildInfo("(devel)"), true)
	type Config struct {
		Foo string
	}
	o, _ := New(&Config{}).Name("myprog").ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, o.Help(), `
  Usage: myprog [options]

  Options:
  --foo, -f
  --help, -h  display help

`)
}
