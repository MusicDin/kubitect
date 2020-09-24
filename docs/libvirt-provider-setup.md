# Setting up libvirt provider

For up to date installation guide and other informations check [libvirt provider's GitHub](https://github.com/dmacvicar/terraform-provider-libvirt).

Create the following directory for Terraform plugins:
```bash
mkdir -p ~/.terraform.d/plugins
```

Download libvirt provider `.tar.gz` archive file 
(all possible assets can be found under their [releases](https://github.com/dmacvicar/terraform-provider-libvirt/releases)):
<pre>
wget -O terraform-provider-libvirt.tar.gz <b><i>https://github.com/dmacvicar/terraform-provider-libvirt/releases/download/v0.6.2/terraform-provider-libvirt-0.6.2+git.1585292411.8cbe9ad0.Ubuntu_18.04.amd64.tar.gz</i></b>
</pre>

Unarchive file in directory for Terraform plugins:
```bash
tar -xzf ~/.terraform.d/plugins/terraform-provider-libvirt.tar.gz -C ~/.terraform.d/plugins/
```

(*Optional*) Remove unnecessary archive file:
```bash
rm terraform-provider-libvirt.tar.gz
```
 
That's it.

