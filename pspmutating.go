package pspmigrator

import (
	"encoding/json"
	"fmt"
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

	b, err := json.Marshal(pspObj.Spec)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("%#v\n", pspObj.Spec)
	fmt.Printf("%s\n", b)
	var pspObjMap map[string]interface{}
	json.Unmarshal(b, &pspObjMap)

	mutatingFields := make(map[string]bool)
	mutatingFields["defaultAddCapabilities"] = true
	mutatingFields["requiredDropCapabilities"] = true
	mutatingFields["seLinux"] = true
	mutatingFields["runAsUser"] = true
	mutatingFields["runAsGroup"] = true
	mutatingFields["supplementalGroups"] = true
	mutatingFields["fsGroup"] = true
	mutatingFields["readOnlyRootFilesystem"] = true
	mutatingFields["defaultAllowPrivilegeEscalation"] = true
	mutatingFields["allowPrivilegeEscalation"] = true

	for field, value := range pspObjMap {
		rValue := reflect.ValueOf(value)
		// Some PSP fields aren't actually set, this captures those
		if rValue.Kind() == reflect.Map && rValue.Len() == 0 {
			continue
		}
		if _, ok := mutatingFields[field]; ok && reflect.ValueOf(value).Len() > 0 {
			fmt.Println(value, reflect.ValueOf(value).Kind())
			fields = append(fields, field)
		}
	}
	//	if len(pspObj.Spec.DefaultAddCapabilities) > 0 {
	//		fields = append(fields, "DefaultAddCapabilities")
	//	}
	//	if len(pspObj.Spec.RequiredDropCapabilities) > 0 {
	//		fields = append(fields, "RequiredDropCapabilities")
	//	}
	//	if (pspObj.Spec.SELinux != v1beta1.SELinuxStrategyOptions{}) {
	//		fields = append(fields, "SELinux")
	//	}
	//	if (pspObj.Spec.RunAsUser != v1beta1.RunAsUserStrategyOptions{}) {
	//		fields = append(fields, "RunAsUser")
	//	}
	//
	//	if pspObj.Spec.DefaultAllowPrivilegeEscalation != nil {
	//		fields = append(fields, "DefaultAllowPrivilegeEscalation")
	//	}

	//
	//	// Still need to filter fields that are nil, broken right now
	//	pspObjFields := reflect.ValueOf(&pspObj.Spec).Elem()
	//	for i := 0; i < pspObjFields.NumField(); i++ {
	//		field := pspObjFields.Type().Field(i)
	//		if _, ok := mutatingFields[field.Name]; ok && !isNil(field.Interface) {
	//			fields = append(fields, field.Name)
	//		}
	//	}
	//
	//	mutatingAnnotations := make(map[string]bool)
	//	mutatingAnnotations["seccomp.security.alpha.kubernetes.io/defaultProfileName"] = true
	//	mutatingAnnotations["apparmor.security.beta.kubernetes.io/defaultProfileName"] = true
	if len(fields) > 0 || len(annotations) > 0 {
		return true, fields, annotations
	}

	return false, fields, annotations
}
