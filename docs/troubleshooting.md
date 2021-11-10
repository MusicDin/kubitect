# Troubleshooting - Common problems and possible solutions

Content:
1. [KVM/Libvirt errors](#kvmlibvirt-errors)
2. [HAProxy load balancer errors](#haproxy-load-balancer-errors)


## KVM/Libvirt errors


### -> Problem 1

#### Error:

*Error: virError(Code=38, Domain=7, Message='Failed to connect socket to '/var/run/libvirt/libvirt-sock': No such file or directory') on libvirt.tf line 1, in provider "libvirt”": 1: provider "libvirt" {...*

#### Explanation:
The problem may occur when libvirt is not started.

#### Solution:
Make sure that the `libvirt` service is running:
```sh
sudo systemctl status libvirtd
```

If the `libvirt` service is not running, start it:
```sh
sudo systemctl start libvirtd
```

*Optional:* Start the `libvirt` service automatically at boot time:
```bash
sudo systemctl enable libvirtd
```


### -> Problem 2

#### Error:

*Error: virError(Code=38, Domain=7, Message='Failed to connect socket to '/var/run/libvirt/libvirt-sock': Permission denied')*

#### Explanation and Possible Solution:
Make sure:
+ the `libvirtd` service is running and
+ the current user is in the `libvirt` and `kvm` groups.

#### Solution

If the `libvirtd` service is not running, start it:
```sh
sudo systemctl start libvirtd
```

Add the current user to the `libvirt` and `kvm` groups if needed:
```sh
# Add current user to groups
sudo usermod -aG libvirt,kvm `id -un`

# Verify groups are added
id -nG

# Reload user session
su - `id -un`
```


### -> Problem 3

#### Error:

*Error: Error creating libvirt domain: … Could not open '/tmp/terraform_libvirt_provider_images/image.qcow2': Permission denied')*

#### Explanation:
This problem can occur when you apply the Terraform plan.
+ Make sure that the directory exists.
+ Make sure that the directory of the file that is being denied has appropriate user permissions.
+ Optionally qemu security driver can be disabled. 

#### Solution:
Make sure the `security_driver` in `/etc/libvirt/qemu.conf` is set to `none` instead of `selinux`.
This line is commented out by default, so you should uncomment it if needed:
```bash
# /etc/libvirt/qemu.conf

...
security_driver = "none"
...
```

Do not forget to restart the `libvirt` service after making the changes:
```bash
sudo systemctl restart libvirtd
```


### -> Problem 4

#### Error:

*Error: Error defining libvirt domain: virError(Code=9, Domain=20, Message='operation failed: domain '**your-domain**' already exists with uuid '...')*

#### Explanation:
This problem can occur when you apply the Terraform plan.

#### Solution:
The resource you are trying to create already exists. 
Make sure you destroy the resource:
<pre>
virsh destroy <b>your-domain</b>
virsh undefine <b>your-domain</b>
</pre>

You can verify that the domain was successfully removed:
<pre>
virsh dominfo --domain <b>your-domain</b>
</pre>

If the domain was successfully removed, the output should look something like this:
<pre>
error: failed to get domain '<b>your-domain</b>'
</pre>


### -> Problem 5

#### Error:

*Error: Error creating libvirt volume: virError(Code=90, Domain=18, Message='storage volume '<b>your-volume</b>.qcow2' exists already')*

and / or

*Error:Error creating libvirt volume for cloudinit device <b>cloud-init</b>.iso: virError(Code=90, Domain=18, Message='storage volume '<b>cloud-init</b>.iso' exists already')*

#### Explanation:
This error can occur if you try to remove a faulty Terraform plan.

#### Solution:
Volumes created by Libvirt are still attached to the images, which prevents a new volume from being created with the same name.
Therefore, these volumes must be removed:
<pre>
virsh vol-delete <b>cloud-init</b>.iso --pool <b>your_resource_pool</b>

# and / or

virsh vol-delete <b>your-volume</b>.qcow2 --pool <b>your_resource_pool</b>
</pre>


### -> Problem 6

#### Error:

*Error: Error storage pool '**your-pool**' already exists*

#### Explanation:
Make sure you delete the created pool by first stopping it and then removing it.

#### Solution:
Remove the libvirt pool that was created during the Terraform process:
<pre>
virsh pool-destroy <b>your-pool</b> && virsh pool-undefine <b>your-pool</b>
</pre>


### -> Problem 7

#### Error:

*Error: Error **your-vm-name** already exists*

#### Explanation:
Your VM was stopped but not completely removed.

#### Solution:
Remove the stopped VM:
<pre>
virsh undefine <b>your-vm-name</b>
</pre>


### -> Problem 8

#### Error:

*Error: internal error: Failed to apply firewall rules /sbin/iptables -w --table filter --insert LIBVIRT_INP --in-interface virbr2 --protocol tcp --destination-port 67 --jump ACCEPT: iptables: No chain/target/match by that name.*

#### Explanation:
Libvirt was already running when Firewalld was installed.
Therefore, `libvirtd` service must be restarted to detect the changes.

#### Solution:
Restart `libvirtd` service:
```sh
sudo systemctl restart libvirtd
```


### -> Problem 9

#### Error:
*Error creating libvirt network: virError(Code=89, Domain=47, Message='COMMAND_FAILED: '/usr/sbin/iptables -w10 -w --table filter --insert LIBVIRT_INP --in-interface virbr1 --protocol tcp --destination-port 67 --jump ACCEPT' failed: iptables: No chain/target/match by that name.*

#### Explanation:
Libvirt was already running when iptables was set up.
Therefore, `libvirtd` service must be restarted to detect the changes.

#### Solution:
Restart `libvirtd` service:
```sh
sudo systemctl restart libvirtd
```


## HAProxy load balancer errors


### -> Problem 10

#### Error:

HAProxy returns a random *HTTP 503 (Bad gateway)* error.

#### Explanation:

More than one HAProxy processes are listening on the same port.

#### Solution 1:

For example, if an error is thrown when accessing port `80`, check which processes are listening on port `80` on the load balancer VM:
<pre>
netstat -lnput | grep <b>80</b>
</pre>
Output:
<pre>
Proto Recv-Q Send-Q Local Address           Foreign Address   State       PID/Program name
tcp        0      0 192.168.113.200:<b>80</b>      0.0.0.0:*         LISTEN      <b>1976</b>/haproxy
tcp        0      0 192.168.113.200:<b>80</b>      0.0.0.0:*         LISTEN      <b>1897</b>/haproxy
</pre>

If you see more than one process, kill the unnecessary process:
<pre>
kill <b>1976</b>
</pre>

*Note: You can kill all HAProxy processes and only one will be automatically recreated.*

#### Solution 2:

Check the HAProxy configuration file (`haproxy.cfg`) that it does not contain 2 frontends bound to the same port.
