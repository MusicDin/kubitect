# Setting up libvirt provider

For up to date installation guide and other informations check [libvirt provider's GitHub](https://github.com/dmacvicar/terraform-provider-libvirt).

Create the following directory:
<pre>
mkdir -p ~/.local/share/terraform/plugins/registry.terraform.io/dmacvicar/libvirt/<b>0.6.2</b>/linux_amd64
</pre>

Move to created directory:
<pre>
cd ~/.local/share/terraform/plugins/registry.terraform.io/dmacvicar/libvirt/<b>0.6.2</b>/linux_amd64
</pre>

Download libvirt provider `.tar.gz` archive file
(all possible assets can be found under their [releases](https://github.com/dmacvicar/terraform-provider-libvirt/releases)):
<pre>
wget https://github.com/dmacvicar/terraform-provider-libvirt/releases/download/v<b>0.6.2</b>/terraform-provider-libvirt-<b>0.6.2</b>+git.1585292411.8cbe9ad0.Ubuntu_18.04.amd64.tar.gz
</pre>

Unarchive file:
<pre>
tar -xzf terraform-provider-libvirt-<b>0.6.2</b>+git.1585292411.8cbe9ad0.Ubuntu_18.04.amd64.tar.gz
</pre>

That's it.

