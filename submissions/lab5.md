# Lab 5 — Virtualization: QuickNotes in a Vagrant VM

## Environment

- Host operating system: Windows 10
- Hypervisor: Oracle VirtualBox 7.1.10
- Vagrant: 2.4.9
- Guest box: `ubuntu/jammy64` — Ubuntu 22.04 LTS
- Guest hostname: `quicknotes-vm`
- Go version: `go1.24.5 linux/amd64`
- VM resources: 2 vCPU and 1024 MB RAM
- Network: NAT with loopback-bound port forwarding
- Port mapping: `127.0.0.1:18080` to guest port `8080`

The implementation is defined in the root-level [Vagrantfile](../Vagrantfile).

---

## Task 1 — Vagrant Up and QuickNotes

### Successful clean provisioning

The VM was created from a clean state with:

```powershell
& "C:\Program Files\Vagrant\bin\vagrant.exe" up --provision
```

The command completed successfully with exit code `0`.

First ten meaningful Vagrant status lines from the successful run:

```text
Bringing machine 'default' up with 'virtualbox' provider...
==> default: Importing base box 'ubuntu/jammy64'...
==> default: Matching MAC address for NAT networking...
==> default: Checking if box 'ubuntu/jammy64' version '20241002.0.0' is up to date...
==> default: Setting the name of the VM: quicknotes-lab5
==> default: Clearing any previously set network interfaces...
==> default: Preparing network interfaces based on configuration...
==> default: Forwarding ports...
==> default: Booting VM...
==> default: Waiting for machine to boot. This may take a few minutes...
```

Later successful status lines included:

```text
==> default: Machine booted and ready!
==> default: Setting hostname...
==> default: Mounting shared folders...
==> default: Running provisioner: shell...
```

The complete first provisioning took:

```text
00:03:47.5465889
```

### Go verification

Command:

```powershell
vagrant ssh -c "go version"
```

Output:

```text
go version go1.24.5 linux/amd64
```

### QuickNotes service verification

Command:

```powershell
vagrant ssh -c "systemctl is-active quicknotes"
```

Output:

```text
active
```

### Guest health check

Command:

```powershell
vagrant ssh -c "curl -s http://127.0.0.1:8080/health"
```

Output:

```json
{"notes":4,"status":"ok"}
```

### Host health check through the forwarded port

Command:

```powershell
curl.exe -s http://127.0.0.1:18080/health
```

Output:

```json
{"notes":4,"status":"ok"}
```

This proves that QuickNotes was running inside the guest and that the loopback-bound host-to-guest port forwarding worked.

### Provisioning idempotency

The provisioner was executed again with:

```powershell
vagrant provision
```

The second run completed successfully with exit code `0`. It detected that Go 1.24.5 was already installed, rebuilt QuickNotes, restarted the systemd service, and preserved the working health endpoint.

### Design questions

#### a) Synced folders

I selected the VirtualBox shared-folder type:

```ruby
config.vm.synced_folder "./app", "/opt/quicknotes-src",
  type: "virtualbox"
```

This provides a live two-way mount between the Windows host directory and the VM and works directly with the native Windows Vagrant and VirtualBox installation. The trade-off is that it depends on compatible VirtualBox Guest Additions and can have slower metadata and small-file performance than a native Linux filesystem. An `rsync` folder may provide better guest-side filesystem performance, but it is normally one-way and requires another synchronization step when host files change.

#### b) NAT, bridged, and host-only networking

The VM uses VirtualBox's default NAT networking with an explicit forwarded port:

```ruby
config.vm.network "forwarded_port",
  guest: 8080,
  host: 18080,
  host_ip: "127.0.0.1"
```

NAT allows the guest to reach external networks while keeping it behind the host. Binding the forwarded port to `127.0.0.1` means only the local host can access QuickNotes. A bridged interface would give the VM an address on the physical LAN and could expose the service to other devices, which is unnecessary and less secure for this course exercise. Host-only networking would isolate communication to the host and guest but would not provide the same default outbound connectivity as NAT.

#### c) Provisioning method

I used Vagrant's shell provisioner because the required configuration is small and linear: install dependencies, download and verify Go, build QuickNotes, create a systemd service, and start it. The shell provisioner requires no additional configuration-management installation and keeps the complete setup visible inside the `Vagrantfile`. For a larger fleet or more complex configuration, Ansible would provide stronger abstractions, reusable roles, inventory management, and more structured idempotency.

#### d) Pinning Go to a point release

Go is pinned to version `1.24.5` rather than the moving `1.24` series. A specific version ensures that all users receive the same compiler, standard library, bug fixes, and build behaviour. The downloaded archive is also checked against a fixed SHA-256 value, so an unexpected or corrupted archive cannot silently change the environment. This makes the clean-clone provisioning process deterministic and reproducible.

---

## Task 2 — Snapshot, Break, and Restore

### 1. Save the clean snapshot

Command:

```powershell
vagrant snapshot save clean-quicknotes
```

Result:

```text
Snapshot save exit code: 0
Snapshot name: clean-quicknotes
```

### 2. Deliberately break the VM

The Go installation and its command links were removed:

