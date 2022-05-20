package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/samos123/pspmigrator"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var MutatingCmd = &cobra.Command{
	Use:   "mutating",
	Short: "Check if pods or PSP objects are mutating",
	Long: `print is for printing anything back to the screen.
					  For many years people have printed back to the screen.`,
}

func initMutating() {
	var namespace string
	podCmd := cobra.Command{
		Use:   "pod [name of pod]",
		Short: "Check if a pod is being mutated by a PSP policy",
		Run: func(cmd *cobra.Command, args []string) {
			// Examples for error handling:
			// - Use helper functions like e.g. errors.IsNotFound()
			// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
			pod := args[0]
			podObj, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				fmt.Printf("Error getting pod %s in namespace %s: %v\n",
					pod, namespace, statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				mutated, diff, err := pspmigrator.IsPodBeingMutatedByPSP(podObj, clientset)
				if err != nil {
					log.Println(err)
				}
				if pspName, ok := podObj.ObjectMeta.Annotations["kubernetes.io/psp"]; ok {
					fmt.Printf("Pod %v is mutated by PSP %v: %v, diff: %v\n", podObj.Name, pspName, mutated, diff)
				}
			}
		},
		Args: cobra.ExactArgs(1),
	}

	podCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "K8s namespace (required)")
	podCmd.MarkFlagRequired("namespace")

	podsCmd := cobra.Command{
		Use:   "pods",
		Short: "Check all pods across all namespaces in a cluster are being mutated by a PSP policy",
		Run: func(cmd *cobra.Command, args []string) {
			pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
			for _, pod := range pods.Items {
				if pspName, ok := pod.ObjectMeta.Annotations["kubernetes.io/psp"]; ok {
					mutated, diff, err := pspmigrator.IsPodBeingMutatedByPSP(&pod, clientset)
					if err != nil {
						log.Println(err)
					}
					fmt.Printf("Pod %v is mutated by PSP %v: %v, diff: %v\n", pod.Name, pspName, mutated, diff)
				}
			}
		},
		Args: cobra.NoArgs,
	}

	MutatingCmd.AddCommand(&podCmd)
	MutatingCmd.AddCommand(&podsCmd)
}
