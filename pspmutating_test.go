package pspmigrator

import (
	"fmt"
	"testing"

	"k8s.io/api/policy/v1beta1"
)

func TestIsPSPMutatingEmptyPolicy(t *testing.T) {
	pspObj := v1beta1.PodSecurityPolicy{}
	yes, fields, annotations := IsPSPMutating(&pspObj)
	fmt.Println("Mutating, fields, annotations:", yes, fields, annotations)
	if yes == true {
		t.Error("Mutating should be false but was true")
	}
}
