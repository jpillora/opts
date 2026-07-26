package opts

import (
	"regexp"
	"runtime/debug"
	"strings"
)

//readBuildInfo is a variable to allow tests to swap it out
var readBuildInfo = debug.ReadBuildInfo

//go module pseudo-versions look like "v0.0.0-20260725231850-e5cbf6f80183",
//they are synthesised by the go tool for untagged commits
var pseudoVersionRe = regexp.MustCompile(`^v.+-[0-9]{14}-[0-9a-f]{12}(\+[0-9A-Za-z.-]+)?$`)

//length of the abbreviated VCS revision, matches "git log --oneline"
const shortRevisionLen = 7

//suffix marking a build which included uncommitted changes, as in
//"this commit, plus whatever source was found". Replaces the "+dirty"
//suffix used by the go tool.
const srcSuffix = "-src"

//buildVersion returns the version of the currently running program, as
//embedded by the go tool at compile time. Returns an empty string when
//the program was built without version control information (go run,
//go test, go build -buildvcs=false, builds from a source tarball, etc).
func buildVersion() string {
	info, ok := readBuildInfo()
	if !ok {
		return ""
	}
	return buildInfoVersion(info)
}

//buildInfoVersion extracts a display version from the given build info,
//preferring the VCS tag ("v1.2.3"), and falling back to the abbreviated
//VCS revision ("e5cbf6f").
func buildInfoVersion(info *debug.BuildInfo) string {
	//module version, set by "go install <pkg>@<version>" and, since go1.24,
	//also by "go build" when the commit is tagged. Untagged builds get a
	//pseudo-version, which is skipped in favour of the revision below.
	if v := info.Main.Version; v != "" && v != "(devel)" && !pseudoVersionRe.MatchString(v) {
		if strings.HasSuffix(v, "+dirty") {
			v = strings.TrimSuffix(v, "+dirty") + srcSuffix
		}
		return v
	}
	//no tag, use the commit hash instead
	revision := ""
	modified := false
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	if revision == "" {
		return ""
	}
	if len(revision) > shortRevisionLen {
		revision = revision[:shortRevisionLen]
	}
	if modified {
		revision += srcSuffix
	}
	return revision
}
