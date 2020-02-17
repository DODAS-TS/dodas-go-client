package apps

// Voms ..
type Voms struct {
	Name    string
	Content string
}

// SlavesStruct ..
type SlavesStruct struct {
	SlaveNum   int
	SlaveMemGB int
	SlaveCPUs  int
	OsImage    string
}

// CoDParams ..
type CoDParams struct {
	Slaves         SlavesStruct
	GsiEnabled     bool
	ExternalIP     string
	VO             string
	VomsFile       Voms
	CacheCert      string
	CacheKey       string
	RedirectorHost string
	RedirectorPort int
	OriginHost     string
	OriginPort     int
}

// CoDTemplate ..
const CoDTemplate = `
tosca_definitions_version: tosca_simple_yaml_1_0

imports:
  - dodas_cod_types: https://raw.githubusercontent.com/dodas-ts/dodas-templates/master/tosca-types/dodas_custom_apps/cod_type.yml

description: TOSCA template for a complete CMS computing cluster on top of K8s orchestrator

topology_template:

  inputs:

    number_of_masters:
      type: integer
      default: 1

    num_cpus_master: 
      type: integer
      default: 2

    mem_size_master:
      type: string
      default: "4 GB"

    number_of_slaves:
      type: integer
      default: {{ .Slaves.SlaveNum }}

    num_cpus_slave: 
      type: integer
      default: 4

    mem_size_slave:
      type: string
      default: "8 GB"

    server_image:
      type: string
      #default: "ost://openstack.fisica.unipg.it/cb87a2ac-5469-4bd5-9cce-9682c798b4e4"
      default: "ost://horizon.cloud.cnaf.infn.it/3d993ab8-5d7b-4362-8fd6-af1391edca39"
      #default: "ost://cloud.recas.ba.infn.it/1113d7e8-fc5d-43b9-8d26-61906d89d479"

    helm_cod_values: 
      type: string
      default: |
        externalIp:
          enabled: true
          ips:
            - {{ .ExternalIP }}
        gsi:
          enabled: {{ .GsiEnabled }}
          {{- if .GsiEnabled }}
          vo: {{ .VO }}
          vomses:
          - filename: {{ .VomsFile.Name }}
            content: | {{ .VomsFile.Content | indent 14 }}
          caCert:
            cert: | {{ .CacheCert | indent 14 }}
            key: | {{ .CacheKey | indent 14 }}
          proxy: true
          {{- end }}

        cache:
          {{- if .RedirectorHost }}
          redirHost: {{ .RedirectorHost }}
          {{- end }}
          originHost: {{ .OriginHost }}
          originXrdPort: {{ .OriginPort }}
        redirector:
          {{- if .RedirectorHost }}
          enabled: false
          {{- end}}
          service:
            cms:
              port: {{ .RedirectorPort }}
        proxy:
          {{- if .RedirectorHost }}
          enabled: false
          {{- end}}

  node_templates:

    helm_cod:
      type: tosca.nodes.DODAS.HelmInstall.CachingOnDemand
      properties:
        repos:
          - { name: dodas, url: "https://dodas-ts.github.io/helm_charts" }
        name: "cachingondemand"
        chart: "dodas/cachingondemand"
        kubeconfig_path: /var/lib/rancher/k3s/server/cred/admin.kubeconfig
        values_file: { get_input: helm_cod_values }
        externalIp: { get_attribute: [ k3s_master_server , public_address, 0 ]  }
      requirements:
        - host: k3s_master_server
        - dependency: k3s_slave_server

    k3s_master:
      type: tosca.nodes.DODAS.k3s 
      properties:
        master_ip: { get_attribute: [ k3s_master_server, private_address, 0 ] }
        mode: server
      requirements:
        - host: k3s_master_server

    k3s_slave:
      type: tosca.nodes.DODAS.k3s 
      properties:
        master_ip: { get_attribute: [ k3s_master_server, private_address, 0 ] }
        mode: node 
      requirements:
        - host: k3s_slave_server
        - dependency: k3s_master

    k3s_master_server:
      type: tosca.nodes.indigo.Compute
      capabilities:
        endpoint:
          properties:
            # network_name: infn-farm.PUBLIC
            network_name: PUBLIC
            ports:
              kube_port:
                protocol: tcp
                source: 30443
              xrd_port:
                protocol: tcp
                source: 31394
        scalable:
          properties:
            count: { get_input: number_of_masters }
        host:
          properties:
            # instance_type:  m1.medium
            num_cpus: { get_input: num_cpus_master }
            mem_size: { get_input: mem_size_master } 
        os:
          properties:
            image: { get_input: server_image }

    k3s_slave_server:
      type: tosca.nodes.indigo.Compute
      capabilities:
        endpoint:
          properties:
            #network_name: test-net.PRIVATE
            network_name: PRIVATE
        scalable:
          properties:
            count: { get_input: number_of_slaves }
        host:
          properties:
            #instance_type:  m1.medium
            num_cpus: { get_input: num_cpus_slave }
            mem_size: { get_input: mem_size_slave } 
        os:
          properties:
            # image: "ost://openstack.fisica.unipg.it/d9a41aed-3ebf-42f9-992e-ef0078d3de95"
            image: { get_input: server_image }

  outputs:
    k8s_endpoint:
      value: { concat: [ 'https://', get_attribute: [ k3s_master_server, public_address, 0 ], ':30443' ] }
`
