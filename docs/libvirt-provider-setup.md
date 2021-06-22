# Setting up libvirt provider

For up to date installation guide and other informations check [libvirt provider's GitHub](https://github.com/dmacvicar/terraform-provider-libvirt).

Create the following directory:
<pre>
mkdir -p ~/.local/share/terraform/plugins/registry.terraform.io/dmacvicar/libvirt/<b>0.6.3</b>/linux_amd64
</pre>

Move to created directory:
<pre>
cd ~/.local/share/terraform/plugins/registry.terraform.io/dmacvicar/libvirt/<b>0.6.3</b>/linux_amd64
</pre>

Download libvirt provider `tar.gz` archive file. 
All possible assets can be found under their [releases](https://github.com/dmacvicar/terraform-provider-libvirt/releases).
Example for *Ubuntu 20.04*:
<pre>
wget -O terraform-provider-libvirt.tar.gz https://github.com/dmacvicar/terraform-provider-libvirt/releases/download/v<b>0.6.3</b>/terraform-provider-libvirt-<b>0.6.3</b>+git.1604843676.67f4f2aa.Ubuntu_20.04.amd64.tar.gz
</pre>

Unarchive downloaded file:
```
tar -xzf terraform-provider-libvirt.tar.gz
```

(Optionally) Remove `tar.gz` archive file:
```
rm terraform-provider-libvirt.tar.gz
```

That's it.

