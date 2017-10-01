// +build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	fmt.Printf("hi, tag is '%s'\n", getTag())
	var version string
	tag := getTag()
	if tag != "" && strings.HasPrefix(tag, "v") {
		version = strings.Split(tag, "v")[1]
	} else {
		c := getLatestTag()
		version = strings.Split(getTagForCommit(c), "v")[1]
		maj, min, patch := splitVersion(version)
		patch += 1
		version = fmt.Sprintf("%d.%d.%d", maj, min, patch)
	}

	fmt.Printf("Version is %s\n", version)
}

// #!/bin/bash
// if (git describe --abbrev=0 --exact-match &>/dev/null); then
//   # We're on a tagged commit - use that as the version
//   git describe --abbrev=0 --exact-match | sed 's/v\(.*\)/\1/'
// else
//   # Get the latest tagged version (if there is one)
//   tags=$(git rev-list --tags --max-count=1 2>/dev/null)
//   if [ "$tags" == "" ]; then
//     v="0.0.0"
//   else
//     v=$(git describe --abbrev=0 --tags $tags 2>/dev/null | sed 's/v\(.*\)/\1/')
//   fi
//   # Split by period into an array
//   a=( ${v//./ } )
//   # Increment the patch-level number
//   (( a[2]++ ))
//   # This is a pre-release - locally it gets '-dev', in CircleCI it gets the build number
//   echo "${a[0]}.${a[1]}.${a[2]}-${CIRCLE_BUILD_NUM:-dev}"
// fi

func splitVersion(version string) (maj, min, patch int64) {
	v := strings.SplitN(version, ".", 3)
	maj, _ = strconv.ParseInt(v[0], 0, 0)
	min, _ = strconv.ParseInt(v[1], 0, 0)
	patch, _ = strconv.ParseInt(v[2], 0, 0)
	return maj, min, patch
}

// getTag - Tet the current commit's tag, if any. Otherwise an empty string.
func getTag() string {
	t, err := runError("git", "describe", "--abbrev=0", "--exact-match")
	if err != nil {
		return ""
	}
	return string(t)
}

func getLatestTag() string {
	t, err := runError("git", "rev-list", "--tags", "--max-count=1")
	if err != nil {
		return ""
	}
	return string(t)
}

func getTagForCommit(commit string) string {
	t, err := runError("git", "describe", "--abbrev=0", "--tags", commit)
	if err != nil {
		return ""
	}
	return string(t)
}

func getGitSha() string {
	v, err := runError("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "unknown-dev"
	}
	return string(v)
}

func run(cmd string, args ...string) []byte {
	bs, err := runError(cmd, args...)
	if err != nil {
		log.Println(cmd, strings.Join(args, " "))
		log.Println(string(bs))
		log.Fatal(err)
	}
	return bytes.TrimSpace(bs)
}

func runError(cmd string, args ...string) ([]byte, error) {
	ecmd := exec.Command(cmd, args...)
	bs, err := ecmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(bs), nil
}

func runPrint(cmd string, args ...string) {
	log.Println(cmd, strings.Join(args, " "))
	ecmd := exec.Command(cmd, args...)
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	err := ecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
