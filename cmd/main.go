/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/samos123/pspmigrator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	//
	// Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	var mutatingCmd = &cobra.Command{
		Use:   "mutating",
		Short: "Check if pods or PSP objects are mutating",
		Long: `print is for printing anything back to the screen.
					  For many years people have printed back to the screen.`,
	}
	mutatingCmd.AddCommand(
		&cobra.Command{
			Use:   "pod",
			Short: "Check if a pod is being mutated by a PSP policy",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("check pod mutating")
			},
		},
	)
	var rootCmd = &cobra.Command{Use: "pspmigrator"}
	rootCmd.AddCommand(mutatingCmd)
	var kubeconfig string

	if home := homedir.HomeDir(); home != "" {
		rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k",
			filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "absolute path to the kubeconfig file")
	}
	rootCmd.Execute()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
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
