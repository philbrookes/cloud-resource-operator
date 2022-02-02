// +kubebuilder:object:generate=true
// +groupName=aro.openshift.io

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +groupName=aro.openshift.io
// +kubebuilder:skip

// Cluster holds cluster-wide information about Infrastructure.  The canonical name is `cluster`
type Cluster struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +required
	Spec ClusterSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ClusterStatus `json:"status"`
}

type ClusterSpec struct {
	VNetID  			string `json:"vnetId,omitempty"`
	Domain                 string              `json:"domain,omitempty"`
	ArchitectureVersion    int                 `json:"architectureVersion,omitempty"`
	ResourceID             string              `json:"resourceId,omitempty"`
	InfraID                string              `json:"infraId,omitempty"`
	ClusterResourceGroupID string              `json:"clusterResourceGroupId,omitempty"`
	ACRDomain              string              `json:"acrDomain,omitempty"`
	Location               string              `json:"location,omitempty"`
	AZEnvironment          string              `json:"azEnvironment,omitempty"`
	IngressIP              string              `json:"ingressIP,omitempty"`
	APIIntIP               string              `json:"apiIntIP,omitempty"`
	Features               ClusterSpecFeatures `json:"features,omitempty"`
}

type ClusterSpecFeatures struct {
	ReconcileWorkaroundsController bool `json:"reconcileWorkaroundsController"`
	ReconcilePullSecret bool `json:"reconcilePullSecret"`
	ReconcileMachineSet bool `json:"reconcileMachineSet"`
	ReconcileAlertWebhook bool `json:"reconcileAlertWebhook"`
	ReconcileDNSMasq bool `json:"reconcileDNSMasq"`
	ReconcileGenevaLogging bool `json:"reconcileGenevaLogging"`
	ReconcileNodeDrainer bool `json:"reconcileNodeDrainer"`
	ReconcileMonitoringConfig bool `json:"reconcileMonitoringConfig"`
	ReconcileRouteFix bool `json:"reconcileRouteFix"`
}


type ClusterStatus struct {
	OperatorVersion string `json:"operatorVersion,omitempty"`
	RedHatKeysPresent []string `json:"redHatKeysPresent,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +groupName=aro.openshift.io

// ClusterList is
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
