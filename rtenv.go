package miscutils

import (
	"fmt"
	"os"
	S "strings"
)

// RTEnv tries to identify the runtime environment.
// It uses lots of stuff from package [os]:
//   - Probably useful:
//   - Environ() []string
//   - Getenv(key string) string // LookupEnv distinguishes btwn empty and unset
//   - Getwd() (dir string, err error)
//   - UserHomeDir() (string, error)
//   - UserConfigDir() (string, error)
//   - Hostname() (name string, err error)
//   - Probably not useful:
//   - Executable() (string, error)
//   - UserCacheDir() (string, error)
//   - Getegid() int
//   - Geteuid() int
//   - Getgid() int
//   - Getuid() int
func RTEnv() string {
	var sb S.Builder
	hnm, _ := os.Hostname()
	cwd, _ := os.Getwd()
	usr, _ := os.UserHomeDir()
	cfg, _ := os.UserConfigDir()
	sb.WriteString(fmt.Sprintf("HN:%s U:%d eU:%d G:%d eG:%d \n",
		hnm, os.Getuid(), os.Geteuid(), os.Getgid(), os.Getegid()))
	sb.WriteString(fmt.Sprintf("cwd: %s\nusr: %s\ncfg: %s\n", cwd, usr, cfg))
	sb.WriteString(fmt.Sprintf("Env: %+v \n", os.Environ()))
	return sb.String()
}
