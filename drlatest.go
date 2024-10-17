/*

GoFmt
GoBuildNull
GoBuild

go get -u -a -v
go mod tidy
*/

package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/rusenask/docker-registry-client/registry"
)

const NL = "\n"

func log(msg string, args ...interface{}) {
	const NL = "\n"
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, msg+NL)
	} else {
		fmt.Fprintf(os.Stderr, msg+NL, args...)
	}
}

var (
	RegistryUsername   string
	RegistryPassword   string
	RegistryUrl        string
	RegistryHost       string
	RegistryRepository string
)

func init() {
	RegistryUsername = os.Getenv("RegistryUsername")
	if RegistryUsername == "" {
		log("warning: RegistryUsername env var empty")
	}
	RegistryPassword = os.Getenv("RegistryPassword")
	if RegistryPassword == "" {
		log("warning: RegistryPassword env var empty")
	}
}

type Versions []string

func (vv Versions) Len() int {
	return len(vv)
}

func (vv Versions) Less(i, j int) bool {
	v1, v2 := vv[i], vv[j]
	v1s := strings.Split(v1, ".")
	v2s := strings.Split(v2, ".")
	if len(v1s) < len(v2s) {
		return true
	} else if len(v1s) > len(v2s) {
		return false
	}
	for e := 0; e < len(v1s); e++ {
		d1, _ := strconv.Atoi(v1s[e])
		d2, _ := strconv.Atoi(v2s[e])
		if d1 < d2 {
			return true
		} else if d1 > d2 {
			return false
		}
	}
	return false
}

func (vv Versions) Swap(i, j int) {
	vv[i], vv[j] = vv[j], vv[i]
}

func main() {
	all := flag.Bool("all", false, "to print all tags, otherwise only the last tag is printed")
	full := flag.Bool("full", false, "to show full image address like registry/path:tag, otherwise only tag is printed")
	flag.Parse()
	if flag.NArg() < 1 {
		log("usage: drlatest docker.registry.repository.url")
		os.Exit(1)
	}

	if u, err := url.Parse(flag.Args()[0]); err != nil {
		log("error: docker.registry.repository.url parse: %v", err)
		os.Exit(1)
	} else {
		RegistryUrl = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
		RegistryHost = u.Host
		RegistryRepository = u.Path
	}
	//log("registry:%s repository:%s", RegistryUrl, RegistryRepository)

	r := registry.NewInsecure(RegistryUrl, RegistryUsername, RegistryPassword)
	r.Logf = registry.Quiet

	tags, err := r.Tags(RegistryRepository)
	if err != nil {
		log("error: list tags: %v", err)
		os.Exit(1)
	}

	sort.Sort(Versions(tags))

	if *all {
		for _, tag := range tags {
			if *full {
				fmt.Printf("%s%s:%s"+NL, RegistryHost, RegistryRepository, tag)
			} else {
				fmt.Printf("%s"+NL, tag)
			}
		}
	} else if len(tags) > 0 {
		if *full {
			fmt.Printf("%s%s:%s"+NL, RegistryHost, RegistryRepository, tags[len(tags)-1])
		} else {
			fmt.Printf("%s"+NL, tags[len(tags)-1])
		}
	}
}
