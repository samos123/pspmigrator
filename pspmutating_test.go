package pspmigrator

import (
	"encoding/json"
	"fmt"
	"testing"

	"k8s.io/api/policy/v1beta1"
)

func GeneratePSPObject() v1beta1.PodSecurityPolicy {
	pspObjJson := []byte(`{"metadata":{"name":"my-psp","uid":"ba0359c9-bc30-407b-97de-7705fe4259c0","resourceVersion":"2302205","creationTimestamp":"2022-05-03T16:28:07Z","annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"policy/v1beta1\",\"kind\":\"PodSecurityPolicy\",\"metadata\":{\"annotations\":{},\"name\":\"my-psp\"},\"spec\":{\"defaultAddCapabilities\":[\"CHOWN\"],\"fsGroup\":{\"rule\":\"RunAsAny\"},\"privileged\":false,\"runAsUser\":{\"rule\":\"RunAsAny\"},\"seLinux\":{\"rule\":\"RunAsAny\"},\"supplementalGroups\":{\"rule\":\"RunAsAny\"},\"volumes\":[\"*\"]}}\n"},"managedFields":[{"manager":"kubectl-client-side-apply","operation":"Update","apiVersion":"policy/v1beta1","time":"2022-05-03T23:31:19Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/last-applied-configuration":{}}},"f:spec":{"f:allowPrivilegeEscalation":{},"f:defaultAddCapabilities":{},"f:fsGroup":{"f:rule":{}},"f:runAsUser":{"f:rule":{}},"f:seLinux":{"f:rule":{}},"f:supplementalGroups":{"f:rule":{}},"f:volumes":{}}}}]},"spec":{"defaultAddCapabilities":["CHOWN"],"volumes":["*"],"seLinux":{"rule":"RunAsAny"},"runAsUser":{"rule":"RunAsAny"},"supplementalGroups":{"rule":"RunAsAny"},"fsGroup":{"rule":"RunAsAny"},"allowPrivilegeEscalation":true}}`)
	var pspObj v1beta1.PodSecurityPolicy
	json.Unmarshal(pspObjJson, &pspObj)
	return pspObj
}

func TestIsPSPNotMutating(t *testing.T) {
	pspSpecJson := []byte(`{"defaultAddCapabilities":[],"volumes":["*"],"seLinux":{"rule":"RunAsAny"},"runAsUser":{"rule":"RunAsAny"},"supplementalGroups":{"rule":"RunAsAny"},"fsGroup":{"rule":"RunAsAny"},"allowPrivilegeEscalation":true}`)

	var pspSpec v1beta1.PodSecurityPolicySpec
	json.Unmarshal(pspSpecJson, &pspSpec)
	pspObj := v1beta1.PodSecurityPolicy{Spec: pspSpec}

	yes, fields, annotations := IsPSPMutating(&pspObj)
	fmt.Println("Mutating, fields, annotations:", yes, fields, annotations)
	if yes == true {
		t.Error("Mutating should be false but was true")
	}
}

func TestIsPSPMutatingDefaultAddCapabilitiesOnly(t *testing.T) {

	//	pspSpecJson := []byte(`{"defaultAddCapabilities":["CHOWN"],"volumes":["*"],"seLinux":{"rule":"RunAsAny"},"runAsUser":{"rule":"RunAsAny"},"supplementalGroups":{"rule":"RunAsAny"},"fsGroup":{"rule":"RunAsAny"},"allowPrivilegeEscalation":true}`)
	//
	//	var pspSpec v1beta1.PodSecurityPolicySpec
	//	json.Unmarshal(pspSpecJson, &pspSpec)
	pspObj := GeneratePSPObject()

	yes, fields, annotations := IsPSPMutating(&pspObj)
	fmt.Println("Mutating, fields, annotations:", yes, fields, annotations)
	if yes == false {
		t.Error("Mutating should be true but was false")
	}

}
