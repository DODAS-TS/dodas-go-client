## dodas app install minio

Install minio

### Synopsis

Install minio

```
dodas app install minio [flags]
```

### Examples

```
  k3sup app install minio
```

### Options

```
      --access-key string   Provide an access key to override the pre-generated value
      --distributed         Deploy Minio in Distributed Mode
  -h, --help                help for minio
      --namespace string    Kubernetes namespace for the application (default "default")
      --persistence         Enable persistence
      --secret-key string   Provide a secret key to override the pre-generated value
      --set stringArray     Use custom flags or override existing flags 
                            (example --set persistence.enabled=true)
      --update-repo         Update the helm repo (default true)
```

### Options inherited from parent commands

```
      --config string       DODAS config file (default is $HOME/.dodas.yaml)
      --kubeconfig string   Local path for your kubeconfig file (default "kubeconfig")
```

### SEE ALSO

* [dodas app install](dodas_app_install.md)	 - Install a Kubernetes app

###### Auto generated by spf13/cobra on 7-Feb-2020
