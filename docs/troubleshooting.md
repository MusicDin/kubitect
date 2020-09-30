# Troubleshooting - Common problems and possible solutions

### -> Problem 1

#### Error:

*Error: virError(Code=38, Domain=7, Message='Failed to connect socket to '/var/run/libvirt/libvirt-sock': No such file or directory') on libvirt.tf line 1, in provider "libvirt”": 1: provider "libvirt" {...*

#### Explanation
The problem can arise when libvirt is not started.

#### Solution:
Verify that `libvirt` service is running:
```bash
sudo systemctl status libvirtd
```

If `libvirt` service is not running, you need to start it:
```bash
sudo systemctl start libvirtd
```

*Optional:* Automatically start `libvirt` service at boot time:
```bash
sudo systemctl enable libvirtd
```

### -> Problem 2

#### Error

*Error: virError(Code=38, Domain=7, Message='Failed to connect socket to '/var/run/libvirt/libvirt-sock': Permission denied')*

#### Explanation and Possible Solution
Check following:
+ Is libvirt running?
+ Is your user in the libvirt group? 
+ If on a virtual machine and you just installed libvirt for the first time, make sure to restart the machine and try again.

### -> Problem 3

#### Error:

*Error: Error creating libvirt domain: … Could not open '/tmp/terraform_libvirt_provider_images/image.qcow2': Permission denied')*

#### Explanation
This problem can occur when applying the Terraform plan on Libvirt provider.
+ Is the directory existing?
+ Make sure the directory of the file that is denied has user permissions.

#### Solution:
Make sure the `security_driver` in `/etc/libvirt/qemu.conf` is set to `none` instead of `selinux`.
This line is by default commented, so if needed uncomment it:
```bash
# /etc/libvirt/qemu.conf

...
security_driver = "none"
...
```

Don't forget to restart `libvirt` service after making changes:
```bash
sudo systemctl restart libvirtd
```

### -> Problem 4

#### Error:

*Error: Error defining libvirt domain: virError(Code=9, Domain=20, Message='operation failed: domain '**your-domain**' already exists with uuid '...')*

#### Explanation
This problem can occur when applying the Terraform plan on Libvirt provider.

#### Solution:
Resource that you are trying to create, already exists. Make sure to destroy the resource:
<pre>
virsh destroy <b>your-domain</b>
virsh undefine <b>your-domain</b>
</pre>

You can verify that the domain is successfully removed with:
<pre>
virsh dominfo --domain <b>your-domain</b>
</pre>

If domain has been removed successfully, output should be something like:
<pre>
error: failed to get domain '<b>your-domain</b>'
</pre>

### -> Problem 5

#### Error:

*Error: Error creating libvirt volume: virError(Code=90, Domain=18, Message='storage volume '<b>your-volume</b>.qcow2' exists already')*

and / or

*Error:Error creating libvirt volume for cloudinit device <b>cloud-init</b>.iso: virError(Code=90, Domain=18, Message='storage volume '<b>cloud-init</b>.iso' exists already')*

#### Explanation
This error can occur when trying to remove a faulty Terraform plan.

#### Solution:
Volumes created by Libvirt are still attached to the images, which prevents a new volume from being applied with the same name. 
Therefore, removal of these volumes is required:
<pre>
virsh vol-delete <b>cloud-init</b>.iso --pool <b>your_resource_pool</b>

# and / or

virsh vol-delete <b>your-volume</b>.qcow2 --pool <b>your_resource_pool</b>
</pre>

### -> Problem 6

#### Error:

*Error: Error storage pool '**your-pool**' already exists*

#### Explanation
Make sure you delete the created pool as well, first by halting it and removing it afterwards.

#### Solution:
Remove the libvirt pool that was created during the Terraform process:
<pre>
virsh pool-destroy <b>your-pool</b> && virsh pool-undefine <b>your-pool</b>
</pre>

### -> Problem 7

#### Error:

*Error: Error **your-vm-name** already exists*

#### Explanation
Your VM has been halted but not removed completely.

#### Solution:
Remove the running VM:
<pre>
virsh undefine <b>your-vm-name</b>
</pre>


