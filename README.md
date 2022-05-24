# pspmigrator: Migrate from PSP to PSA

pspmigrator is a tool to make it easier for K8s users to migrate from
PodSecurityPolicy (PSP) to PodSecurity Standards/Admission (PSA). The tool has
the following features:

- CLI tool to interactively migrate you from PSP to PSA by looking
  at the running pods in a namespace and suggesting a Pod Security Standard
- Detect if PSP object is potentially mutating Pods
- Detect if a Pod is being mutated by a PSP object

## Installation

```
go install github.com/samos123/pspmigrator/cmd/pspmigrator
```

Alternatively, you can download a release from [here](https://github.com/samos123/pspmigrator/releases/latest)

## Usage
Start the migration from PSP to PSA
```
pspmigrator migrate
# example output
Checking if any pods are being mutated by a PSP object
Suggest using baseline in namespace default
✔ enforce
Applied pod security level baseline on namespace default in enforce control mode
Review the labels by running `kubectl get ns default -o yaml`
There are no pods running in namespace empty. Skipping and going to the next one.
Done with migrating namespaces with pods to PSA
```

Full help docs available by running:
```
pspmigrator -h
pspmigrator is a tool to help migrate from PSP to PSA

Usage:
  pspmigrator [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  migrate     Interactive command to migrate from PSP to PSA
  mutating    Check if pods or PSP objects are mutating

Flags:
  -h, --help                help for pspmigrator
  -k, --kubeconfig string   (optional) absolute path to the kubeconfig file (default "/Users/stoelinga/.kube/config")

Use "pspmigrator [command] --help" for more information about a command.
```

Check if any pods are being mutated by a Pod Security Policy in your K8s cluster:
```
pspmigrator mutating pods
# example output
+----------------------------------------------------+-------------+---------+------------------------+
|                        NAME                        |  NAMESPACE  | MUTATED |          PSP           |
+----------------------------------------------------+-------------+---------+------------------------+
| nginx-nonpriv-66b6c48dd5-rl6jt                     | default     | true    | my-psp                 |
| event-exporter-gke-5479fd58c8-k8f5k                | kube-system | true    | gce.event-exporter     |
| fluentbit-gke-4hcg8                                | kube-system | true    | gce.fluentbit-gke      |
| gke-metadata-server-8bbrf                          | kube-system | false   | gce.privileged         |
+----------------------------------------------------+-------------+---------+------------------------+
```

Check if a specific pod called `my-pod` in namespace `my-namespace` is being
mutated by a Pod Security Policy:
```
pspmigrator mutating pod my-pod -n my-namespace
# example output
Pod nginx-nonpriv-66b6c48dd5-rl6jt is mutated by PSP my-psp: true, diff: [slice[0]: <nil pointer> != v1.SecurityContext]
PSP profile my-psp has the following mutating fields: [DefaultAddCapabilities] and annotations: []
```

## License
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
