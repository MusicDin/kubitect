<h1 align="center"> Troubleshooting</h1>

!!! question "Is your issue not listed here?"

    If the troubleshooting page is missing an error you encountered, please report it on GitHub by [opening an issue](https://github.com/MusicDin/kubitect/issues/new).
    By doing so, you will help improve the project and help others find the solution to the same problem faster.

## General errors

### Virtualenv not found

=== ":material-close-thick: Error"

    !!! failure "Error"

        Output: /bin/sh: 1: virtualenv: not found

        /bin/sh: 2: ansible-playbook: not found

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        The error indicates that the `virtualenv` is not installed.

=== ":material-check-bold: Solution"

    !!! success "Solution"

        There are many ways to install `virtualenv`.
        For all installation options you can refere to their official documentation - [Virtualenv installation](https://virtualenv.pypa.io/en/latest/installation.html).

        For example, virtualenv can be installed using `pip`.

        First install pip.
        ```sh
        sudo apt install python3-pip
        ```

        Then install virtualenv using pip3.
        ```sh
        pip3 install virtualenv
        ```

## KVM/Libvirt errors

### Failed to connect socket (No such file or directory)

=== ":material-close-thick: Error"

    !!! failure "Error"

        Error: virError(Code=38, Domain=7, Message='Failed to connect socket to '/var/run/libvirt/libvirt-sock': No such file or directory')

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        The problem may occur when libvirt is not started.

=== ":material-check-bold: Solution"

    !!! success "Solution"

        Make sure that the `libvirt` service is running:
        ```sh
        sudo systemctl status libvirtd
        ```

        If the `libvirt` service is not running, start it:
        ```sh
        sudo systemctl start libvirtd
        ```

        *Optional:* Start the `libvirt` service automatically at boot time:
        ```sh
        sudo systemctl enable libvirtd
        ```

### Failed to connect socket (Permission denied)

=== ":material-close-thick: Error"

    !!! failure "Error"

        Error: virError(Code=38, Domain=7, Message='Failed to connect socket to '/var/run/libvirt/libvirt-sock': Permission denied')

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        The error indicates that either the `libvirtd` service is not running or the current user is not in the `libvirt` (or `kvm`) group.

=== ":material-check-bold: Solution"

    !!! success "Solution"

        If the `libvirtd` service is not running, start it:
        ```sh
        sudo systemctl start libvirtd
        ```

        Add the current user to the `libvirt` and `kvm` groups if needed:
        ```sh
        # Add current user to groups
        sudo adduser $USER libvirt
        sudo adduser $USER kvm

        # Verify groups are added
        id -nG

        # Reload user session
        ```

### Error creating libvirt domain

=== ":material-close-thick: Error"

    !!! failure "Error"

        Error: Error creating libvirt domain: … Could not open '/tmp/terraform_libvirt_provider_images/image.qcow2': Permission denied')

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        The error indicates that the file cannot be created in the specified location due to missing permissions.

        + Make sure the directory exists.
        + Make sure the directory of the file that is being denied has appropriate user permissions.
        + Optionally qemu security driver can be disabled. 

=== ":material-check-bold: Solution"

    !!! success "Solution"

        Make sure the `security_driver` in `/etc/libvirt/qemu.conf` is set to `none` instead of `selinux`.
        This line is commented out by default, so you should uncomment it if needed:
        ```sh
        # /etc/libvirt/qemu.conf

        ...
        security_driver = "none"
        ...
        ```

        Do not forget to restart the `libvirt` service after making the changes:
        ```sh
        sudo systemctl restart libvirtd
        ```

### Libvirt domain already exists

=== ":material-close-thick: Error"

    !!! failure "Error"

        Error: Error defining libvirt domain: virError(Code=9, Domain=20, Message='operation failed: domain '**your-domain**' already exists with uuid '...')

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        The error indicates that the libvirt domain (virtual machine) already exists.

=== ":material-check-bold: Solution"

    !!! success "Solution"

        The resource you are trying to create already exists. 
        Make sure you destroy the resource:
        ```
        virsh destroy your-domain
        virsh undefine your-domain
        ```

        You can verify that the domain was successfully removed:
        ```
        virsh dominfo --domain your-domain
        ```

        If the domain was successfully removed, the output should look something like this:

        <code>
        error: failed to get domain '<b>your-domain</b>'
        </code>

### Libvirt volume already exists

=== ":material-close-thick: Error"

    !!! failure "Error"

        Error: Error creating libvirt volume: virError(Code=90, Domain=18, Message='storage volume '<b>your-volume</b>.qcow2' exists already')

        and / or

        Error:Error creating libvirt volume for cloudinit device <b>cloud-init</b>.iso: virError(Code=90, Domain=18, Message='storage volume '<b>cloud-init</b>.iso' exists already')

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        The error indicates that the specified volume already exists.

=== ":material-check-bold: Solution"

    !!! success "Solution"

        Volumes created by Libvirt are still attached to the images, which prevents a new volume from being created with the same name.
        Therefore, these volumes must be removed:

        <code>
        virsh vol-delete <b>cloud-init</b>.iso --pool <b>your_resource_pool</b>
        </code>

        and / or

        <code>
        virsh vol-delete <b>your-volume</b>.qcow2 --pool <b>your_resource_pool</b>
        </code>

### Libvirt storage pool already exists

=== ":material-close-thick: Error"

    !!! failure "Error"

        Error: Error storage pool '**your-pool**' already exists

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        The error indicates that the libvirt storage pool already exists.

=== ":material-check-bold: Solution"

    !!! success "Solution"

        Remove the existing libvirt storage pool.

        <code>
        virsh pool-destroy <b>your-pool</b> && virsh pool-undefine <b>your-pool</b>
        </code>

### Failed to apply firewall rules

=== ":material-close-thick: Error"

    !!! failure "Error"

        Error: internal error: Failed to apply firewall rules /sbin/iptables -w --table filter --insert LIBVIRT_INP --in-interface virbr2 --protocol tcp --destination-port 67 --jump ACCEPT: iptables: No chain/target/match by that name.

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        Libvirt was already running when firewall (usually FirewallD) was started/installed.
        Therefore, `libvirtd` service must be restarted to detect the changes.

=== ":material-check-bold: Solution"

    !!! success "Solution"

        Restart the `libvirtd` service:
        ```sh
        sudo systemctl restart libvirtd
        ```

### Failed to remove storage pool

=== ":material-close-thick: Error"

    !!! failure "Error"

        Error: error deleting storage pool: failed to remove pool '/var/lib/libvirt/images/local-k8s-cluster-main-resource-pool': Directory not empty

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        The pool cannot be deleted because there are still some volumes in the pool.
        Therefore, the volumes should be removed before the pool can be deleted.

=== ":material-check-bold: Solution"

    !!! success "Solution"

        1. Make sure the pool is running.
        ```sh
        virsh pool-start --pool local-k8s-cluster-main-resource-pool
        ```

        2. List volumes in the pool.
        ```sh
        virsh vol-list --pool local-k8s-cluster-main-resource-pool

        #  Name         Path
        # -------------------------------------------------------------------------------------
        #  base_volume  /var/lib/libvirt/images/local-k8s-cluster-main-resource-pool/base_volume
        ```

        3. Delete listed volumes from the pool.
        ```sh
        virsh vol-delete --pool local-k8s-cluster-main-resource-pool --vol base_volume
        ```

        4. Destroy and undefine the pool.
        ```sh
        virsh pool-destroy --pool local-k8s-cluster-main-resource-pool
        virsh pool-undefine --pool local-k8s-cluster-main-resource-pool
        ```


## HAProxy load balancer errors

### Random HAProxy (503) bad gateway


=== ":material-close-thick: Error"

    !!! failure "Error"

        HAProxy returns a random *HTTP 503 (Bad gateway)* error.

=== ":material-information-outline: Explanation"

    !!! info "Explanation"

        More than one HAProxy processes are listening on the same port.

=== ":material-check-bold: Solution"

    !!! success "Solution 1"

        For example, if an error is thrown when accessing port `80`, check which processes are listening on port `80` on the load balancer VM:
        ```sh
        netstat -lnput | grep 80

        # Proto Recv-Q Send-Q Local Address           Foreign Address   State       PID/Program name
        # tcp        0      0 192.168.113.200:80      0.0.0.0:*         LISTEN      1976/haproxy
        # tcp        0      0 192.168.113.200:80      0.0.0.0:*         LISTEN      1897/haproxy
        ``` 

        If you see more than one process, kill the unnecessary process:
        ```sh
        kill 1976
        ```

        *Note: You can kill all HAProxy processes and only one will be automatically recreated.*

    !!! success "Solution 2"

        Check the HAProxy configuration file (`config/haproxy/haproxy.cfg`) that it does not contain 2 frontends bound to the same port.

