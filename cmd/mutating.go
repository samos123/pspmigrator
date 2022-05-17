package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/samos123/pspmigrator"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var MutatingCmd = &cobra.Command{
	Use:   "mutating",
	Short: "Check if pods or PSP objects are mutating",
	Long: `print is for printing anything back to the screen.
					  For many years people have printed back to the screen.`,
}

func initMutating() {
	MutatingCmd.AddCommand(
		&cobra.Command{
			Use:   "pod",
			Short: "Check if a pod is being mutated by a PSP policy",
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) == 0 {
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
				}
				if len(args) == 1 {
					fmt.Println("get pod")
				}
			},
			Args: cobra.MaximumNArgs(1),
		},
	)
}