```powershell
vagrant ssh -c "sudo rm -rf /usr/local/go; sudo rm -f /usr/local/bin/go /usr/local/bin/gofmt"
```

### 3. Prove that the VM was broken

Command:

```powershell
vagrant ssh -c "go version"
```

Output:

```text
bash: line 1: go: command not found
```

Exit code:

```text
127
```

### 4. Restore and time the snapshot

The restore was performed without rerunning the provisioner so that recovery came directly from the snapshot:

```powershell
Measure-Command {
    vagrant snapshot restore clean-quicknotes --no-provision
}
```

Result:

```text
Restore exit code: 0
Elapsed restore time: 00:00:33.1116211
```

### 5. Verify recovery

Go verification:

```powershell
vagrant ssh -c "go version"
```

Output:

```text
go version go1.24.5 linux/amd64
```

Service verification:

```powershell
vagrant ssh -c "systemctl is-active quicknotes"
```

Output:

```text
active
```

Host health verification:

```powershell
curl.exe -s http://127.0.0.1:18080/health
```

Output:

```json
{"notes":4,"status":"ok"}
```

### 6. Delete the temporary snapshot

Commands:

```powershell
vagrant snapshot delete clean-quicknotes
vagrant snapshot list
```

Result:

```text
Snapshot delete exit code: 0
No snapshots have been taken yet!
```

The snapshot was removed after the experiment to prevent unnecessary differencing-disk growth.

### Snapshot design questions

#### e) Why snapshots are not backups

A snapshot depends on the original VM storage and is normally stored on the same host disk as the VM. A physical-disk failure, host loss, serious filesystem corruption, ransomware incident, or accidental deletion of the complete VM can therefore destroy both the VM and its snapshots. A backup must be an independent copy stored separately and preferably off-site.

#### f) Copy-on-write disk usage

A VirtualBox snapshot does not immediately create a complete duplicate of the virtual disk. It preserves the existing state and writes later changes into a differencing disk. Ten snapshots therefore do not initially require ten full copies, but every snapshot creates another delta layer and total usage grows with the number of blocks changed after each snapshot. With sufficient write activity, a long chain can consume a large amount of disk space.

#### g) When snapshotting becomes an antipattern

Snapshotting becomes an antipattern when snapshots are retained as a permanent versioning or backup system, particularly when they form long chains. Reads and restores may need to traverse several dependent differencing disks, increasing I/O overhead, recovery complexity, storage consumption, and the consequences of corruption in the chain. Snapshots should be short-lived checkpoints that are restored or deleted after the immediate operation.

---

## Bonus — VM versus Docker Container

The measurements below were collected on the same host hardware.

### VM measurements

Cold boot command:

```powershell
vagrant halt
Measure-Command {
    vagrant up --no-provision
}
```

Cold boot result:

```text
00:00:39.6203004
```

Idle memory command:

```powershell
vagrant ssh -c "free -h"
```

Relevant output:

```text
               total        used        free      shared  buff/cache   available
Mem:           957Mi       193Mi       477Mi       1.0Mi       286Mi       616Mi
```

Process count command:

```powershell
vagrant ssh -c "ps -A --no-headers | wc -l"
```

Output:

```text
107
```

VM directory size:

```text
2,886,078,502 bytes
2.69 GiB
```

### Docker measurements

The same QuickNotes application was run using the pinned `golang:1.24.5-alpine` image:

```bash
docker run -d \
  --name quicknotes-lab5-container \
  -p 28080:8080 \
  -v "$PWD/app:/src" \
  -w /src \
  golang:1.24.5-alpine \
  sh -c 'go build -o /tmp/quicknotes && /tmp/quicknotes'
```

Health output:

```json
{"notes":4,"status":"ok"}
```

Cold-start measurement:

```bash
docker stop quicknotes-lab5-container
time docker start quicknotes-lab5-container
```

Output:

```text
real    0m0.406s
user    0m0.022s
sys     0m0.025s
```

Idle-memory output:

```text
Name=quicknotes-lab5-container Memory=8.543MiB / 5.788GiB CPU=0.00%
```

Container process table:

```text
UID    PID    PPID    C    STIME    TTY    TIME        CMD
root   1872   1849    1    11:29    ?      00:00:00    /tmp/quicknotes
```

Container process count:

```text
1
```

Image size:

```text
Repository=golang Tag=1.24.5-alpine Size=262MB
```

### Comparison

| Dimension | Vagrant VM | Docker container |
|---|---:|---:|
| Cold start | 39.620 s | 0.406 s |
| Idle RAM | 193 MiB | 8.543 MiB |
| On-disk size | 2.69 GiB | 262 MB |
| Process count | 107 | 1 |

The largest difference was the start time: the VM required approximately 39.6 seconds, while the container started in approximately 0.4 seconds. The container also used far less idle memory because it shared the host kernel and ran only the QuickNotes process instead of an entire guest operating system. A VM remains appropriate when a separate kernel, complete operating-system environment, stronger isolation boundary, or several system services are required. Containers are better suited to stateless services, CI jobs, and rapidly scaled application workloads. These measurements help explain why containers became dominant for stateless microservices between 2014 and 2020: they provided faster deployment, higher workload density, lower resource overhead, and simpler replacement of individual service instances.
