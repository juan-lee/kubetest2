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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"k8s.io/klog/v2"

	"sigs.k8s.io/kubetest2/pkg/exec"
)

type armDeployment struct {
	Properties armDeploymentProperties `json:"properties"`
}

type armDeploymentProperties struct {
	OutputResources []armDeploymentOutputResources `json:"outputResources"`
}

type armDeploymentOutputResources struct {
	ID            string `json:"id"`
	ResourceGroup string `json:"resourceGroup"`
}

// Up should provision a new cluster for testing
func (ad *aksDeployer) Up() error {
	klog.Infof("Creating deployment: %+v", *ad)

	resourceID, err := ensureDeployment(ad.ClusterOptions.ResourceGroup, ad.ClusterOptions.Template)
	if err != nil {
		klog.ErrorS(err, "Failed to create deployment")
		return err
	}
	ad.clusterResourceID = resourceID
	klog.Info("Successful deployment")

	clusterName, err := parseClusterName(resourceID)
	if err != nil {
		return err
	}

	err = ensureKubeconfig(ad.ClusterOptions.ResourceGroup, clusterName)
	if err != nil {
		klog.ErrorS(err, "Failed to create deployment")
		return err
	}

	return nil
}

// IsUp should return true if a test cluster is successfully provisioned
func (ad *aksDeployer) IsUp() (up bool, err error) {
	out, err := runWithErrorOutput(
		exec.Command("az", "deployment", "group", "create", "-w", "-r", "ResourceIdOnly",
			"-g", ad.ClusterOptions.ResourceGroup, "-f", ad.ClusterOptions.Template,
		))
	if err != nil {
		return false, fmt.Errorf("Failed to query deployment [%s]: %w", ad.ClusterOptions.Template, err)
	}
	resourceID, err := parseWhatIf(out)
	if err != nil {
		return false, nil
	}
	_, err = runWithErrorOutput(exec.Command("az", "resource", "show", "--ids", resourceID))
	if err != nil {
		klog.Infof("Cluster does not exist: %v", err)
		return false, nil
	}
	ad.clusterResourceID = resourceID
	return true, nil
}

func ensureDeployment(group, template string) (string, error) {
	out, err := runWithErrorOutput(exec.Command("az", "deployment", "group", "create", "-g", group, "-f", template))
	if err != nil {
		return "", fmt.Errorf("Failed to create deployment [%s]: %w", template, err)
	}

	resourceID, err := parseResourceID(out)
	if err != nil {
		return "", err
	}
	return resourceID, err
}

func ensureKubeconfig(group, clusterName string) error {
	kubeconfig, err := runWithErrorOutput(exec.Command("az", "aks", "get-credentials", "-g", group, "-n", clusterName, "-f", "-"))
	if err != nil {
		return fmt.Errorf("Failed to get kubeconfig: %w", err)
	}

	tmpdir, err := ioutil.TempDir("", "kubetest2-aks")
	if err != nil {
		return err
	}

	filename := filepath.Join(tmpdir, fmt.Sprintf("kubeconfig-%s-%s", group, clusterName))
	if err := ioutil.WriteFile(filename, kubeconfig, 0644); err != nil {
		return err
	}

	if err := os.Setenv("KUBECONFIG", filename); err != nil {
		return err
	}
	return nil
}

func runWithErrorOutput(cmd exec.Cmd) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd.SetStdout(&stdout)
	cmd.SetStderr(&stderr)
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", stderr.Bytes(), err)
	}
	return stdout.Bytes(), nil
}

func parseWhatIf(output []byte) (string, error) {
	lines := bytes.Split(output, []byte{'\n'})
	var scope string
	for _, line := range lines {
		scopeRegex := regexp.MustCompile(`Scope: (.+)`)
		scopeMatch := scopeRegex.FindStringSubmatch(string(line))
		if len(scopeMatch) == 2 {
			scope = scopeMatch[1]
			continue
		}
		re := regexp.MustCompile(`\+ Microsoft.ContainerService/managedClusters/(.+)`)
		match := re.FindStringSubmatch(string(line))
		if len(match) == 2 {
			return scope + "/Microsoft.ContainerService/managedClusters/" + match[1], nil
		}

	}
	return "", fmt.Errorf("resourceID not found")
}

func parseResourceID(output []byte) (string, error) {
	var armOutput armDeployment
	if err := json.Unmarshal(output, &armOutput); err != nil {
		return "", err
	}
	return armOutput.Properties.OutputResources[0].ID, nil
}

func parseClusterName(resourceID string) (string, error) {
	re := regexp.MustCompile(`(?i)/subscriptions/(.+)/resourcegroups/(.+)/providers/Microsoft.ContainerService/managedClusters/(.+)$`)
	match := re.FindStringSubmatch(resourceID)
	if len(match) != 4 {
		return "", fmt.Errorf("invalid cluster resourceID")
	}
	return match[3], nil
}
