package pspmigrator

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"k8s.io/api/policy/v1beta1"
)

type PSPOptions struct {
	Annotations            map[string]string
	DefaultAddCapabilities []string
	RunAsGroup             map[string]string
}

func GeneratePSPObject(options PSPOptions) v1beta1.PodSecurityPolicy {
	annotations := []byte("{}")
	if len(options.Annotations) > 0 {
		annotations, _ = json.Marshal(options.Annotations)
	}
	defaultAddCapabilities, err := json.Marshal(options.DefaultAddCapabilities)
	if err != nil {
		log.Panic(err)
	}
	pspObjJson := fmt.Sprintf(`{
		"metadata":{
			"name":"my-psp",
			"annotations": %s
		},
		"spec":{
			"defaultAddCapabilities":%s,
			"volumes":["*"],
			"seLinux":{"rule":"RunAsAny"},
			"runAsUser":{"rule":"RunAsAny"},
			"supplementalGroups":{"rule":"RunAsAny"},
			"fsGroup":{"rule":"RunAsAny"},
			"allowPrivilegeEscalation":true}
	}`, annotations, defaultAddCapabilities)
	fmt.Println(pspObjJson)
	var pspObj v1beta1.PodSecurityPolicy
	if err := json.Unmarshal([]byte(pspObjJson), &pspObj); err != nil {
		log.Panic(err)
	}
	return pspObj
}

func TestIsPSPNotMutating(t *testing.T) {
	pspObj := GeneratePSPObject(PSPOptions{})
	yes, fields, annotations := IsPSPMutating(&pspObj)
	fmt.Println("Mutating, fields, annotations:", yes, fields, annotations)
	if yes == true {
		t.Error("Mutating should be false but was true")
	}
	if len(fields) > 0 {
		t.Errorf("Expected fields to be empty, but got: %s", fields)
	}
	if len(annotations) > 0 {
		t.Errorf("Expected annoations to be empty, but got: %s", annotations)
	}
}

func TestIsPSPMutatingDefaultAddCapabilitiesOnly(t *testing.T) {
	pspObj := GeneratePSPObject(PSPOptions{DefaultAddCapabilities: []string{"CHOWN"}})
	yes, fields, annotations := IsPSPMutating(&pspObj)
	fmt.Println("Mutating, fields, annotations:", yes, fields, annotations)
	if yes == false {
		t.Error("Mutating should be true but was false")
	}
	if len(fields) != 1 {
		t.Errorf("Only DefaultAddCapabilities should have been reported as mutating, but got %s", fields)
	}
	if fields[0] != "DefaultAddCapabilities" {
		t.Errorf("Expected DefaultAddCapabilities to be mutating but got %s", fields[0])
	}
}

func TestIsPSPMutatingAnnotation(t *testing.T) {
	pspObj := GeneratePSPObject(PSPOptions{Annotations: map[string]string{"seccomp.security.alpha.kubernetes.io/defaultProfileName": "a"}})
	yes, fields, annotations := IsPSPMutating(&pspObj)
	fmt.Println("Mutating, fields, annotations:", yes, fields, annotations)
	if yes == false {
		t.Error("Mutating should be true but was false")
	}
}

func TestIsPSPMutatingPodTrue(t *testing.T) {

}
