# Self-hosted Kubernetes assets (kubeconfig, manifests)
module "bootkube" {
  source = "github.com/kinvolk/terraform-render-bootkube?ref=0f3156c73b846bf1bc5dd4681879048c37fb2f25"

  cluster_name          = var.cluster_name
  api_servers           = [format("%s.%s", var.cluster_name, var.dns_zone)]
  etcd_servers          = aws_route53_record.etcds.*.fqdn
  asset_dir             = var.asset_dir
  networking            = var.networking
  network_mtu           = var.network_mtu
  pod_cidr              = var.pod_cidr
  service_cidr          = var.service_cidr
  cluster_domain_suffix = var.cluster_domain_suffix
  enable_reporting      = var.enable_reporting
  enable_aggregation    = var.enable_aggregation

  certs_validity_period_hours = var.certs_validity_period_hours
}
