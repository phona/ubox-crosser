# REQ-fresh-1776858473: Add Git SHA to /version Endpoint

## Summary
Enhance the existing `/version` endpoint to include the git commit SHA of the build. This allows operators and monitoring systems to identify the exact source code version deployed, enabling better tracking of which features and fixes are in a running instance.

## Problem Statement
The current `/version` endpoint returns:
- `version`: Build version (from go build info)
- `module`: Module path
- `go_os`: Target OS
- `go_arch`: Target architecture

However, it does not include the git commit SHA, which means:
- Operators cannot verify the exact commit deployed
- Monitoring systems cannot correlate running code with git history
- Debugging issues requires manual git log correlation
- CI/CD pipelines cannot reliably track which commit is running in production

## Solution Overview
Add a `git_sha` field to the `/version` endpoint response that contains the git commit SHA. This field will be:
1. Populated from git build metadata during compilation
2. Included in every `/version` response
3. Consistent across requests (static build info)
4. Validated by contract tests to ensure proper format

## Acceptance Criteria
- [ ] `/version` GET response includes `git_sha` field
- [ ] `git_sha` contains valid git commit hash (hexadecimal)
- [ ] `git_sha` is at least 7 characters (short SHA format)
- [ ] `git_sha` is consistent across multiple requests
- [ ] `git_sha` format is validated by contract tests
- [ ] Existing fields (version, module, go_os, go_arch) remain unchanged
- [ ] Response is still valid JSON with Content-Type `application/json`

## Stage Scope
**This is a contract-spec stage task.** The scope is to define and validate the contract for how the `/version` endpoint should behave when git SHA is included. Implementation in the dev stage will follow these contract specifications.

## Target Branch
`stage/REQ-fresh-1776858473-dev`
