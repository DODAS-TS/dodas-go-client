## dodas

A self-sufficient client for DODAS deployments

### Synopsis

A self-sufficient client for DODAS deployments.
Default configuration file searched in $HOME/.dodas.yaml

Usage examples:
"""
# CREATE A CLUSTER FROM TEMPLATE
dodas create --template my_tosca_template.yml

# VALIDATE TOSCA TEMPLATE
dodas validate --template my_tosca_template.yml
"""

```
dodas [flags]
```

### Options

```
      --config string       DODAS config file (default is $HOME/.dodas.yaml)
  -h, --help                help for dodas
      --kubeconfig string   Kubernetes config file (default is /etc/kubernetes/admin.conf) (default "/etc/kubernetes/admin.conf")
  -v, --version             DODAS client version
```

### SEE ALSO

* [dodas app](dodas_app.md)	 - Install Kubernetes cluster and apps from helm charts or YAML files
* [dodas autocomplete](dodas_autocomplete.md)	 - Generate script for bash autocomplete
* [dodas create](dodas_create.md)	 - Create a cluster from a TOSCA template
* [dodas destroy](dodas_destroy.md)	 - Destroy infrastructure with this InfID
* [dodas get](dodas_get.md)	 - Wrapper command for get operations
* [dodas iam](dodas_iam.md)	 - Wrapper command for IAM interaction
* [dodas list](dodas_list.md)	 - Wrapper function for list operations
* [dodas login](dodas_login.md)	 - ssh login into a deployed vm
* [dodas reboot](dodas_reboot.md)	 - reboot a vm in the cluster
* [dodas reconfig](dodas_reconfig.md)	 - restart cluster configuration
* [dodas refresh](dodas_refresh.md)	 - Not implemented yet
* [dodas update](dodas_update.md)	 - Update the number of vms to satisfy the new template
* [dodas validate](dodas_validate.md)	 - Validate your tosca template
* [dodas version](dodas_version.md)	 - Client version
* [dodas zsh-autocomplete](dodas_zsh-autocomplete.md)	 - Generate script for zsh autocomplete

###### Auto generated by spf13/cobra on 21-Mar-2020
