package tpl

import (
	"bytes"
	"text/template"
)

func GitIgnoreTemplate() []byte {
	return []byte(`
**/kubeconfs/*.yml
`)
}

type GenericClusterInfo struct {
	ClusterName                 string
	ClusterProvider             string
	ClusterVariant              string
	CloudFlareDnsEnabled        bool
	CloudFlareDnsApiToken       string
	CloudFlareDnsZoneIdentifier string
	BaseDomain                  string
}

func (ClusterInfo GenericClusterInfo) RandomHex() string {
	str, _ := randomHex(16)
	return str
}

func SharedVault(vars GenericClusterInfo) []byte {
	tmpl, err := template.New("shared_vault").Parse(`vault_global_base_domain: {{ .BaseDomain }}
vault_global_basic_auth_password: {{ .RandomHex }}
vault_global_grafana_password: {{ .RandomHex }}
vault_global_drone_bearer_token: {{ .RandomHex }}
vault_global_longhorn_aws_backup_access_key_id: 
vault_global_longhorn_aws_backup_access_key: 
vault_global_gitea_admin_password: {{ .RandomHex }}
{{ if .CloudFlareDnsEnabled }}
vault_global_cloudflare_api_token: {{ .CloudflareDnsApiToken }}
vault_global_cloudflare_zone_identifier: {{ .CloudflareDnsZoneIdentifier }}
{{ end }}
	`)
	CheckIfError(err)

	var outputData bytes.Buffer
	err = tmpl.Execute(&outputData, vars)
	CheckIfError(err)

	return outputData.Bytes()
}

func SharedValues(vars GenericClusterInfo) []byte {
	tmpl, err := template.New("shared_values").Parse(`	global_cluster:
	state: present
	name: {{ .ClusterName }}
	kpas_config_version: 2
	kpas_provider: {{ .ClusterProvider }}
	kpas_variant: {{ .ClusterVariant }}

{{ if eq .ClusterVariant "k3s" }}
k3s:
  shared_secret: <%= SecureRandom.hex %>
  multimaster: <%= options.fetch(:k3s_multimaster, false) == 'true' %>
{{ end }}

global_master_kubeconf_base_dir: <%= path_facts.kubeconf_directory %>
global_master_kubeconf: "{{"{{ global_master_kubeconf_base_dir }}"}}/master.yml"

# Configure apps to pre-configure on cluster
global_install:
  monitoring_apps: true
  registry: true
  ci: true
  longhorn: true

<% if options[:dns_cloudflare] %>
global_cloudflare:
  enable: true
  api_token: "{{"{{vault_global_cloudflare_api_token}}" }}"
  zone_identifier: "{{"{{vault_global_cloudflare_zone_identifier}}"}}"
<% else %>
global_cloudflare:
  enable: false
<% end %>

# Configure ingress
global_ingress:
  type: nginx
  base_domain: "{{"{{vault_global_base_domain}}"}}"
  ## SSL Certificate Issuance
  # values:
  #   letsencrypt-staging
  #   letsencrypt-production
  certificate_issuer: <%= options[:letsencrypt_production] ? 'letsencrypt-production' : 'letsencrypt-staging' %>
  tls: true
  force_tls: true
  basic_auth_username: admin
  basic_auth_password:  "{{"{{vault_global_basic_auth_password}}"}}"

global_grafana:
  username: admin
  password: "{{"{{vault_global_grafana_password}}"}}"

global_drone:
  bearer_token: "{{"{{vault_global_drone_bearer_token}}"}}"

# Longhorn configuration
global_longhorn:
  default_replica_count: 3
  backups: 
    enabled: false
    backup_schedule: '[{"name":"snap", "task":"snapshot", "cron":"*/5 * * * *", "retain":6},{"name":"backup", "task":"backup", "cron":"*/10 * * * *", "retain":5}]'
    aws_backup_access_key_id: "{{"{{vault_global_longhorn_aws_backup_access_key_id}}"}}"
    aws_backup_secret_access_key: "{{"{{vault_global_longhorn_aws_backup_access_key}}"}}"
    backup_target: 

global_gitea:
  admin_username: administrator
  admin_email: user@example.com
  admin_password: "{{"{{vault_global_gitea_admin_password}}"}}"
  ssh_port: 30010
  http_domain: "gitea.{{"{{global_ingress.base_domain}}"}}"
  git_domain: "gitea.{{"{{global_ingress.base_domain}}"}}"

# Namespaces used by core cluster apps, we create these in advance
# So that in the event of having to restore longhorn from a backup,
# we can just batch re-create from previous PVC's and have all volumes
# waiting. This will only work for stateful sets, not for simple 
# deployments.
cluster_app_namespaces:
  - monitoring
  - cluster-registry

# These are namespaces for deploying workloads to, for each of these
# a service account will be created which only has access to that
# namespace. A suitable kubeconfig file will be created in 
# kubeconf/{{"{{global_cluster.name}}"}}/namespace.yaml
#
# These namespaces all have a network policy applied which
# blocks all ingress from other namespaces except for required
# access to monitoring and logging and cluster ingress. Ingress and Egress
# within the namespace is unrestricted.
user_app_namespaces:
  - myapp
	`)
	CheckIfError(err)

	var outputData bytes.Buffer
	err = tmpl.Execute(&outputData, vars)
	CheckIfError(err)

	return outputData.Bytes()
}
