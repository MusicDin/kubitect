terraform {
  required_version = ">= 1.3.7"

  backend "local" {
    path = "../config/terraform/terraform.tfstate"
  }

  required_providers {
    libvirt = {
      source  = "dmacvicar/libvirt"
      version = "0.9.1"
    }
  }
}
