/*
Copyright 2018 The Kubernetes Authors.

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

package main


import (
	"encoding/json"

	"k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const (
	addFirstLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels", "value": {"added-label": "yes"}}
     ]`
	addAdditionalLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels/added-label", "value": "yes" }
     ]`
	updateLabelPatch string = `[
         { "op": "replace", "path": "/metadata/labels/added-label", "value": "yes" }
     ]`
)

// Add a label {"added-label": "yes"} to the object
func AddLabel(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("calling add-label")
	obj := struct {
		metav1.ObjectMeta `json:"metadata,omitempty"`
	}{}
	// raw中的内容
	//`{"kind":"Pod","apiVersion":"v1","metadata":{"name":"my-curl","namespace":"default","creationTimestamp":null,"labels":{"app":"my-curl"},"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"labels\":{\"app\":\"my-curl\"},\"name\":\"my-curl\",\"namespace\":\"default\"},\"spec\":{\"containers\":[{\"image\":\"hysyeah/my-curl:v1\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"my-curl\",\"ports\":[{\"containerPort\":8080}]}]}}\n"},"managedFields":[{"manager":"kubectl-client-side-apply","operation":"Update","apiVersion":"v1","time":"2023-04-22T12:29:48Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/last-applied-configuration":{}},"f:labels":{".":{},"f:app":{}}},"f:spec":{"f:containers":{"k:{\"name\":\"my-curl\"}":{".":{},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":8080,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:enableServiceLinks":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}}}]},"spec":{"volumes":[{"name":"kube-api-access-29gqq","projected":{"sources":[{"serviceAccountToken":{"expirationSeconds":3607,"path":"token"}},{"configMap":{"name":"kube-root-ca.crt","items":[{"key":"ca.crt","path":"ca.crt"}]}},{"downwardAPI":{"items":[{"path":"namespace","fieldRef":{"apiVersion":"v1","fieldPath":"metadata.namespace"}}]}}],"defaultMode":420}}],"containers":[{"name":"my-curl","image":"hysyeah/my-curl:v1","ports":[{"containerPort":8080,"protocol":"TCP"}],"resources":{},"volumeMounts":[{"name":"kube-api-access-29gqq","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}],"priority":0,"enableServiceLinks":true,"preemptionPolicy":"PreemptLowerPriority"},"status":{}}`
	raw := ar.Request.Object.Raw
	// 将raw中数据反序列化到obj
	err := json.Unmarshal(raw, &obj)
	if err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}

	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true

	pt := v1.PatchTypeJSONPatch
	// 获取added-label标签
	labelValue, hasLabel := obj.ObjectMeta.Labels["added-label"]
	switch {
	case len(obj.ObjectMeta.Labels) == 0:
		// 添加label
		reviewResponse.Patch = []byte(addFirstLabelPatch)
		reviewResponse.PatchType = &pt
	case !hasLabel:
		reviewResponse.Patch = []byte(addAdditionalLabelPatch)
		reviewResponse.PatchType = &pt
	case labelValue != "yes":
		reviewResponse.Patch = []byte(updateLabelPatch)
		reviewResponse.PatchType = &pt
	default:
		// already set
	}
	return &reviewResponse
}
