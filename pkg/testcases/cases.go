/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testcases

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type OpTestCase struct {
	Category           string   `yaml:"category,omitempty"`
	Focus              []string `yaml:"focus,omitempty"`
	Skip               []string `yaml:"skip,omitempty"`
	KubernetesVersions []string `yaml:"kubernetesVersions,omitempty"`
	WindowsPodImage    string   `yaml:"windows_image,omitempty"`
	LinuxPodImage      string   `yaml:"linux_image,omitempty"`
	Description        string   `yaml:"description,omitempty"`
}

// RunTest runs the binary set in the test context with the parameters from flags
func (o *OpTestCase) RunTest(testCtx *TestContext) error {
	args := []string{
		"--provider", testCtx.Provider,
		"--kubeconfig", testCtx.KubeConfig,
		"--node-os-distro", "windows",
		"--non-blocking-taints", "os",
		"--ginkgo.flakeAttempts", "1",
	}

	if len(o.Focus) > 0 {
		focus := o.Focus[0]
		for k, f := range o.Focus {
			if k > 0 {
				focus = focus + "|" + f
			}
		}
		args = append(args, "--ginkgo.focus")
		args = append(args, focus)
	}

	if len(o.Skip) > 0 {
		skip := o.Skip[0]
		for k, s := range o.Skip {
			if k > 0 {
				skip = skip + "|" + s
			}
		}
		args = append(args, "--ginkgo.skip")
		args = append(args, skip)
	}

	cmd := exec.Command(testCtx.E2EBinary, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	cmd.Start()

	redirectOutput(stdout)
	redirectOutput(stderr)

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func redirectOutput(stdout io.ReadCloser) {
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}
