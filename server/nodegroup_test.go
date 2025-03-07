package server

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	managednodeClientset "github.com/Fred78290/kubernetes-vmware-autoscaler/pkg/generated/clientset/versioned"
	"github.com/Fred78290/kubernetes-vmware-autoscaler/types"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
)

type mockupClientGenerator struct {
}

func (m mockupClientGenerator) KubeClient() (kubernetes.Interface, error) {
	return nil, nil
}

func (m mockupClientGenerator) NodeManagerClient() (managednodeClientset.Interface, error) {
	return nil, nil
}

func (m mockupClientGenerator) ApiExtentionClient() (apiextension.Interface, error) {
	return nil, nil
}

func (m mockupClientGenerator) PodList(nodeName string, podFilter types.PodFilterFunc) ([]apiv1.Pod, error) {
	return nil, nil
}

func (m mockupClientGenerator) NodeList() (*apiv1.NodeList, error) {
	return &apiv1.NodeList{}, nil
}

func (m mockupClientGenerator) UncordonNode(nodeName string) error {
	return nil
}

func (m mockupClientGenerator) CordonNode(nodeName string) error {
	return nil
}

func (m mockupClientGenerator) SetProviderID(nodeName, providerID string) error {
	return nil
}

func (m mockupClientGenerator) MarkDrainNode(nodeName string) error {
	return nil
}

func (m mockupClientGenerator) DrainNode(nodeName string, ignoreDaemonSet, deleteLocalData bool) error {
	return nil
}

func (m mockupClientGenerator) GetNode(nodeName string) (*apiv1.Node, error) {
	return nil, nil
}

func (m mockupClientGenerator) DeleteNode(nodeName string) error {
	return nil
}

func (m mockupClientGenerator) AnnoteNode(nodeName string, annotations map[string]string) error {
	return nil
}

func (m mockupClientGenerator) LabelNode(nodeName string, labels map[string]string) error {
	return nil
}

func (m mockupClientGenerator) TaintNode(nodeName string, taints ...apiv1.Taint) error {
	return nil
}

func (m mockupClientGenerator) WaitNodeToBeReady(nodeName string, timeToWaitInSeconds int) error {
	return nil
}

func createTestNode(ng *AutoScalerServerNodeGroup) *AutoScalerServerNode {
	return &AutoScalerServerNode{
		NodeGroupID:   testGroupID,
		NodeName:      testNodeName,
		Memory:        2048,
		CPU:           2,
		Disk:          5120,
		IPAddress:     "127.0.0.1",
		State:         AutoScalerServerNodeStateNotCreated,
		NodeType:      AutoScalerServerNodeAutoscaled,
		VSphereConfig: ng.configuration.GetVSphereConfiguration(testGroupID),
		serverConfig:  ng.configuration,
	}
}

func newTestNode() (*types.AutoScalerServerConfig, *AutoScalerServerNodeGroup, *AutoScalerServerNode, error) {
	config, ng, err := newTestNodeGroup()

	if err == nil {
		vm := createTestNode(ng)

		ng.Nodes[testGroupID] = vm

		return config, ng, vm, err
	}

	return config, ng, nil, err
}

func newTestNodeGroup() (*types.AutoScalerServerConfig, *AutoScalerServerNodeGroup, error) {
	config, err := newTestConfig()

	if err == nil {
		ng := &AutoScalerServerNodeGroup{
			ServiceIdentifier:          testServiceIdentifier,
			NodeGroupIdentifier:        testGroupID,
			ProvisionnedNodeNamePrefix: "autoscaled",
			ManagedNodeNamePrefix:      "worker",
			ControlPlaneNamePrefix:     "master",
			Machine: &types.MachineCharacteristic{
				Memory: 4096,
				Vcpu:   4,
				Disk:   5120,
			},
			Status:      NodegroupNotCreated,
			MinNodeSize: 0,
			MaxNodeSize: 5,
			NodeLabels: types.KubernetesLabel{
				"monitor":  "true",
				"database": "true",
			},
			SystemLabels:  types.KubernetesLabel{},
			Nodes:         make(map[string]*AutoScalerServerNode),
			pendingNodes:  make(map[string]*AutoScalerServerNode),
			configuration: config,
		}

		return config, ng, err
	}

	return nil, nil, err
}

