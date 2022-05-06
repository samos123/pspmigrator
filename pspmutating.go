package pspmigrator

import (
	"reflect"

	"k8s.io/api/policy/v1beta1"
)

// source https://stackoverflow.com/questions/13476349/check-for-nil-and-nil-interface-in-go
func isNil(c interface{}) bool {
	return c == nil || (reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil())
}

func IsPSPMutating(pspObj *v1beta1.PodSecurityPolicy) (bool, []string, []string) {
	// check if associated PSP object is using any mutating fields
	// if yes then lookup ownerReferences and see if the field is actually mutating
	// if no continue check for next pod
	fields := make([]string, 0)
	annotations := make([]string, 0)

	//	mutatingFields := make(map[string]bool)
	//	mutatingFields["defaultAddCapabilities"] = true
	//	mutatingFields["requiredDropCapabilities"] = true
	//	mutatingFields["seLinux"] = true
	//	mutatingFields["runAsUser"] = true
	//	mutatingFields["runAsGroup"] = true
	//	mutatingFields["supplementalGroups"] = true
	//	mutatingFields["fsGroup"] = true
	//	mutatingFields["readOnlyRootFilesystem"] = true
	//	mutatingFields["defaultAllowPrivilegeEscalation"] = true
	//	mutatingFields["allowPrivilegeEscalation"] = true
	//
	//	for field, value := range pspObjMap {
	//		rValue := reflect.ValueOf(value)
	//		// Some PSP fields aren't actually set, this captures those
	//		if rValue.Kind() == reflect.Map && rValue.Len() == 0 {
	//			continue
	//		}
	//		//		if (field == "seLinux" || field == "runAsUser") {
	//		//			continue
	//		//		}
	//		if _, ok := mutatingFields[field]; ok {
	//			fmt.Println(value, reflect.ValueOf(value).Kind())
	//			fields = append(fields, field)
	//		}
	//	}
	if len(pspObj.Spec.DefaultAddCapabilities) > 0 {
		fields = append(fields, "DefaultAddCapabilities")
	}
	if len(pspObj.Spec.RequiredDropCapabilities) > 0 {
		fields = append(fields, "RequiredDropCapabilities")
	}
	if pspObj.Spec.SELinux.Rule != v1beta1.SELinuxStrategyRunAsAny {
		fields = append(fields, "SELinux")
	}
	if pspObj.Spec.RunAsUser.Rule != v1beta1.RunAsUserStrategyRunAsAny {
		fields = append(fields, "RunAsUser")
	}
	if pspObj.Spec.RunAsGroup != nil && pspObj.Spec.RunAsGroup.Rule != v1beta1.RunAsGroupStrategyRunAsAny {
		fields = append(fields, "RunAsGroup")
	}
	if pspObj.Spec.SupplementalGroups.Rule != v1beta1.SupplementalGroupsStrategyRunAsAny {
		fields = append(fields, "SupplementalGroups")
	}
	if pspObj.Spec.FSGroup.Rule != v1beta1.FSGroupStrategyRunAsAny {
		fields = append(fields, "FSGroup")
	}
	if pspObj.Spec.ReadOnlyRootFilesystem != false {
		fields = append(fields, "ReadOnlyRootFilesystem")
	}
	if pspObj.Spec.DefaultAllowPrivilegeEscalation != nil {
		fields = append(fields, "DefaultAllowPrivilegeEscalation")
	}
	if pspObj.Spec.AllowPrivilegeEscalation != nil && *pspObj.Spec.AllowPrivilegeEscalation != true {
		fields = append(fields, "AllowPrivilegeEscalation")
	}

	mutatingAnnotations := make(map[string]bool)
	mutatingAnnotations["seccomp.security.alpha.kubernetes.io/defaultProfileName"] = true
	mutatingAnnotations["apparmor.security.beta.kubernetes.io/defaultProfileName"] = true

	for k, _ := range pspObj.Annotations {
		if _, ok := mutatingAnnotations[k]; ok {
			annotations = append(annotations, k)
		}
	}

	if len(fields) > 0 || len(annotations) > 0 {
		return true, fields, annotations
	}

	return false, fields, annotations
}
