# Design: Git SHA in /version Endpoint

## Architecture

### Response Structure Enhancement
The `VersionInfo` struct will be extended to include git SHA:
```go
type VersionInfo struct {
    Version string `json:"version"`
    Module  string `json:"module"`
    GoOS    string `json:"go_os"`
    GoArch  string `json:"go_arch"`
    GitSha  string `json:"git_sha"`  // NEW
}
```

### Git SHA Retrieval
Git SHA will be extracted during build/runtime from available sources (in priority order):
1. **Go build info** (via `runtime/debug.ReadBuildInfo()`)
   - Modern Go (1.12+) provides vcs metadata
   - Includes commit hash when built from git repository
2. **Fallback** if not available:
   - "unknown" or empty string
   - Operators can provide via `-ldflags` if needed

### Endpoint Handler
The existing `handleVersion()` method in `server/management.go` will be enhanced:
- Extract git SHA using `runtime/debug.ReadBuildInfo()`
- Include `GitSha` field in response
- Maintain backward compatibility with existing fields
- No changes to HTTP method validation (GET only)

### Response Format Example
```json
{
  "version": "v1.0.0",
  "module": "github.com/phona/ubox-crosser",
  "go_os": "linux",
  "go_arch": "amd64",
  "git_sha": "abc123def456..."
}
```

## Implementation Details

### Changes to `server/management.go`
1. Add `GitSha string` field to `VersionInfo` struct
2. Enhance `handleVersion()` to populate git SHA
3. Extract VCS revision from `runtime/debug.ReadBuildInfo()`
4. Handle cases where git metadata is not available

### Git SHA Extraction Logic
```go
// Pseudo-code for git SHA extraction
func getGitSha() string {
    if buildInfo, ok := debug.ReadBuildInfo(); ok {
        for _, setting := range buildInfo.Settings {
            if setting.Key == "vcs.revision" {
                return setting.Value  // Returns full or short SHA
            }
        }
    }
    return "unknown"  // Fallback
}
```

### Build Considerations
- Go 1.18+ automatically includes VCS revision when building from git
- No special build flags required for git metadata
- Works with `go build`, `go install`, and standard CI/CD
- Git metadata requires `.git` directory to be present at build time

## Contract Specification
The contract spec (`specs/version-endpoint-git-sha/spec.md`) defines:
- Expected response structure with all 5 fields
- Git SHA format validation (hexadecimal, 7-40 characters)
- Consistency across requests
- HTTP method validation (GET only, 405 for others)
- Edge cases (development builds, missing git metadata)

## Testing Strategy
- **Contract tests** (`tests/contract/version_test.go`):
  - Validate git_sha field presence
  - Verify hexadecimal format
  - Check consistency across requests
  - Verify all existing fields still present
  - Test HTTP method restrictions
  - Measure response time

- **Unit tests** (in `server/management_test.go`):
  - Test git SHA extraction logic
  - Test fallback behavior when git metadata unavailable
  - Test JSON marshaling of VersionInfo

## Backward Compatibility
- Fully backward compatible
- Existing fields (version, module, go_os, go_arch) unchanged in meaning
- Only adds new `git_sha` field
- Clients expecting 4 fields will get 5; well-designed clients should handle this
- No changes to existing endpoints (`/health`, `/healthz`)

## Edge Cases Handled
1. **No git metadata**: `git_sha` contains "unknown"
2. **Development build**: git metadata still available if in git repo
3. **Shallow clone**: git SHA still available (commit hash is included)
4. **Build outside repo**: `git_sha` will be "unknown"