func newTestConfig() (*types.AutoScalerServerConfig, error) {
	var config types.AutoScalerServerConfig

	configStr, _ := ioutil.ReadFile("./masterkube/config/config.json")
	err := json.Unmarshal(configStr, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func Test_AutoScalerNode_launchVM(t *testing.T) {
	kubeClient := &mockupClientGenerator{}
	_, ng, testNode, err := newTestNode()

	if assert.NoError(t, err) {
		t.Run("Launch VM", func(t *testing.T) {
			if err := testNode.launchVM(kubeClient, ng.NodeLabels, ng.SystemLabels); err != nil {
				t.Errorf("AutoScalerNode.launchVM() error = %v", err)
			}
		})
	}
}

func Test_AutoScalerNode_startVM(t *testing.T) {
	kubeClient := &mockupClientGenerator{}
	_, _, testNode, err := newTestNode()

	if assert.NoError(t, err) {
		t.Run("Start VM", func(t *testing.T) {
			if err := testNode.startVM(kubeClient); err != nil {
				t.Errorf("AutoScalerNode.startVM() error = %v", err)
			}
		})
	}
}

func Test_AutoScalerNode_stopVM(t *testing.T) {
	kubeClient := &mockupClientGenerator{}
	_, _, testNode, err := newTestNode()

	if assert.NoError(t, err) {
		t.Run("Stop VM", func(t *testing.T) {
			if err := testNode.stopVM(kubeClient); err != nil {
				t.Errorf("AutoScalerNode.stopVM() error = %v", err)
			}
		})
	}
}

func Test_AutoScalerNode_deleteVM(t *testing.T) {
	kubeClient := &mockupClientGenerator{}
	_, _, testNode, err := newTestNode()

	if assert.NoError(t, err) {
		t.Run("Delete VM", func(t *testing.T) {
			if err := testNode.deleteVM(kubeClient); err != nil {
				t.Errorf("AutoScalerNode.deleteVM() error = %v", err)
			}
		})
	}
}

func Test_AutoScalerNode_statusVM(t *testing.T) {
	_, _, testNode, err := newTestNode()

	if assert.NoError(t, err) {
		t.Run("Status VM", func(t *testing.T) {
			if got, err := testNode.statusVM(); err != nil {
				t.Errorf("AutoScalerNode.statusVM() error = %v", err)
			} else if got != AutoScalerServerNodeStateRunning {
				t.Errorf("AutoScalerNode.statusVM() = %v, want %v", got, AutoScalerServerNodeStateRunning)
			}
		})
	}
}

func Test_AutoScalerNodeGroup_addNode(t *testing.T) {
	kubeClient := &mockupClientGenerator{}
	_, ng, err := newTestNodeGroup()

	if assert.NoError(t, err) {
		t.Run("addNode", func(t *testing.T) {
			if _, err := ng.addNodes(kubeClient, 1); err != nil {
				t.Errorf("AutoScalerServerNodeGroup.addNode() error = %v", err)
			}
		})
	}
}

func Test_AutoScalerNodeGroup_deleteNode(t *testing.T) {
	kubeClient := &mockupClientGenerator{}
	_, ng, testNode, err := newTestNode()

	if assert.NoError(t, err) {
		t.Run("Delete VM", func(t *testing.T) {
			if err := ng.deleteNodeByName(kubeClient, testNode.NodeName); err != nil {
				t.Errorf("AutoScalerServerNodeGroup.deleteNode() error = %vv", err)
			}
		})
	}
}

func Test_AutoScalerNodeGroup_deleteNodeGroup(t *testing.T) {
	kubeClient := &mockupClientGenerator{}
	_, ng, _, err := newTestNode()

	if assert.NoError(t, err) {

		t.Run("Delete node group", func(t *testing.T) {
			if err := ng.deleteNodeGroup(kubeClient); err != nil {
				t.Errorf("AutoScalerServerNodeGroup.deleteNodeGroup() error = %v", err)
			}
		})
	}
}
