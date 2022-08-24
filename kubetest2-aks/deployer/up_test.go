package deployer

import "testing"

func TestParseOutput(t *testing.T) {
	output := []byte(`
{
  "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tekton-clusters/providers/Microsoft.Resources/deployments/cluster",
  "location": null,
  "name": "cluster",
  "properties": {
    "correlationId": "300460db-412b-443b-94e4-8de7f967049b",
    "debugSetting": null,
    "dependencies": [],
    "duration": "PT1M38.6080921S",
    "error": null,
    "mode": "Incremental",
    "onErrorDeployment": null,
    "outputResources": [
      {
        "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tekton-clusters/providers/Microsoft.ContainerService/managedClusters/bicep-test",
        "resourceGroup": "tekton-clusters"
      }
    ],
    "outputs": null,
    "parameters": {
      "clusterName": {
        "type": "String",
        "value": "bicep-test"
      },
      "location": {
        "type": "String",
        "value": "southcentralus"
      }
    },
    "parametersLink": null,
    "providers": [
      {
        "id": null,
        "namespace": "Microsoft.ContainerService",
        "providerAuthorizationConsentState": null,
        "registrationPolicy": null,
        "registrationState": null,
        "resourceTypes": [
          {
            "aliases": null,
            "apiProfiles": null,
            "apiVersions": null,
            "capabilities": null,
            "defaultApiVersion": null,
            "locationMappings": null,
            "locations": [
              "southcentralus"
            ],
            "properties": null,
            "resourceType": "managedClusters",
            "zoneMappings": null
          }
        ]
      }
    ],
    "provisioningState": "Succeeded",
    "templateHash": "13029678253376033917",
    "templateLink": null,
    "timestamp": "2022-08-17T13:34:40.217477+00:00",
    "validatedResources": null
  },
  "resourceGroup": "tekton-clusters",
  "tags": null,
  "type": "Microsoft.Resources/deployments"
}`)
	expectedResult := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tekton-clusters/providers/Microsoft.ContainerService/managedClusters/bicep-test"
	resourceID, err := parseResourceID(output)
	if err != nil {
		t.Fatalf("parseOutput failed: %s", err)
	}
	if resourceID != expectedResult {
		t.Fatalf("expected [%s], actual [%s]", expectedResult, resourceID)
	}
}

func TestParseResourceID(t *testing.T) {
	resourceID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/tekton-clusters/providers/Microsoft.ContainerService/managedClusters/bicep-test"
	expectedResult := "bicep-test"
	clusterName, err := parseClusterName(resourceID)
	if err != nil {
		t.Fatalf("parseOutput failed: %s", err)
	}
	if clusterName != expectedResult {
		t.Fatalf("expected [%s], actual [%s]", expectedResult, clusterName)
	}
}

func TestParseWhatIf(t *testing.T) {
	output := []byte(`Note: The result may contain false positive predictions (noise).
You can help us improve the accuracy of the result by opening an issue here: https://aka.ms/WhatIfIssues

Resource and property changes are indicated with these symbols:
  + Create
  * Ignore

The deployment will update the following scope:

Scope: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-clusters

  + Microsoft.ContainerService/managedClusters/bicep-test
  * Microsoft.ContainerService/managedClusters/foobar

Resource changes: 1 to create, 9 to ignore.
`)
	expectedResult := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-clusters/Microsoft.ContainerService/managedClusters/bicep-test"
	resourceID, err := parseWhatIf(output)
	if err != nil {
		t.Fatalf("parseOutput failed: %s", err)
	}
	if resourceID != expectedResult {
		t.Fatalf("expected [%s], actual [%s]", expectedResult, resourceID)
	}
}
