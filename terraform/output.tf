output "output-yaml" {
  value = replace(yamlencode(module.output), "/((?:^|\n)[\\s-]*)\"([\\w-]+)\":/", "$1$2:")
}