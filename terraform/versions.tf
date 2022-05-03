terraform {

  experiments      = [module_variable_optional_attrs]
  required_version = ">= 1.1.4"

  backend "local" {
    path = "../config/terraform/terraform.tfstate"
  }

  required_providers {
    libvirt = {
      source  = "dmacvicar/libvirt"
      version = "~> 0.6.12"
    }
  }
}
