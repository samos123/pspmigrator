package pspmigrator

import (
	"reflect"

	"k8s.io/api/policy/v1beta1"
)

func IsPSPMutating(pspObj *v1beta1.PodSecurityPolicy) (bool, []string, []string) {
	// check if associated PSP object is using any mutating fields
	// if yes then lookup ownerReferences and see if the field is actually mutating
	// if no continue check for next pod

	mutatingFields := make(map[string]bool)
	mutatingFields["DefaultAddCapabilities"] = true
	mutatingFields["RequiredDropCapabilities"] = true
	mutatingFields["SeLinux"] = true
	mutatingFields["RunAsUser"] = true
	mutatingFields["RunAsGroup"] = true
	mutatingFields["SupplementalGroups"] = true
	mutatingFields["FsGroup"] = true
	mutatingFields["ReadOnlyRootFilesystem"] = true
	mutatingFields["DefaultAllowPrivilegeEscalation"] = true
	mutatingFields["AllowPrivilegeEscalation"] = true

	fields := make([]string, 0)
	annotations := make([]string, 0)

	// Still need to filter fields that are nil, broken right now
	pspObjFields := reflect.ValueOf(&pspObj.Spec).Elem()
	for i := 0; i < pspObjFields.NumField(); i++ {
		field := pspObjFields.Type().Field(i).Name
		if _, ok := mutatingFields[field]; ok {
			fields = append(fields, field)
		}
	}

	mutatingAnnotations := make(map[string]bool)
	mutatingAnnotations["seccomp.security.alpha.kubernetes.io/defaultProfileName"] = true
	mutatingAnnotations["apparmor.security.beta.kubernetes.io/defaultProfileName"] = true
	if len(fields) > 0 || len(annotations) > 0 {
		return true, fields, annotations
	}

	return false, fields, annotations
}
