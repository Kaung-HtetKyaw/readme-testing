package postgres

import (
	"encoding/json"

	"github.com/kubegrade/matrix"
)

func objectHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	// NOTE let's try switch approch until we come up with better ideas (generic?).
	switch kind {
	case "Pod":
		return podHealthStatus(kind, raw)
	case "Deployment":
		return deploymentHealthStatus(kind, raw)
	case "Job":
		return jobHealthStatus(kind, raw)
	case "ReplicaSet":
		return replicaSetHealthStatus(kind, raw)
	case "ReplicationController":
		return replicationControllerHealthStatus(kind, raw)
	case "StatefulSet":
		return statefulSetHealthStatus(kind, raw)
	case "PersistentVolume":
		return persistentVolumeHealthStatus(kind, raw)
	case "PersistentVolumeClaim":
		return persistentVolumeClaimHealthStatus(kind, raw)
	case "VolumeAttachment":
		return volumeAttachmentHealthStatus(kind, raw)
	case "ValidatingAdmissionPolicyBinding":
		return validatingAdmissionPolicyBindingHealthStatus(kind, raw)
	case "ClusterRole":
		return clusterRoleHealthStatus(kind, raw)
	case "Role":
		return roleHealthStatus(kind, raw)
	case "APIService":
		return apiServiceHealthStatus(kind, raw)
	case "CertificateSigningRequest":
		return certificateSigningRequestHealthStatus(kind, raw)
	case "ComponentStatus":
		return componentStatusHealthStatus(kind, raw)
	case "DeviceClass":
		return deviceClassHealthStatus(kind, raw)
	case "FlowSchema":
		return flowSchemaHealthStatus(kind, raw)
	case "PriorityLevelConfiguration":
		return priorityLevelConfigurationHealthStatus(kind, raw)
	case "ResourceClaim":
		return resourceClaimHealthStatus(kind, raw)
	case "ResourceSlice":
		return resourceSliceHealthStatus(kind, raw)
	case "ServiceCIDR":
		return serviceCIDRHealthStatus(kind, raw)
	case "StorageVersion":
		return storageVersionHealthStatus(kind, raw)
	case "StorageVersionMigration":
		return storageVersionMigrationHealthStatus(kind, raw)
	case "Namespace":
		return namespaceHealthStatus(kind, raw)
	case "Node":
		return nodeHealthStatus(kind, raw)
	case "CronJob":
		return cronJobHealthStatus(kind, raw)
	case "Ingress":
		return ingressHealthStatus(kind, raw)
	case "IngressClass":
		return ingressClassHealthStatus(kind, raw)
	case "NetworkPolicy":
		return networkPolicyHealthStatus(kind, raw)
	case "HorizontalPodAutoscaler":
		return horizontalPodAutoscalerHealthStatus(kind, raw)
	case "PodDisruptionBudget":
		return podDisruptionBudgetHealthStatus(kind, raw)
	case "PriorityClass":
		return priorityClassHealthStatus(kind, raw)
	case "Service":
		return serviceHealthStatus(kind, raw)
	case "ConfigMap":
		return configMapHealthStatus(kind, raw)
	case "LimitRange":
		return limitRangeHealthStatus(kind, raw)
	case "ResourceQuota":
		return resourceQuotaHealthStatus(kind, raw)
	case "Secret":
		return secretHealthStatus(kind, raw)
	case "CSIDriver":
		return csiDriverHealthStatus(kind, raw)
	case "CSINode":
		return csiNodeHealthStatus(kind, raw)
	case "CSIStorageCapacity":
		return csiStorageCapacityHealthStatus(kind, raw)
	case "StorageClass":
		return storageClassHealthStatus(kind, raw)
	case "MutatingAdmissionPolicy":
		return mutatingAdmissionPolicyHealthStatus(kind, raw)
	case "MutatingAdmissionPolicyBinding":
		return mutatingAdmissionPolicyBindingHealthStatus(kind, raw)
	case "MutatingWebhookConfiguration":
		return mutatingWebhookConfigurationHealthStatus(kind, raw)
	case "RuntimeClass":
		return runtimeClassHealthStatus(kind, raw)
	case "ValidatingAdmissionPolicy":
		return validatingAdmissionPolicyHealthStatus(kind, raw)
	case "ValidatingWebhookConfiguration":
		return validatingWebhookConfigurationHealthStatus(kind, raw)
	case "ClusterRoleBinding":
		return clusterRoleBindingHealthStatus(kind, raw)
	case "RoleBinding":
		return roleBindingHealthStatus(kind, raw)
	case "ServiceAccount":
		return serviceAccountHealthStatus(kind, raw)
	case "ClusterTrustBundle":
		return clusterTrustBundleHealthStatus(kind, raw)
	case "ControllerRevision":
		return controllerRevisionHealthStatus(kind, raw)
	case "DeviceTaintRule":
		return deviceTaintRuleHealthStatus(kind, raw)
	case "IPAddress":
		return ipAddressHealthStatus(kind, raw)
	case "LeaseCandidate":
		return leaseCandidateHealthStatus(kind, raw)
	case "PodTemplate":
		return podTemplateHealthStatus(kind, raw)
	case "ResourceClaimTemplate":
		return resourceClaimTemplateHealthStatus(kind, raw)
	case "VolumeAttributesClass":
		return volumeAttributesClassHealthStatus(kind, raw)
	}
	return matrix.ObjectHealthStatusUnknown, nil
}

func podHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var pod struct {
		Status struct {
			Phase      string `json:"phase"`
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
			ContainerStatuses []struct {
				RestartCount int `json:"restartCount"`
			} `json:"containerStatuses"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &pod); err != nil {
		return "", err
	}

	// Priority 1: Broken
	// Check Ready condition
	for _, condition := range pod.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "False" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}
	// Check phase
	if pod.Status.Phase == "Failed" {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// Check phase=Pending
	if pod.Status.Phase == "Pending" {
		return matrix.ObjectHealthStatusWarning, nil
	}
	// Check restartCount > 5 (Kubernetes considers 5+ restarts as problematic but not necessarily broken)
	for _, container := range pod.Status.ContainerStatuses {
		if container.RestartCount > 5 {
			return matrix.ObjectHealthStatusWarning, nil
		}
	}
	// Check PodScheduled=False
	for _, condition := range pod.Status.Conditions {
		if condition.Type == "PodScheduled" && condition.Status == "False" {
			return matrix.ObjectHealthStatusWarning, nil
		}
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func deploymentHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var deployment struct {
		Spec struct {
			Replicas *int32 `json:"replicas"`
		} `json:"spec"`
		Status struct {
			Replicas            int32 `json:"replicas"`
			ReadyReplicas       int32 `json:"readyReplicas"`
			UnavailableReplicas int32 `json:"unavailableReplicas"`
			Conditions          []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &deployment); err != nil {
		return "", err
	}

	// Priority 1: Broken
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == "Available" && condition.Status == "False" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}
	if deployment.Status.Replicas == 0 && deployment.Spec.Replicas != nil && *deployment.Spec.Replicas > 0 {
		return matrix.ObjectHealthStatusBroken, nil
	}
	// Check if deployment is not progressing
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == "Progressing" && condition.Status == "False" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	if deployment.Status.ReadyReplicas < deployment.Status.Replicas {
		return matrix.ObjectHealthStatusWarning, nil
	}
	if deployment.Status.UnavailableReplicas > 0 {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func jobHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var job struct {
		Status struct {
			Active     int32 `json:"active"`
			Succeeded  int32 `json:"succeeded"`
			Failed     int32 `json:"failed"`
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &job); err != nil {
		return "", err
	}

	// Priority 1: Broken
	// Check if there is a failed condition
	for _, condition := range job.Status.Conditions {
		if condition.Type == "Failed" && condition.Status == "True" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}
	// Check if there are failed jobs
	if job.Status.Failed > 0 {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// Check if there are no active or succeeded jobs (stuck)
	if job.Status.Active == 0 && job.Status.Succeeded == 0 {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func replicaSetHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var replicaSet struct {
		Spec struct {
			Replicas *int32 `json:"replicas"`
		} `json:"spec"`
		Status struct {
			Replicas          int32  `json:"replicas"`
			ReadyReplicas     *int32 `json:"readyReplicas"`
			AvailableReplicas *int32 `json:"availableReplicas"`
			Conditions        []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &replicaSet); err != nil {
		return "", err
	}

	// NOTE: Missing readyReplicas field could mean that
	// 1. it is at initial states before controller updates,
	// 2. it is in between controler update patches
	// 3. no pods are ready - broken (This is a common pattern in Kubernetes API objects where zero values for optional fields are often omitted to reduce payload size.)
	// TODO: Find a way to differentiate case 3 from case 1 and case 2

	// Priority 1: Broken
	// Check if current replicas is 0 but desired replicas is not 0
	if replicaSet.Status.ReadyReplicas != nil && *replicaSet.Status.ReadyReplicas == 0 && replicaSet.Spec.Replicas != nil && *replicaSet.Spec.Replicas > 0 {
		return matrix.ObjectHealthStatusBroken, nil
	}
	// Check if there is a replica failure
	for _, condition := range replicaSet.Status.Conditions {
		if condition.Type == "ReplicaFailure" && condition.Status == "True" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// Check if ready replicas is less than current replicas
	if replicaSet.Status.ReadyReplicas == nil || (*replicaSet.Status.ReadyReplicas < replicaSet.Status.Replicas) {
		return matrix.ObjectHealthStatusWarning, nil
	}
	// Check if available replicas is less than ready replicas
	if replicaSet.Status.AvailableReplicas == nil || (*replicaSet.Status.AvailableReplicas < *replicaSet.Status.ReadyReplicas) {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func replicationControllerHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var replicationController struct {
		Spec struct {
			Replicas *int32 `json:"replicas"`
		} `json:"spec"`
		Status struct {
			Replicas          int32 `json:"replicas"`
			ReadyReplicas     int32 `json:"readyReplicas"`
			AvailableReplicas int32 `json:"availableReplicas"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &replicationController); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check if current replicas is 0 but desired replicas is not 0
	if replicationController.Status.ReadyReplicas == 0 && replicationController.Spec.Replicas != nil && *replicationController.Spec.Replicas > 0 {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// Check if ready replicas is less than current replicas
	if replicationController.Status.ReadyReplicas < replicationController.Status.Replicas {
		return matrix.ObjectHealthStatusWarning, nil
	}
	// Check if available replicas is less than ready replicas
	if replicationController.Status.AvailableReplicas < replicationController.Status.ReadyReplicas {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func statefulSetHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var statefulSet struct {
		Spec struct {
			Replicas *int32 `json:"replicas"`
		} `json:"spec"`
		Status struct {
			Replicas        int32 `json:"replicas"`
			ReadyReplicas   int32 `json:"readyReplicas"`
			CurrentReplicas int32 `json:"currentReplicas"`
			Conditions      []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &statefulSet); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check status.readyReplicas = 0 AND spec.replicas > 0
	if statefulSet.Status.ReadyReplicas == 0 && statefulSet.Spec.Replicas != nil && *statefulSet.Spec.Replicas > 0 {
		return matrix.ObjectHealthStatusBroken, nil
	}
	// Check status.currentReplicas < status.replicas OR status.conditions[Available].status = False
	if statefulSet.Status.CurrentReplicas < statefulSet.Status.Replicas {
		return matrix.ObjectHealthStatusBroken, nil
	}
	for _, condition := range statefulSet.Status.Conditions {
		if condition.Type == "Available" && condition.Status == "False" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// Check status.readyReplicas < status.replicas
	if statefulSet.Status.ReadyReplicas < statefulSet.Status.Replicas {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func persistentVolumeHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var pv struct {
		Status struct {
			Phase string `json:"phase"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &pv); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	if pv.Status.Phase == "Failed" || pv.Status.Phase == "Lost" {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	if pv.Status.Phase == "Pending" {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func persistentVolumeClaimHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var pvc struct {
		Status struct {
			Phase      string `json:"phase"`
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &pvc); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	if pvc.Status.Phase == "Lost" {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	if pvc.Status.Phase == "Pending" {
		return matrix.ObjectHealthStatusWarning, nil
	}
	// Check if PVC is resizing
	for _, condition := range pvc.Status.Conditions {
		if condition.Type == "Resizing" && condition.Status == "True" {
			return matrix.ObjectHealthStatusWarning, nil
		}
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func volumeAttachmentHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var va struct {
		Status struct {
			Attached    bool   `json:"attached"`
			FailedError string `json:"failedError"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &va); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	if va.Status.FailedError != "" || !va.Status.Attached {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func validatingAdmissionPolicyBindingHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var vab struct {
		Spec struct {
			MatchConstraints struct {
				ResourceRules []interface{} `json:"resourceRules"`
			} `json:"matchConstraints"`
			Validations []interface{} `json:"validations"`
		} `json:"spec"`
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &vab); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, condition := range vab.Status.Conditions {
		if condition.Type == "Ready" && condition.Status != "True" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// Check if resource rules is empty or validations is missing
	if len(vab.Spec.MatchConstraints.ResourceRules) == 0 || len(vab.Spec.Validations) == 0 {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func clusterRoleHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var clusterRole struct {
		Spec struct {
			Rules []interface{} `json:"rules"`
		} `json:"spec"`
	}

	if err := json.Unmarshal(raw, &clusterRole); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check if rules array is empty
	if clusterRole.Spec.Rules != nil && len(clusterRole.Spec.Rules) == 0 {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func roleHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var role struct {
		Rules []struct {
			APIGroups []string `json:"apiGroups"`
			Resources []string `json:"resources"`
			Verbs     []string `json:"verbs"`
		} `json:"rules"`
	}

	if err := json.Unmarshal(raw, &role); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Rules are optional, but when present, all three fields (apiGroups, resources, verbs) are mandatory
	if role.Rules != nil {
		// Priority 1: Broken
		// Check if rules array is empty when it should contain rules
		for _, rule := range role.Rules {
			if len(rule.APIGroups) == 0 || len(rule.Resources) == 0 || len(rule.Verbs) == 0 {
				return matrix.ObjectHealthStatusBroken, nil
			}
		}

		// Priority 2: Warning
		// Check if any rule is missing required fields or has empty arrays
		if len(role.Rules) == 0 {
			return matrix.ObjectHealthStatusWarning, nil
		}
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func apiServiceHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var apiService struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &apiService); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, condition := range apiService.Status.Conditions {
		if condition.Type == "Available" && condition.Status == "False" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func certificateSigningRequestHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var csr struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &csr); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, condition := range csr.Status.Conditions {
		if condition.Type == "Failed" && condition.Status == "True" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func componentStatusHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var componentStatus struct {
		Conditions []struct {
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &componentStatus); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, condition := range componentStatus.Conditions {
		if condition.Type == "Healthy" && (condition.Status == "False" || condition.Status == "Unknown") {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func deviceClassHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var deviceClass struct {
		Spec struct {
			Selectors []interface{} `json:"selectors"`
		} `json:"spec"`
	}

	if err := json.Unmarshal(raw, &deviceClass); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// TODO: define broken criterias

	// Priority 2: Warning
	if len(deviceClass.Spec.Selectors) == 0 {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func flowSchemaHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var flowSchema struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &flowSchema); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, condition := range flowSchema.Status.Conditions {
		if condition.Type == "Dangling" && condition.Status == "True" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func priorityLevelConfigurationHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var plc struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &plc); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, condition := range plc.Status.Conditions {
		if condition.Type == "Available" && condition.Status == "False" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func resourceClaimHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var resourceClaim struct {
		Status struct {
			Devices []struct {
				Conditions []struct {
					Type   string `json:"type"`
					Status string `json:"status"`
				} `json:"conditions"`
			} `json:"devices"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &resourceClaim); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, device := range resourceClaim.Status.Devices {
		for _, condition := range device.Conditions {
			if condition.Type == "Available" && condition.Status == "False" {
				return matrix.ObjectHealthStatusBroken, nil
			}
		}
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func resourceSliceHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var resourceSlice struct {
		Spec struct {
			AllNodes bool `json:"allNodes"`
		} `json:"spec"`
	}

	if err := json.Unmarshal(raw, &resourceSlice); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// TODO: define broken criterias

	// Priority 2: Warning
	// not all nodes can access resources in the pool
	if !resourceSlice.Spec.AllNodes {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func serviceCIDRHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var serviceCIDR struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &serviceCIDR); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, condition := range serviceCIDR.Status.Conditions {
		if condition.Type == "Available" && condition.Status != "True" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func storageVersionHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var storageVersion struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &storageVersion); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, condition := range storageVersion.Status.Conditions {
		if condition.Type == "Available" && condition.Status != "True" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func storageVersionMigrationHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var svm struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &svm); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	for _, condition := range svm.Status.Conditions {
		if condition.Type == "Failed" && condition.Status == "True" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func namespaceHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var namespace struct {
		Metadata struct {
			DeletionTimestamp *string `json:"deletionTimestamp"`
		} `json:"metadata"`
		Status struct {
			Phase string `json:"phase"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &namespace); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check if namespace is in terminating phase
	if namespace.Status.Phase == "Terminating" {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// Check if namespace has deletion timestamp (marked for deletion)
	if namespace.Metadata.DeletionTimestamp != nil {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func nodeHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var node struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &node); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check if node is not ready or has unknown ready status
	for _, condition := range node.Status.Conditions {
		if condition.Type == "Ready" && (condition.Status == "False" || condition.Status == "Unknown") {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// Check if node has resource pressure conditions
	for _, condition := range node.Status.Conditions {
		if (condition.Type == "MemoryPressure" || condition.Type == "DiskPressure" || condition.Type == "PIDPressure") && condition.Status == "True" {
			return matrix.ObjectHealthStatusWarning, nil
		}
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func cronJobHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var cronJob struct {
		Spec struct {
			Suspend *bool `json:"suspend"`
		} `json:"spec"`
		Status struct {
			LastSuccessfulTime string `json:"lastSuccessfulTime"`
			LastScheduleTime   string `json:"lastScheduleTime"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &cronJob); err != nil {
		return "", err
	}

	// Priority 1: Broken
	// Check spec.suspend = true
	if cronJob.Spec.Suspend != nil && *cronJob.Spec.Suspend {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// NOTE: using string comparison since it will be in ISO 8601 format but can be converted to time.Time to handle more edge cases
	if cronJob.Status.LastSuccessfulTime != "" && cronJob.Status.LastScheduleTime != "" {
		if cronJob.Status.LastSuccessfulTime < cronJob.Status.LastScheduleTime {
			return matrix.ObjectHealthStatusWarning, nil
		}
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func ingressHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var ingress struct {
		Spec struct {
			Rules          []interface{} `json:"rules"`
			DefaultBackend interface{}   `json:"defaultBackend"`
		} `json:"spec"`
		Status struct {
			LoadBalancer struct {
				Ingress []interface{} `json:"ingress"`
			} `json:"loadBalancer"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &ingress); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check if load balancer ingress is empty while rules or default backend is configured
	hasRules := len(ingress.Spec.Rules) > 0
	hasDefaultBackend := ingress.Spec.DefaultBackend != nil
	hasLoadBalancerIngress := len(ingress.Status.LoadBalancer.Ingress) > 0

	if (hasRules || hasDefaultBackend) && !hasLoadBalancerIngress {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func ingressClassHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var ingressClass struct {
		Spec struct {
			Controller string `json:"controller"`
		} `json:"spec"`
	}

	if err := json.Unmarshal(raw, &ingressClass); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check if controller is not specified
	// Optional: matches no pods
	if ingressClass.Spec.Controller == "" {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func networkPolicyHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var networkPolicy struct {
		Spec struct {
			PodSelector map[string]interface{} `json:"podSelector"`
		} `json:"spec"`
	}

	if err := json.Unmarshal(raw, &networkPolicy); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check if pod selector is empty (no pods will match this policy)
	if len(networkPolicy.Spec.PodSelector) == 0 {
		return matrix.ObjectHealthStatusBroken, nil
	}

	// Priority 2: Warning
	// TODO: define warning criterias

	return matrix.ObjectHealthStatusHealthy, nil
}

func horizontalPodAutoscalerHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var hpa struct {
		Spec struct {
			MaxReplicas int32 `json:"maxReplicas"`
		} `json:"spec"`
		Status struct {
			CurrentReplicas int32 `json:"currentReplicas"`
			Conditions      []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"status"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &hpa); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check if HPA is unable to scale
	for _, condition := range hpa.Status.Conditions {
		if condition.Type == "AbleToScale" && condition.Status == "False" {
			return matrix.ObjectHealthStatusBroken, nil
		}
	}

	// Priority 2: Warning
	// Check if current replicas have reached the maximum limit
	if hpa.Status.CurrentReplicas >= hpa.Spec.MaxReplicas {
		return matrix.ObjectHealthStatusWarning, nil
	}
	// Check if scaling is limited
	for _, condition := range hpa.Status.Conditions {
		if condition.Type == "ScalingLimited" && condition.Status == "True" {
			return matrix.ObjectHealthStatusWarning, nil
		}
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func podDisruptionBudgetHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var pdb struct {
		Status struct {
			CurrentHealthy int32 `json:"currentHealthy"`
			DesiredHealthy int32 `json:"desiredHealthy"`
		} `json:"status"`
	}

	if err := json.Unmarshal(raw, &pdb); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// TODO: define broken criterias

	// Priority 2: Warning
	// Check if current healthy pods are less than desired healthy pods
	if pdb.Status.CurrentHealthy < pdb.Status.DesiredHealthy {
		return matrix.ObjectHealthStatusWarning, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func priorityClassHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	var priorityClass struct {
		Spec struct {
			Value *int32 `json:"value"`
		} `json:"spec"`
	}

	if err := json.Unmarshal(raw, &priorityClass); err != nil {
		return matrix.ObjectHealthStatusUnknown, err
	}

	// Priority 1: Broken
	// Check if value field is negative
	if priorityClass.Spec.Value != nil && *priorityClass.Spec.Value < 0 {
		return matrix.ObjectHealthStatusBroken, nil
	}

	return matrix.ObjectHealthStatusHealthy, nil
}

func serviceHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func configMapHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func limitRangeHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func resourceQuotaHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func secretHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func csiDriverHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	// NOTE: requires cluster context (node info, pod status) to determine CSI driver health
	return matrix.ObjectHealthStatusUnknown, nil
}

func csiNodeHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	// NOTE: requires cluster context (node info, pod status) to determine CSI node health
	return matrix.ObjectHealthStatusUnknown, nil
}

func csiStorageCapacityHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func storageClassHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func mutatingAdmissionPolicyHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func mutatingAdmissionPolicyBindingHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func mutatingWebhookConfigurationHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func runtimeClassHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func validatingAdmissionPolicyHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func validatingWebhookConfigurationHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func clusterRoleBindingHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func roleBindingHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func serviceAccountHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func clusterTrustBundleHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func controllerRevisionHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func deviceTaintRuleHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func ipAddressHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func leaseCandidateHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func podTemplateHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func resourceClaimTemplateHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}

func volumeAttributesClassHealthStatus(kind string, raw json.RawMessage) (matrix.ObjectHealthStatus, error) {
	return matrix.ObjectHealthStatusUnknown, nil
}
