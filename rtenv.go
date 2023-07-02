package miscutils

import (
	"fmt"
	"golang.org/x/exp/slices"
	"os"
	RT "runtime"
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
	sb.WriteString(fmt.Sprintf(
		"Hostname: %s \n\t GO-OS: %s \n\t GO-ARCH: %s \n\t GO-ROOT: %s \n",
		hnm, RT.GOOS, RT.GOARCH, RT.GOROOT()))
	sb.WriteString(fmt.Sprintf("\t U:%d G:%d \n", os.Getuid(), os.Getgid()))
	sb.WriteString(fmt.Sprintf("cwd: %s\nusr: %s\ncfg: %s\n", cwd, usr, cfg))
	sb.WriteString("Env:\n")
	env := os.Environ()
	slices.Sort(env)
	for _, s := range env {
		if !S.HasPrefix(s, "BASH_FUNC") {
			sb.WriteString(s + "\n")
		}
	}
	return sb.String()
}
