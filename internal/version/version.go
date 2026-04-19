package version

import "encoding/json"

var (
	Version   = "dev"
	GitCommit = ""
	BuildTime = ""
)

type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
}

func GetInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
	}
}

func (i Info) JSON() []byte {
	b, _ := json.Marshal(i)
	return b
}
