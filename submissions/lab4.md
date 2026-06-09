# Lab 4 Submission — OS & Networking: Trace, Debug, and Read the Substrate

## Task 1 — Trace a Request End-to-End

### 1.1 Capture setup

QuickNotes was started locally from the `app/` directory using:

    go run .

A packet capture was started on loopback for TCP port 8080:

    sudo tcpdump -i lo -nn -s 0 -A 'tcp port 8080' -w lab4-trace.pcap

Then one request was sent:

    curl -v -X POST http://localhost:8080/notes \
      -H 'Content-Type: application/json' \
      -d '{"title":"trace me","body":"in flight"}'

The request returned:

    HTTP/1.1 201 Created

Response body:

    {"id":6,"title":"trace me","body":"in flight","created_at":"2026-06-09T20:21:16.071257766Z"}

### 1.2 Annotated packet trace

The decoded trace was saved in:

    lab4-trace.txt

TCP three-way handshake:

    23:21:16.065482 IP6 ::1.32966 > ::1.8080: Flags [S]
    23:21:16.066930 IP6 ::1.8080 > ::1.32966: Flags [S.]
    23:21:16.067406 IP6 ::1.32966 > ::1.8080: Flags [.]

This shows SYN, SYN/ACK, and ACK.

HTTP request:

    POST /notes HTTP/1.1
    Host: localhost:8080
    User-Agent: curl/8.5.0
    Accept: */*
    Content-Type: application/json
    Content-Length: 39

    {"title":"trace me","body":"in flight"}

HTTP response:

    HTTP/1.1 201 Created
    Content-Type: application/json
    Content-Length: 93

    {"id":6,"title":"trace me","body":"in flight","created_at":"2026-06-09T20:21:16.071257766Z"}

Connection close:

    23:21:16.090747 IP6 ::1.32966 > ::1.8080: Flags [F.]
    23:21:16.091016 IP6 ::1.8080 > ::1.32966: Flags [F.]
    23:21:16.091238 IP6 ::1.32966 > ::1.8080: Flags [.]

The connection closed with FIN/FIN/ACK.

### 1.3 Five debugging commands

#### 1. What is listening?

Command:

    ss -tlnp | grep :8080

Output:

    LISTEN 0 4096 *:8080 *:* users:(("quicknotes",pid=9136,fd=3))

Decision:

    QuickNotes is listening on TCP port 8080.

#### 2. What routes does the host have?

Command:

    ip route show

Output:

    default via 192.168.240.1 dev eth0 proto kernel
    192.168.240.0/20 dev eth0 proto kernel scope link src 192.168.254.234

Decision:

    The host has a default route through eth0 and a directly connected local subnet route.

#### 3. Can localhost be reached?

Command:

    mtr -rwc 5 localhost

Output:

    HOST: DESKTOP-U1R4GKD Loss% Snt Last Avg Best Wrst StDev
      1.|-- localhost 0.0% 5 0.1 0.3 0.1 1.4 0.6

Decision:

    Localhost is reachable with 0% packet loss over the loopback path.

#### 4. Does DNS work?

Command:

    dig +short example.com @1.1.1.1

Output:

    8.47.69.0
    8.6.112.0

Decision:

    DNS resolution works using Cloudflare resolver 1.1.1.1.

#### 5. Are there user service logs?

Command:

    journalctl --user -u quicknotes -n 20 || true

Output:

    -- No entries --

Decision:

    QuickNotes was run manually with `go run .`, not installed as a user systemd service, so there are no journald unit logs for `quicknotes`.

### 1.4 502 debugging reflection

If QuickNotes returned 502, I would debug outside-in. First I would check whether the service process is running, then whether it is listening on the expected port with `ss -tlnp`. Next I would test local reachability with `curl localhost:8080/health`. If the application is healthy locally, I would check the reverse proxy or load balancer path, firewall rules, and DNS resolution. This order separates application failure from transport, routing, firewall, and name-resolution problems.
