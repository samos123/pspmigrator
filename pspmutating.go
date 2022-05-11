package pspmigrator

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-test/deep"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetContainerSecurityContexts(podSpec v1.PodSpec) []*v1.SecurityContext {
	scs := make([]*v1.SecurityContext, 0)
	for _, c := range podSpec.Containers {
		scs = append(scs, c.SecurityContext)
	}
	fmt.Println("scs", scs)
	return scs
}

func GetPSPAnnotations(annotations map[string]string) map[string]string {
	pspAnnotations := make(map[string]string)
	for ann, val := range annotations {
		if strings.Contains(ann, "seccomp.security") || strings.Contains(ann, "apparmor.security") {
			pspAnnotations[ann] = val
		}
	}
	return pspAnnotations
}

func IsPodBeingMutatedByPSP(pod *v1.Pod, clientset *kubernetes.Clientset) (bool, error) {
	if len(pod.ObjectMeta.OwnerReferences) > 0 {
		var owner metav1.OwnerReference
		for _, reference := range pod.ObjectMeta.OwnerReferences {
			if reference.Controller != nil && *reference.Controller == true {
				owner = reference
				break
			}
		}
		var parentPod v1.PodTemplateSpec
		if owner.Kind == "ReplicaSet" {
			rs, err := clientset.AppsV1().ReplicaSets(pod.Namespace).Get(context.TODO(), owner.Name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			parentPod = rs.Spec.Template
		}
		if diff := deep.Equal(GetContainerSecurityContexts(parentPod.Spec), GetContainerSecurityContexts(pod.Spec)); diff != nil {
			return true, nil
		}
		if diff := deep.Equal(parentPod.Spec.SecurityContext, pod.Spec.SecurityContext); diff != nil {
			return true, nil
		}
		if diff := deep.Equal(GetPSPAnnotations(parentPod.ObjectMeta.Annotations), GetPSPAnnotations(pod.ObjectMeta.Annotations)); diff != nil {
			return true, nil
		}
	}
	return false, nil
}

func IsPodBeingMutatedByPSPOld(pod *v1.Pod, clientset *kubernetes.Clientset) (bool, error) {
	// check if associated PSP object is using any mutating fields
	// if yes then lookup ownerReferences and see if the field is actually mutating
	// if no continue check for next pod
	if pspName, ok := pod.ObjectMeta.Annotations["kubernetes.io/psp"]; ok {
		pspObj, err := clientset.PolicyV1beta1().PodSecurityPolicies().Get(context.TODO(), pspName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return false, fmt.Errorf("PodSecurityPolicy %s not found: %w\n", pspName, err)
		} else if err != nil {
			return false, fmt.Errorf("Error getting PodSecurityPolicy %s: %w\n",
				pspName, err)
		} else if err != nil {
			mutating, fields, annotations := IsPSPMutating(pspObj)
			if mutating == false {
				return false, nil
			}
			fmt.Println(fields, annotations)
			// Lookup ownerReferences and compare pod spec with owner pod spec
			if len(pod.ObjectMeta.OwnerReferences) > 0 {
				var owner metav1.OwnerReference
				for _, reference := range pod.ObjectMeta.OwnerReferences {
					if reference.Controller != nil && *reference.Controller == true {
						owner = reference
						break
					}
				}
				if owner.Kind == "ReplicaSet" {
					rs, err := clientset.AppsV1().ReplicaSets(pod.Namespace).Get(context.TODO(), owner.Name, metav1.GetOptions{})
					if err != nil {
						return false, err
					}
					parentPodSpec := rs.Spec.Template.Spec
					log.Println(parentPodSpec)
				}
				return true, nil
			}
		}
	}
	return false, nil
}

// IsPSPMutating checks wheter a PodSecurityPolicy is potentially mutating
// pods. It returns true if one of the fields or annotations used in the
// PodSecurityPolicy is suspected to be mutating pods. The field or annotations
// that are suspected to be mutating are returned as well.
func IsPSPMutating(pspObj *v1beta1.PodSecurityPolicy) (mutating bool, fields, annotations []string) {
	fields = make([]string, 0)
	annotations = make([]string, 0)

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
	if pspObj.Spec.RunAsGroup != nil && pspObj.Spec.RunAsGroup.Rule == v1beta1.RunAsGroupStrategyMustRunAs {
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
