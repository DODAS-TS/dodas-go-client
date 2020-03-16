# DODAS client

<p align="center">
<img src="https://github.com/DODAS-TS/dodas-templates/raw/master/logo.png" width="200" height="200" />
</p>

[![Build Status](https://travis-ci.org/DODAS-TS/dodas-go-client.svg?branch=master)](https://travis-ci.org/DODAS-TS/dodas-go-client)

## Installation and usage

Download the binary from the latest release on [github](https://github.com/DODAS-TS/dodas-go-client/releases). For instance:

```bash
wget https://github.com/DODAS-TS/dodas-go-client/releases/download/v1.3.0/dodas.zip
unzip dodas.zip
cp dodas /usr/local/bin
```

In alternative you can also run the dodas command inside the client container `dodasts/dodas-client:v1.3.0`.

> **CLI autocomplete**
>
> Autocompletion for bash and zsh is supported
>
> - **bash** add the following line to ~/.bashrc: `. <(dodas autocomplete)`
> - **zsh** add the following line to ~/.zshrc: `source <(dodas zsh-autocomplete)`

You can find now a template for creating your client configuration file in [config/client_config.yaml](https://raw.githubusercontent.com/DODAS-TS/dodas-go-client/master/config/client_config.yaml). Note that by default the client will look for `$HOME/.dodas.yaml`.

Now you are ready to go. For instance you can validate a tosca template like this:

```bash
dodas validate --template tests/tosca/valid_template.yml
```

> **Tip:** you can find supported tosca templates for applications and k8s deployments on the [DODAS template repo](https://github.com/DODAS-TS/dodas-templates)

or you can create a cluster through the InfrastructureManager configured in your configuration file:

```bash
dodas create --config my_client_conf.yaml my_template.yaml
```

To list the Infrastructure ID of all your deployments:

```bash
dodas list infIDs
```

You can eventually login into a vm in created cluster that has a public IP address with:

```bash
dodas login <infID> <vmID>
# e.g. dodas login cb585e5c-33b6-11ea-8776-0242ac150003 0
```

## Quick start

Your deployments will be created and managed by the [InfrastructureManager](https://www.grycap.upv.es/im/index.php)(IM).
To start playing with DODAS please refer to this two quick start guides on the official site:

- using the **[community instance of IM](https://dodas-ts.github.io/dodas-templates/quick-start-community/)** (required free registration for evaluation purpose [here](https://dodas-iam.cloud.cnaf.infn.it))
- a **[standalone setup](https://dodas-ts.github.io/dodas-templates/quick-start/)** where IM will be deployed in a docker container and used with the client

> **N.B** All of the pre-compiled templates provided by DODAS use the helm charts defined and documented [here](https://github.com/DODAS-TS/helm_charts/tree/master/stable).
>
> Therefore **all the available applications can be installed as they are on top of any k8s instance with [Helm](https://helm.sh/)**

## Building from source

To compile on a linux machine (go version that supports `go modules` is required for building from source: e.g. >= v1.12):

```bash
make build
```

while to compile with Docker:

```bash
make docker-build
```

It's also possible to cross compile for windows and macOS with:

```bash
make windows-build
make macos-build
```

## Contributing

If you want to contribute:

1. create a branch
2. upload your changes
3. create a pull request

Thanks!

## Acknowledgement

**This work is co-funded by the EOSC-hub project (Horizon 2020) under Grant number 777536.**                          

![EU logo](https://github.com/DODAS-TS/dodas-templates/raw/master/docs/img/eu-logo.jpeg)                              
![EOSC hub logo](https://github.com/DODAS-TS/dodas-templates/raw/master/docs/img/eosc-hub-web.png)

## Contact us

DODAS Team provides two support channels, email and Slack channel.

- **mailing list**: send a message to the following list dodas-support@lists.infn.it
- **slack channel**: join us on [Slack Channel](https://dodas-infn.slack.com/archives/CAJ6VG71A)
