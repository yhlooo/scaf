package version

import (
	"runtime"
	"runtime/debug"
)

// 构建时注入的版本信息
var (
	version = ""
)

// Version 版本信息
type Version struct {
	Version   string `json:"version" yaml:"version"`
	GitCommit string `json:"gitCommit" yaml:"gitCommit"`
	GoVersion string `json:"goVersion" yaml:"goVersion"`
	Arch      string `json:"arch" yaml:"arch"`
	OS        string `json:"os" yaml:"os"`
}

// GetVersion 获取版本信息
func GetVersion() Version {
	ret := Version{
		Version: version,
		Arch:    runtime.GOARCH,
		OS:      runtime.GOOS,
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		// 获取 go module 版本
		if ret.Version == "" {
			ret.Version = buildInfo.Main.Version
		}
		// 获取 Go 版本
		ret.GoVersion = buildInfo.GoVersion
		// 获取 Git 提交信息
		vcsRevision := ""
		vcsDirty := false
		for _, s := range buildInfo.Settings {
			switch s.Key {
			case "vcs.revision":
				vcsRevision = s.Value
			case "vcs.modified":
				vcsDirty = s.Value == "true"
			}
		}
		if vcsRevision != "" {
			ret.GitCommit = vcsRevision
			if vcsDirty {
				ret.GitCommit += "-dirty"
			}
		}
	}

	return ret
}
