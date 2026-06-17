Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/jammy64"
  config.vm.hostname = "quicknotes-vm"
  config.vm.boot_timeout = 600

  # NAT port forwarding: host-only access to QuickNotes.
  config.vm.network "forwarded_port",
    guest: 8080,
    host: 18080,
    host_ip: "127.0.0.1",
    auto_correct: false

  # Disable the default project-root mount and sync only the application.
  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.synced_folder "./app", "/opt/quicknotes-src",
    type: "virtualbox"

  config.vm.provider "virtualbox" do |vb|
    vb.name = "quicknotes-lab5"
    vb.memory = 1024
    vb.cpus = 2
    vb.gui = false
    vb.customize ["modifyvm", :id, "--uartmode1", "disconnected"]
  end

  config.vm.provision "shell", privileged: true, inline: <<-SHELL
    set -euxo pipefail

    GO_VERSION="1.24.5"
    GO_ARCHIVE="go${GO_VERSION}.linux-amd64.tar.gz"
    GO_SHA256="10ad9e86233e74c0f6590fe5426895de6bf388964210eac34a6d83f38918ecdc"

    export DEBIAN_FRONTEND=noninteractive
    apt-get update
    apt-get install -y ca-certificates curl

    CURRENT_GO_VERSION=""
    if [ -x /usr/local/go/bin/go ]; then
      CURRENT_GO_VERSION=$(/usr/local/go/bin/go version | awk '{print $3}')
    fi

    if [ "${CURRENT_GO_VERSION}" != "go${GO_VERSION}" ]; then
      curl -fsSLo "/tmp/${GO_ARCHIVE}" \
        "https://go.dev/dl/${GO_ARCHIVE}"
      echo "${GO_SHA256}  /tmp/${GO_ARCHIVE}" | sha256sum -c -

      rm -rf /usr/local/go
      tar -C /usr/local -xzf "/tmp/${GO_ARCHIVE}"
      rm -f "/tmp/${GO_ARCHIVE}"
    fi

    ln -sf /usr/local/go/bin/go /usr/local/bin/go
    ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt

    install -d -o vagrant -g vagrant /var/lib/quicknotes

    if [ ! -f /var/lib/quicknotes/notes.json ]; then
      cp /opt/quicknotes-src/seed.json /var/lib/quicknotes/notes.json
      chown vagrant:vagrant /var/lib/quicknotes/notes.json
    fi

    cd /opt/quicknotes-src
    /usr/local/go/bin/go build -o /usr/local/bin/quicknotes .

    cat > /etc/systemd/system/quicknotes.service <<'UNIT'
    [Unit]
    Description=QuickNotes service
    After=network-online.target

    [Service]
    Type=simple
    User=vagrant
    Group=vagrant
    WorkingDirectory=/opt/quicknotes-src
    Environment=ADDR=:8080
    Environment=DATA_PATH=/var/lib/quicknotes/notes.json
    ExecStart=/usr/local/bin/quicknotes
    Restart=on-failure
    RestartSec=2

    [Install]
    WantedBy=multi-user.target
UNIT

    systemctl daemon-reload
    systemctl enable quicknotes
    systemctl restart quicknotes
  SHELL
end
