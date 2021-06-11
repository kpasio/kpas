package tpl

import (
	"bytes"
	"text/template"
)

type HCloudClusterInfo struct {
	ClusterName    string
	HCloudApiToken string
	SshKeyName     string
}

func HCloudInventory(vars HCloudClusterInfo) []byte {
	tmpl, err := template.New("hcloud_inventory").Parse(`plugin: hcloud
strict: true
groups:
	first_master: (labels.first_master == 'true') and (labels.kpas_cluster_name == '{{.ClusterName}}')
	secondary_masters: (labels.secondary_master == 'true') and (labels.kpas_cluster_name == '{{.ClusterName}}')
	workers: (labels.worker == 'true') and (labels.kpas_cluster_name == '{{.ClusterName}}')
	masters_and_workers: ((labels.worker == 'true') or (labels.first_master == 'true') or (labels.secondary_master == 'true')) and (labels.kpas_cluster_name == '{{.ClusterName}}')
	etcd: ((labels.first_master == 'true') or (labels.secondary_master == 'true')) and (labels.kpas_cluster_name == '{{.ClusterName}}')
`)
	CheckIfError(err)

	var outputData bytes.Buffer
	err = tmpl.Execute(&outputData, vars)
	CheckIfError(err)

	return outputData.Bytes()
}

func HCloudProviderVault(vars HCloudClusterInfo) []byte {
	tmpl, err := template.New("hcloud_provider_vault").Parse(`vault_global_hetzner_api_token: {{ .HCloudApiToken }}
	`)
	CheckIfError(err)

	var outputData bytes.Buffer
	err = tmpl.Execute(&outputData, vars)
	CheckIfError(err)

	return outputData.Bytes()
}

func HCloudProviderValues(vars HCloudClusterInfo) []byte {
	tmpl, err := template.New("hcloud_provider_values").Parse(`## Hetzner Cloud Configuration
	global_hetzner:
		api_token: "{{"{{ vault_global_hetzner_api_token }}"}}"
		master_count: 1
		worker_count: 3
		master_type: cx11
		worker_type: cx41
		ssh_keys:
			- {{ .SshKeyName }}
	
	ansible_user: root
	ansible_ssh_common_args: '-o StrictHostKeyChecking=no'
	`)
	CheckIfError(err)

	var outputData bytes.Buffer
	err = tmpl.Execute(&outputData, vars)
	CheckIfError(err)

	return outputData.Bytes()
}
