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

package cmd

import (
	"fmt"
	"os"

	"github.com/k8sbykeshed/op-readiness/pkg/flags"
	"github.com/k8sbykeshed/op-readiness/pkg/testcases"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/tools/clientcmd"
)

// NewLoggerConfig return the configuration object for the logger
func NewLoggerConfig(options ...zap.Option) *zap.Logger {
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		NameKey:     "logger",
		TimeKey:     "timer",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		EncodeTime:  zapcore.RFC3339TimeEncoder,
	}), os.Stdout, zap.InfoLevel)
	return zap.New(core).WithOptions(options...)
}

var (
	E2EBinary  string
	provider   string
	testFile   string
	kubeconfig string
	categories flags.ArrayFlags

	rootCmd = &cobra.Command{
		Use:   "op-readiness",
		Short: "The Windows Operational Readiness testing suite",
		Long:  "Run this software and make sure your Windows node is suitable for Kubernetes operations.",
		Run: func(cmd *cobra.Command, args []string) {
			zap.ReplaceGlobals(NewLoggerConfig())

			opTestConfig, err := testcases.NewOpTestConfig(testFile)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Create op-readiness context failed, error is %v", zap.Error(err)))
				os.Exit(1)
			}
			testCtx := testcases.NewTestContext(E2EBinary, kubeconfig, provider, opTestConfig, categories)
			failedTest := false

			for i, t := range opTestConfig.OpTestCases {
				if !testCtx.CategoryEnabled(t.Category) {
					zap.L().Info(fmt.Sprintf("Skipping Operational Readiness Test %v / %v : %v is not in Category %v", i+1, len(opTestConfig.OpTestCases), t.Description, t.Category))
					continue
				}

				zap.L().Info(fmt.Sprintf("Running Operational Readiness Test %v / %v : %v on %v", i+1, len(opTestConfig.OpTestCases), t.Description, t.Category))
				if err = t.RunTest(testCtx); err != nil {
					zap.L().Error(fmt.Sprintf("Operational Readiness Test %v failed, error is %v", t.Description, zap.Error(err)))
					failedTest = true
				}
			}

			if failedTest {
				os.Exit(1)
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&testFile, "test-file", "tests.yaml", "Path to YAML file containing the tests.")
	rootCmd.PersistentFlags().StringVar(&E2EBinary, "e2e-binary", "./e2e.test", "The E2E Ginkgo default binary used to run the tests.")
	rootCmd.PersistentFlags().StringVar(&provider, "provider", "local", "The name of the Kubernetes provider (gce, gke, local, skeleton (the fallback if not set), etc.)")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, clientcmd.RecommendedConfigPathFlag, os.Getenv(clientcmd.RecommendedConfigPathEnvVar), "Path to kubeconfig containing embedded authinfo.")
	rootCmd.PersistentFlags().Var(&categories, "category", "Append category of tests you want to run, default empty will run all tests.")
}
