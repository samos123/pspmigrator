package pspmigrator

import (
	"fmt"

	v1 "k8s.io/api/core/v1"

	"k8s.io/pod-security-admission/api"
	psaApi "k8s.io/pod-security-admission/api"
	"k8s.io/pod-security-admission/policy"
)

func SuggestedPodSecurityStandard(pod *v1.Pod) (psaApi.Level, error) {
	evaluator, err := policy.NewEvaluator(policy.DefaultChecks())
	if err != nil {
		return "", err
	}
	apiVersion, err := api.ParseVersion("latest")
	if err != nil {
		return "", err
	}
	for _, level := range []string{"restricted", "baseline"} {
		apiLevel, err := psaApi.ParseLevel(level)
		if err != nil {
			return "", err
		}
		result := policy.AggregateCheckResults(evaluator.EvaluatePod(
			psaApi.LevelVersion{Level: apiLevel, Version: apiVersion}, &pod.ObjectMeta, &pod.Spec))

		if result.Allowed {
			return apiLevel, nil
		}
		forbiddenReasons := make([]string, len(result.ForbiddenReasons))
		for i := range result.ForbiddenReasons {
			forbiddenReasons = append(forbiddenReasons,
				fmt.Sprintf("%s: %s, ", result.ForbiddenReasons[i], result.ForbiddenDetails[i]))
		}
	}
	return api.LevelPrivileged, nil
}
