# Windows Operational Readiness


Define an operational readiness standard for Kubernetes clusters supporting Windows that certifies the readiness of Windows clusters for running production workloads

Related KEP: https://github.com/kubernetes/enhancements/tree/master/keps/sig-windows/2578-windows-conformance

## Build the project

#### Build the project with the default Kubernetes version (e.g. v1.24.0)
```
$ make build
```
#### Build the project with a specific Kubernetes version
```
$ make build KUBERNETES_HASH=<Kubernetes commit sha>
```

## Run the tests

To specify your windows cluster's readiness to run workflows:

- define the tests.yaml (you can use the "tests.yaml" in the repo as a template).

Tests categories can be passed in the flag `--category`, this allows users to pick a category of tests by run.
To run ALL tests do not pass the flag.

```
./op-readiness --provider=local --kubeconfig=<path-to-kubeconfig> --category=Core.Network --category=networkpolicy

Running Operational Readiness Test 1 / 10 : Ability to access Windows container IP by pod IP on Core.Network
...
Running Operational Readiness Test 2 / 10 : Ability to expose windows pods by creating the service ClusterIP on Core.Network
...
```

#### Run the Sonobuoy plugin

We support an OCI image and a Sonobuoy plugin, so the user don't need to compile the binary locally
by default the latest version of the E2E binary is builtin the image, if you need to add a custom file
just mount your local version in the plugin at `/app/e2e.test`.

Before running sonobuoy, taint the windows worker node. Sonobuoy pod should be scheduled on the control plane node:

```
kubectl taint node <windows-worker-node> sonobuoy:NoSchedule
```

To run the plugin with the default image:

```
make sonobuoy-plugin
```

To retrieve the sonobuoy result:

```
make sonobuoy-results
```

The result can be found in the `./sonobuoy-results` folder.

##### Set a particular category

It's possible to choose one or more categories of tests to run the plugin. In the following
example both `core` and `activedirectory` tests are being use by the plugin.

To allow all tests in the YAML don't use the `--category` flag.

```
spec:
  command:
    - /app/op-readiness
  args:
  - --category 
  - core
  - --category 
  - activedirectory
  - --test-file 
  - tests.yaml
````

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack channel](https://kubernetes.slack.com/messages/sig-windows) 
- [Mailing list](https://groups.google.com/g/kubernetes-sig-windows)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE