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
	"fmt"

	"k8s.io/klog/v2"
	"sigs.k8s.io/kubetest2/pkg/exec"
)

// Down should tear down the test cluster if any
func (ad *aksDeployer) Down() error {
	up, err := ad.IsUp()
	if err != nil {
		return fmt.Errorf("failed to get cluster state: %w", err)
	}

	if !up {
		klog.Info("No cluster, skipping down.")
		return nil
	}

	klog.Infof("Deleting Resources: %+v", *ad)
	if err := deleteResource(ad.clusterResourceID); err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}
	klog.Info("Deleted resources")
	return nil
}

func deleteResource(resourceID string) error {
	_, err := runWithErrorOutput(exec.Command("az", "resource", "delete", "--ids", resourceID))
	return err
}
