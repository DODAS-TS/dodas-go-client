## dodas app install cod

Install CachingOnDemand

### Synopsis

Install CachingOnDemand

```
dodas app install cod [flags]
```

### Examples

```
dodas app install cod
```

### Options

```
      --gsi-enabled              Enable GSI-VO auth
  -h, --help                     help for cod
      --nslaves int              Number of cache servers (default 2)
      --origin-host string       Origin server host (default "xrootd-cms.infn.it")
      --origin-port int          Origin server port (default 1094)
      --redirector-host string   Cache redirector host
      --redirector-port int      Cache redirector port (default 31213)
      --set stringArray          Use custom flags or override existing flags 
                                 (example --set persistence.enabled=true)
      --voms-file string         Provide a voms file for a vo (example --voms-file cms.txt=/tmp/voms/cms.txt)
      --x509-cert string         path to the certificate for cache server auth with remote
      --x509-key string          path to the certificate key for cache server auth with remote
```

### Options inherited from parent commands

```
      --config string       DODAS config file (default is $HOME/.dodas.yaml)
      --kubeconfig string   Kubernetes config file (default is /etc/kubernetes/admin.conf) (default "/etc/kubernetes/admin.conf")
```

### SEE ALSO

* [dodas app install](dodas_app_install.md)	 - Install a DODAS cluster with Kubernetes apps

###### Auto generated by spf13/cobra on 21-Mar-2020
