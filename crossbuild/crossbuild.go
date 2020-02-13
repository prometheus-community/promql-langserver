// Copyright 2019 Tobias Guggenmos
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type architecture struct {
	os   string
	arch string
}

var targets = []architecture{
	architecture{os: "linux", arch: "386"},
	architecture{os: "linux", arch: "amd64"},
	architecture{os: "linux", arch: "arm"},
	architecture{os: "linux", arch: "arm64"},
	architecture{os: "darwin", arch: "386"},
	architecture{os: "darwin", arch: "amd64"},
	architecture{os: "darwin", arch: "arm"},
	architecture{os: "darwin", arch: "arm64"},
	architecture{os: "windows", arch: "386"},
	architecture{os: "windows", arch: "amd64"},
}

func main() {
	versionFileContent, err := ioutil.ReadFile("VERSION")
	if err != nil {
		log.Fatal(err)
	}

	version := strings.TrimSpace(string(versionFileContent))

	for _, target := range targets {
		fmt.Printf("Building %s/%s:\n", target.os, target.arch)

		outputDir := fmt.Sprint(".build/promql-langserver-", version, ".", target.os, "-", target.arch)
		tarName := fmt.Sprint("out/promql-langserver-", version, ".", target.os, "-", target.arch, ".tar.gz")
		//zipName := fmt.Sprint("promql-langserver-", version, ".", target.os, "-", target.arch, ".zip")

		if err := exec.Command("rm", "-rf", outputDir).Run(); err != nil {
			log.Fatal(err)
		}

		if err := exec.Command("mkdir", "-p", outputDir, "out").Run(); err != nil {
			log.Fatal(err)
		}

		if err := exec.Command("cp", "LICENSE", outputDir).Run(); err != nil {
			log.Fatal(err)
		}

		cmd := exec.Command("go", "build", "-o", outputDir, "./cmd/promql-langserver")

		cmd.Env = append(os.Environ(), fmt.Sprint("GOOS=", target.os), fmt.Sprint("GOARCH", target.arch))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		if err := exec.Command("tar", "-zcvf", tarName, "-C", outputDir, ".").Run(); err != nil {
			log.Fatal(err)
		}

		// We don't need zip files at the moment.
		/*
			if err := exec.Command("zip", "-r", zipName, outputDir).Run(); err != nil {
				log.Fatal(err)
			}
		*/
	}
}
