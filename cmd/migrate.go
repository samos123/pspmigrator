package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/samos123/pspmigrator"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Interactive command to migrate from PSP to PSA ",
	Long: `The interactive command will help with setting a suggested a
	Suggested Pod Security Standard for each namespace. In addition, it also
	checks whether a PSP object is mutating pods in every namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		pods := GetPods()
		fmt.Println("Checking if any pods are being mutated by a PSP object")
		mutatedPods := make([]v1.Pod, 0)
		for _, pod := range pods.Items {
			mutated, _, err := pspmigrator.IsPodBeingMutatedByPSP(&pod, clientset)
			if err != nil {
				log.Println(err)
			}
			if mutated {
				mutatedPods = append(mutatedPods, pod)
			}
		}
		if len(mutatedPods) > 0 {
			fmt.Println("The table below shows the pods that were mutated by a PSP object")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Namespace", "PSP"})
			for _, pod := range mutatedPods {
				if pspName, ok := pod.ObjectMeta.Annotations["kubernetes.io/psp"]; ok {
					table.Append([]string{pod.Name, pod.Namespace, pspName})
				}
			}
			table.Render()
			pod := mutatedPods[0]
			fmt.Printf("There were %v pods mutated. Please modify the PodSpec such that PSP no longer needs to mutate your pod.\n", len(mutatedPods))
			fmt.Printf("You can for `pspmigrator mutating pod %v -n %v` to learn more why and how your pod is being mutated.", pod.Name, pod.Namespace)
			fmt.Printf("Please re-run the tool again after you've modified your PodSpecs.")
			os.Exit(1)
		}
		for _, namespace := range GetNamespaces().Items {
			suggestions := make(map[string]bool)
			for _, pod := range GetPodsByNamespace(namespace.Name).Items {
				level, err := pspmigrator.SuggestedPodSecurityStandard(&pod)
				if err != nil {
					fmt.Println("error occured checking the suggested pod security standard", err)
				}
				suggestions[string(level)] = true
			}
			var suggested string
			if suggestions["restricted"] {
				suggested = "restricted"
			}
			if suggestions["baseline"] {
				suggested = "baseline"
			}
			if suggestions["privileged"] {
				suggested = "privileged"
			}
			fmt.Printf("Suggest using %v in namespace %v\n", suggested, namespace.Name)
		}

	},
}
