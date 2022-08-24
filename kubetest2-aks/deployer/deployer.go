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

package deployer

import (
	"flag"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"sigs.k8s.io/kubetest2/kubetest2-aks/deployer/options"
	"sigs.k8s.io/kubetest2/pkg/types"
)

// Name is the name of the kubetest2 deployer
const Name = "aks"

var (
	GitTag string
	_      types.NewDeployer = New
	_      types.Deployer    = &aksDeployer{}
)

// New implements deployer.New for aks
func New(opts types.Options) (types.Deployer, *pflag.FlagSet) {
	ad := newDeployer(opts)

	// register flags
	fs := bindFlags(ad)

	// register flags for klog
	klog.InitFlags(nil)
	fs.AddGoFlagSet(flag.CommandLine)
	return ad, fs
}

func newDeployer(opts types.Options) *aksDeployer {
	return &aksDeployer{
		ClusterOptions: &options.ClusterOptions{
			Template: "",
		},
	}
}

type aksDeployer struct {
	*options.ClusterOptions

	clusterResourceID string
}

// DumpClusterLogs should export logs from the cluster. It may be called
// multiple times. Options for this should come from New(...)
func (ad *aksDeployer) DumpClusterLogs() error {
	panic("not implemented") // TODO: Implement
}

// Build should build kubernetes and package it in whatever format
// the deployer consumes
func (ad *aksDeployer) Build() error {
	panic("not implemented") // TODO: Implement
}

func bindFlags(d *aksDeployer) *pflag.FlagSet {
	flags, err := gpflag.Parse(d)
	if err != nil {
		klog.Fatalf("unable to generate flags from deployer")
		return nil
	}

	return flags
}
