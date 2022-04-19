resource "aws_prometheus_workspace" "monitoring" {
  alias = local.environment_name
}

data "aws_iam_policy_document" "monitoring" {
  statement {
    effect = "Allow"

    actions = [
      "aps:RemoteWrite",
      "aps:QueryMetrics",
      "aps:GetSeries",
      "aps:GetLabels",
      "aps:GetMetricMetadata"
    ]

    resources = [aws_prometheus_workspace.monitoring.arn]
  }
}

resource "aws_iam_policy" "monitoring" {
  name   = "${local.environment_name}_monitoring"
  policy = data.aws_iam_policy_document.monitoring.json
}

module "monitoring_role" {
  source  = "registry.terraform.io/terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version = "4.17.1"

  create_role = true

  role_name    = "monitoring"
  role_path = local.iam_path
  provider_url = module.eks.oidc_provider

  role_policy_arns = [
    aws_iam_policy.monitoring.arn,
  ]

  oidc_fully_qualified_subjects = ["system:serviceaccount:monitoring:prometheus-k8s"]
}

resource "aws_security_group" "vpc_tls" {
  name        = "${local.environment_name}_vpc_tls"
  description = "Allow TLS inbound traffic"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "TLS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }
}

module "endpoints" {
  source  = "registry.terraform.io/terraform-aws-modules/vpc/aws//modules/vpc-endpoints"
  version = "3.14.0"

  vpc_id             = module.vpc.vpc_id
  security_group_ids = [aws_security_group.vpc_tls.id]
  subnet_ids         = module.vpc.private_subnets

  endpoints = {
    aps-workspaces = {
      service             = "aps-workspaces"
      vpc_endpoint_type   = "Interface"
      private_dns_enabled = true
    }
  }
}
