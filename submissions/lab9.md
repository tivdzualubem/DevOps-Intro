# Lab 9 — DevSecOps: Trivy, ZAP, and `govulncheck`

## Submission links

- [Course repository PR #1298](https://github.com/inno-devops-labs/DevOps-Intro/pull/1298)
- [Fork PR #9](https://github.com/tivdzualubem/DevOps-Intro/pull/9)
- Red CI demonstration commit: [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c)
- Green recovery commit: [053da8c](https://github.com/tivdzualubem/DevOps-Intro/commit/053da8ce5c4c6883d210b1f9332912a010446669)

## Tool versions and scope

| Tool | Pinned version | Purpose |
| --- | --- | --- |
| Trivy | `aquasec/trivy:0.59.1` | Image, filesystem, configuration, and CycloneDX SBOM |
| OWASP ZAP | `ghcr.io/zaproxy/zaproxy:2.16.1` | Passive baseline scan only |
| govulncheck | `golang.org/x/vuln/cmd/govulncheck@v1.5.0` | Reachable Go vulnerability CI gate |
| Final image builder | `golang:1.26.4-alpine3.24` | Patched production build toolchain |

The ZAP baseline was run against `http://localhost:8080`. No active
scan was performed.

---

# Task 1 — Trivy image, filesystem, configuration, and SBOM

## 1. Scan artifacts

- [Image scan before remediation](../artifacts/lab9/trivy/image-before.txt)
- [Image scan before remediation — JSON](../artifacts/lab9/trivy/image-before.json)
- [Image scan after remediation](../artifacts/lab9/trivy/image-after.txt)
- [Image scan after remediation — JSON](../artifacts/lab9/trivy/image-after.json)
- [Filesystem scan](../artifacts/lab9/trivy/filesystem.txt)
- [Filesystem scan — JSON](../artifacts/lab9/trivy/filesystem.json)
- [Configuration scan](../artifacts/lab9/trivy/config.txt)
- [Configuration scan — JSON](../artifacts/lab9/trivy/config.json)
- [Pre-remediation CycloneDX SBOM](../artifacts/lab9/trivy/sbom-before.cdx.json)
- [Final CycloneDX SBOM](../artifacts/lab9/trivy/sbom.cdx.json)

## 2. Results summary

| Scan | HIGH | CRITICAL | Decision |
| --- | ---: | ---: | --- |
| Image before remediation | 20 occurrences / 10 unique CVEs | 0 | All fixed |
| Image after remediation | 0 | 0 | Clean |
| Repository filesystem | 0 | 0 | No triage required |
| Repository configuration | 0 | 0 | No triage required |

The image findings were Go standard-library vulnerabilities embedded in
both `/quicknotes` and `/healthcheck`. The distroless Debian runtime layer
itself reported zero HIGH/CRITICAL vulnerabilities.

The final SBOM records the Go standard library changing from
`v1.24.13` to `v1.26.4`.

## 3. Top of image scan output

```text
2026-07-01T14:15:43Z	INFO	[vuln] Vulnerability scanning is enabled
2026-07-01T14:15:43Z	INFO	[secret] Secret scanning is enabled
2026-07-01T14:15:43Z	INFO	[secret] If your scanning is slow, please try '--scanners vuln' to disable secret scanning
2026-07-01T14:15:43Z	INFO	[secret] Please see also https://aquasecurity.github.io/trivy/v0.59/docs/scanner/secret#recommendation for faster secret detection
2026-07-01T14:15:43Z	INFO	Detected OS	family="debian" version="12.14"
2026-07-01T14:15:43Z	INFO	[debian] Detecting vulnerabilities...	os_version="12" pkg_num=4
2026-07-01T14:15:43Z	INFO	Number of language-specific files	num=2
2026-07-01T14:15:43Z	INFO	[gobinary] Detecting vulnerabilities...
2026-07-01T14:15:43Z	WARN	Using severities from other vendors for some vulnerabilities. Read https://aquasecurity.github.io/trivy/v0.59/docs/scanner/vulnerability#severity-selection for details.

quicknotes:lab6 (debian 12.14)
==============================
Total: 0 (HIGH: 0, CRITICAL: 0)


healthcheck (gobinary)
======================
Total: 10 (HIGH: 10, CRITICAL: 0)

┌─────────┬────────────────┬──────────┬────────┬───────────────────┬─────────────────┬──────────────────────────────────────────────────────────────┐
│ Library │ Vulnerability  │ Severity │ Status │ Installed Version │  Fixed Version  │                            Title                             │
├─────────┼────────────────┼──────────┼────────┼───────────────────┼─────────────────┼──────────────────────────────────────────────────────────────┤
│ stdlib  │ CVE-2026-25679 │ HIGH     │ fixed  │ v1.24.13          │ 1.25.8, 1.26.1  │ net/url: Incorrect parsing of IPv6 host literals in net/url  │
│         │                │          │        │                   │                 │ https://avd.aquasec.com/nvd/cve-2026-25679                   │
│         ├────────────────┤          │        │                   ├─────────────────┼──────────────────────────────────────────────────────────────┤
│         │ CVE-2026-27145 │          │        │                   │ 1.25.11, 1.26.4 │ crypto/x509: golang: golang crypto/x509: Denial of Service   │
│         │                │          │        │                   │                 │ via excessive processing of DNS...                           │
│         │                │          │        │                   │                 │ https://avd.aquasec.com/nvd/cve-2026-27145                   │
│         ├────────────────┤          │        │                   ├─────────────────┼──────────────────────────────────────────────────────────────┤
│         │ CVE-2026-32280 │          │        │                   │ 1.25.9, 1.26.2  │ crypto/x509: crypto/tls: golang: Go: Denial of Service       │
│         │                │          │        │                   │                 │ vulnerability in certificate chain building...               │
│         │                │          │        │                   │                 │ https://avd.aquasec.com/nvd/cve-2026-32280                   │
│         ├────────────────┤          │        │                   │                 ├──────────────────────────────────────────────────────────────┤
│         │ CVE-2026-32281 │          │        │                   │                 │ crypto/x509: golang: Go crypto/x509: Denial of Service via   │
│         │                │          │        │                   │                 │ inefficient certificate chain validation...                  │
```

## 4. Top of filesystem scan output

```text
2026-07-01T15:51:45Z	INFO	[vuln] Vulnerability scanning is enabled
2026-07-01T15:51:45Z	INFO	Number of language-specific files	num=1
2026-07-01T15:51:45Z	INFO	[gomod] Detecting vulnerabilities...
```

## 5. Top of configuration scan output

```text
2026-07-01T15:53:30Z	INFO	[misconfig] Misconfiguration scanning is enabled
2026-07-01T15:53:30Z	INFO	No downloadable policies were loaded as --skip-check-update is enabled
2026-07-01T15:53:33Z	INFO	Detected config files	num=1
```

Trivy `0.59.1` used its pinned embedded configuration checks with
`--skip-check-update`. This avoided downloading an incompatible newer
policy bundle while retaining the checks shipped with the pinned scanner.

## 6. First 30 lines of the final CycloneDX SBOM

```json
{
    "$schema": "http://cyclonedx.org/schema/bom-1.6.schema.json",
    "bomFormat": "CycloneDX",
    "specVersion": "1.6",
    "serialNumber": "urn:uuid:b9f901b6-f51e-4c6c-a2e9-42728d892424",
    "version": 1,
    "metadata": {
        "timestamp": "2026-07-01T15:52:03+00:00",
        "tools": {
            "components": [
                {
                    "type": "application",
                    "group": "aquasecurity",
                    "name": "trivy",
                    "version": "0.59.1"
                }
            ]
        },
        "component": {
            "bom-ref": "64929bbd-9cda-4e78-a49e-d1ff7e94ad44",
            "type": "container",
            "name": "quicknotes:lab6",
            "properties": [
                {
                    "name": "aquasecurity:trivy:DiffID",
                    "value": "sha256:03496a3bc73427beb81a3d6bdb0fe7b93f9d4b2bef226f2bc0264de2800aa4e7"
                },
                {
                    "name": "aquasecurity:trivy:DiffID",
                    "value": "sha256:114dde0fefebbca13165d0da9c500a66190e497a82a53dcaabc3172d630be1e9"
```

## 7. Complete HIGH/CRITICAL triage

Each row below represents one scanner occurrence. The ten unique CVEs
appeared in both Go binaries, producing twenty total findings.

| CVE | Severity | Target | Package | Installed | Fixed version | Disposition | Reason |
| --- | --- | --- | --- | --- | --- | --- | --- |
| CVE-2026-25679 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.8, 1.26.1 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-25679 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.8, 1.26.1 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-27145 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.11, 1.26.4 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-27145 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.11, 1.26.4 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-32280 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.9, 1.26.2 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-32280 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.9, 1.26.2 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-32281 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.9, 1.26.2 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-32281 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.9, 1.26.2 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-32283 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.9, 1.26.2 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-32283 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.9, 1.26.2 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-33811 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-33811 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-33814 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-33814 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-39820 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-39820 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-39836 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-39836 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-42499 | HIGH | healthcheck | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |
| CVE-2026-42499 | HIGH | quicknotes | stdlib | v1.24.13 | 1.25.10, 1.26.3 | FIX | Upgraded the Docker builder from Go 1.24.13 to Go 1.26.4 in app/Dockerfile. The remediated image scan contains zero HIGH/CRITICAL findings. See the final PR and implementation commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c). |

The remediation changed the builder image in
[`app/Dockerfile`](../app/Dockerfile) from Go `1.24.13` to
Go `1.26.4`. The after-scan result is:

```text
2026-07-01T15:51:41Z	INFO	[vuln] Vulnerability scanning is enabled
2026-07-01T15:51:43Z	INFO	Detected OS	family="debian" version="12.14"
2026-07-01T15:51:43Z	INFO	[debian] Detecting vulnerabilities...	os_version="12" pkg_num=4
2026-07-01T15:51:43Z	INFO	Number of language-specific files	num=2
2026-07-01T15:51:43Z	INFO	[gobinary] Detecting vulnerabilities...

quicknotes:lab6 (debian 12.14)
==============================
Total: 0 (HIGH: 0, CRITICAL: 0)
```

## 8. Design questions a–d

### a) CVE severity is one input, not the answer. What else matters?

Severity measures the potential impact of a vulnerability in a general
environment; it does not establish the risk to this deployment. Triage
also needs:

- **Reachability:** whether QuickNotes actually calls the vulnerable code.
- **Exposure:** whether the affected path is reachable from an untrusted
  network or only from a restricted internal context.
- **Exploit maturity:** whether reliable public exploits exist and whether
  exploitation requires authentication or special preconditions.
- **Asset and data sensitivity:** the privileges of the process and the
  confidentiality or integrity of the data it can access.
- **Deployment controls:** non-root execution, a read-only filesystem,
  dropped Linux capabilities, network controls, and other compensating
  controls.
- **Fix availability and operational cost:** whether a patched version
  exists and whether upgrading introduces compatibility or availability
  risk.

A lower-severity reachable issue on an internet-facing service may deserve
a faster fix than a higher-severity issue in unreachable code.

### b) Why is a minimal or distroless base such a strong security control?

A minimal image removes packages before they can become vulnerabilities.
It normally contains no shell, package manager, compiler, debugging tools,
or unrelated utilities. This produces fewer CVEs, fewer executable
post-exploitation tools, fewer configuration paths, and a smaller image
whose contents are easier to audit. It also reduces patching workload
because there are fewer components to maintain.

Minimal images are not sufficient by themselves: the application binary
and language runtime can still contain vulnerabilities, as the original
Go `1.24.13` binaries demonstrated.

### c) When is `.trivyignore` appropriate, and when is it security theater?

An ignore entry is appropriate only after the finding has been reviewed
and shown to be non-applicable, temporarily accepted, or blocked on an
upstream fix. It should be narrowly scoped to the exact finding, include
an owner and rationale, and have an expiry or review date.

It becomes security theater when it is used to make the report green
without investigation, suppresses broad classes of findings, has no
expiry, remains after the affected component changes, or hides a finding
merely because remediation is inconvenient. No `.trivyignore` was needed
for this lab because all image findings had an available upgrade.

### d) What future problem does an SBOM solve?

The SBOM is a versioned inventory of the components in a specific release
artifact. When a new vulnerability such as Log4Shell is disclosed, the
inventory can be queried immediately to determine which images contain
the affected component and version. Without it, teams must search source
repositories, build systems, and running hosts manually.

It also supports targeted recalls, patch prioritization, incident response,
licence review, and evidence that the rebuilt artifact no longer contains
the affected version.

---

# Task 2 — OWASP ZAP baseline and code remediation

## 1. ZAP artifacts

### Before remediation

- [HTML report](../artifacts/lab9/zap/before/baseline-before.html)
- [JSON report](../artifacts/lab9/zap/before/baseline-before.json)
- [Console log](../artifacts/lab9/zap/before/baseline-before.log)
- [Markdown report](../artifacts/lab9/zap/before/baseline-before.md)
- [Raw `/health` headers](../artifacts/lab9/zap/before/health-headers-before.txt)
- [Raw `/` headers](../artifacts/lab9/zap/before/root-headers-before.txt)

### After remediation

- [HTML report](../artifacts/lab9/zap/after/baseline-after.html)
- [JSON report](../artifacts/lab9/zap/after/baseline-after.json)
- [Console log](../artifacts/lab9/zap/after/baseline-after.log)
- [Markdown report](../artifacts/lab9/zap/after/baseline-after.md)
- [Raw `/health` headers](../artifacts/lab9/zap/after/health-headers-after.txt)
- [Raw `/` headers](../artifacts/lab9/zap/after/root-headers-after.txt)
- [Machine-generated comparison](../artifacts/lab9/zap/zap-comparison.txt)

## 2. Complete ZAP triage

The table includes every alert instance from both saved JSON reports.

| State | ID / alert reference | Name | Risk | Method | URL | Parameter | Disposition | Reason |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Before fix | 10116 / 10116 | ZAP is Out of Date | Low (High) | GET | http://localhost:8080 | none | ACCEPT | This is a scanner self-version warning, not a QuickNotes application vulnerability. ZAP was intentionally pinned to 2.16.1 for reproducibility as required by the lab. Re-evaluate the scanner version before production use. |
| Before fix | 10049 / 10049-3 | Storable and Cacheable Content | Informational (Medium) | GET | http://localhost:8080 | none | FIX | Responses could be stored by intermediaries or browsers. Added router-wrapping middleware that sets `Cache-Control: no-store` on every route, including 404 responses. Alert reference 10049-3 disappeared after rebuild. |
| Before fix | 10049 / 10049-3 | Storable and Cacheable Content | Informational (Medium) | GET | http://localhost:8080/robots.txt | none | FIX | Responses could be stored by intermediaries or browsers. Added router-wrapping middleware that sets `Cache-Control: no-store` on every route, including 404 responses. Alert reference 10049-3 disappeared after rebuild. |
| After fix | 10116 / 10116 | ZAP is Out of Date | Low (High) | GET | http://localhost:8080 | none | ACCEPT | This is a scanner self-version warning, not a QuickNotes application vulnerability. ZAP was intentionally pinned to 2.16.1 for reproducibility as required by the lab. Re-evaluate the scanner version before production use. |
| After fix | 10049 / 10049-1 | Non-Storable Content | Informational (Medium) | GET | http://localhost:8080 | none | ACCEPT | This is the intended result of the remediation. ZAP reference 10049-1 confirms that the response is non-storable because `Cache-Control: no-store` is present. |

## 3. Code fix

The remediation is implemented in
[`app/middleware.go`](../app/middleware.go), wraps the router in
[`app/handlers.go`](../app/handlers.go), and is guarded by
[`app/middleware_test.go`](../app/middleware_test.go).

The implementation is visible in both the
[course PR](https://github.com/inno-devops-labs/DevOps-Intro/pull/1298) and [fork PR](https://github.com/tivdzualubem/DevOps-Intro/pull/9).

```go
package main

import "net/http"

// securityHeaders applies response-security policy consistently to every route,
// including responses generated by the router for unmatched paths.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Cache-Control", "no-store")
		headers.Set("Content-Security-Policy", "default-src 'none'")
		headers.Set(
			"Permissions-Policy",
			"camera=(), geolocation=(), microphone=()",
		)
		headers.Set("Referrer-Policy", "no-referrer")
		headers.Set("X-Content-Type-Options", "nosniff")
		headers.Set("X-Frame-Options", "DENY")

		next.ServeHTTP(w, r)
	})
}
```

The middleware applies the following headers to matched and unmatched
routes:

- `Cache-Control: no-store`
- `Content-Security-Policy: default-src 'none'`
- `Permissions-Policy: camera=(), geolocation=(), microphone=()`
- `Referrer-Policy: no-referrer`
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`

`TestSecurityHeadersApplyToAllRoutes` checks `/health`, `/notes`, and an
unmatched route. Removing `securityHeaders(mux)` causes the assertions to
fail.

## 4. Before/after proof

### Before

```text
WARN-NEW: Storable and Cacheable Content [10049] x 2
	http://localhost:8080 (404 Not Found)
	http://localhost:8080/robots.txt (404 Not Found)
WARN-NEW: ZAP is Out of Date [10116] x 1
	http://localhost:8080 (404 Not Found)
FAIL-NEW: 0	FAIL-INPROG: 0	WARN-NEW: 2	WARN-INPROG: 0	INFO: 0	IGNORE: 0	PASS: 65
```

The before headers contained no `Cache-Control` field.

### After

```text
WARN-NEW: Non-Storable Content [10049] x 1
	http://localhost:8080 (404 Not Found)
WARN-NEW: ZAP is Out of Date [10116] x 1
	http://localhost:8080 (404 Not Found)
FAIL-NEW: 0	FAIL-INPROG: 0	WARN-NEW: 2	WARN-INPROG: 0	INFO: 0	IGNORE: 0	PASS: 65
```

The alert changed from:

- `10049-3 — Storable and Cacheable Content`

to:

- `10049-1 — Non-Storable Content`

The plugin ID is shared by multiple cacheability sub-rules, so comparing
only `10049` would be incorrect. The alert name and alert reference prove
that the vulnerable cacheable condition disappeared and the intended
`no-store` condition replaced it.

## 5. Design questions e–g

### e) Why middleware instead of setting headers in every handler?

Middleware creates one policy enforcement point around the entire router.
It prevents handlers from being forgotten, covers framework-generated
responses such as 404 and method errors, keeps business logic separate
from transport security policy, and makes the behavior easy to review and
test. Per-handler calls would duplicate code and would eventually produce
inconsistent protection.

### f) What does `Content-Security-Policy: default-src 'none'` break?

It blocks all browser-loaded resources by default, including scripts,
stylesheets, images, fonts, media, frames, and network connections unless
another CSP directive explicitly allows them. A normal website would lose
its styling, client-side JavaScript, images, and API calls.

QuickNotes is a JSON API and does not serve a browser user interface or
legitimate page resources, so a deny-by-default policy is appropriate.
A website would need narrowly scoped directives such as `script-src`,
`style-src`, `img-src`, and `connect-src`.

### g) What is the cost of accepting every informational alert without reading it?

Blanket acceptance destroys the distinction between harmless noise and an
early indicator of a real weakness. It creates an unreliable audit trail,
normalizes risk, hides regressions, and trains reviewers to ignore the
scanner. Each alert still needs context: the ZAP version warning is a
tooling concern, while the cacheability alert represented application
behavior that could be corrected.

---

# Bonus — `govulncheck` CI PR gate

## 1. CI job

The separate status check is defined in
[`.github/workflows/ci.yml`](../.github/workflows/ci.yml).

```yaml
  govulncheck:
    name: govulncheck-go-1.24
    runs-on: ubuntu-24.04
    defaults:
      run:
        working-directory: app
    steps:
      - name: Check out repository
        uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

      - name: Set up required Go 1.24 environment
        uses: actions/setup-go@d35c59abb061a4a6fb18e82ac0862c26744d6ab5 # v5.5.0
        with:
          go-version: '1.24.13'
          cache: false

      - name: Verify required Go 1.24 setup
        env:
          GOTOOLCHAIN: local
        run: go version

      - name: Install pinned govulncheck
        env:
          GOTOOLCHAIN: go1.26.4
        run: go install golang.org/x/vuln/cmd/govulncheck@v1.5.0

      - name: Run reachable-vulnerability scan
        env:
          GOTOOLCHAIN: go1.26.4
        run: '"$(go env GOPATH)/bin/govulncheck" ./...'
```

The workflow provisions the required Go `1.24.13` CI environment. The
scanner installation and source analysis use `GOTOOLCHAIN=go1.26.4`,
matching the patched production builder. This is necessary because the
current vulnerability database reports eight reachable standard-library
advisories in Go `1.24.13`; using that vulnerable toolchain would keep the
clean branch permanently red independently of the deliberately introduced
dependency.

## 2. Red demonstration

The red commit deliberately added `golang.org/x/text v0.3.5` and called
`language.Parse`, making `GO-2021-0113` reachable.

Local scanner evidence:

```text
=== Symbol Results ===

Vulnerability #1: GO-2021-0113
    Out-of-bounds read in golang.org/x/text/language
  More info: https://pkg.go.dev/vuln/GO-2021-0113
  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.3.5
    Fixed in: golang.org/x/text@v0.3.7
    Example traces found:
      #1: vuln_demo.go:8:23: quicknotes.runVulnerableDependencyDemo calls language.Parse

Your code is affected by 1 vulnerability from 1 module.
This scan also found 1 vulnerability in packages you import and 0
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
Use '-show verbose' for more details.
```

GitHub CI evidence:

```text
RED CI CHECK EVIDENCE
========================================
Commit: 332bec4b61fa457fe9be771378b2bc7238c8aa3c
Fork PR: https://github.com/tivdzualubem/DevOps-Intro/pull/9

govulncheck-go-1.24 | status=completed | conclusion=failure | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529310985/job/84574247121
lint | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529310985/job/84574247235
test-go-1.23 | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529310985/job/84574247138
test-go-1.24 | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529310985/job/84574247167
vet-go-1.23 | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529310985/job/84574247108
vet-go-1.24 | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529310985/job/84574247258
```

At commit [332bec4](https://github.com/tivdzualubem/DevOps-Intro/commit/332bec4b61fa457fe9be771378b2bc7238c8aa3c):

- `govulncheck-go-1.24` failed.
- The other five CI checks passed.

## 3. Green recovery

The vulnerable source file, call, dependency, and `go.sum` entries were
removed in commit [053da8c](https://github.com/tivdzualubem/DevOps-Intro/commit/053da8ce5c4c6883d210b1f9332912a010446669).

Local scanner evidence:

```text
No vulnerabilities found.
```

GitHub CI evidence:

```text
GREEN CI CHECK EVIDENCE
========================================
Commit: 053da8ce5c4c6883d210b1f9332912a010446669
Fork PR: https://github.com/tivdzualubem/DevOps-Intro/pull/9

govulncheck-go-1.24 | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529636247/job/84575398381
lint | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529636247/job/84575398472
test-go-1.23 | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529636247/job/84575398528
test-go-1.24 | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529636247/job/84575398482
vet-go-1.23 | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529636247/job/84575398432
vet-go-1.24 | status=completed | conclusion=success | url=https://github.com/tivdzualubem/DevOps-Intro/actions/runs/28529636247/job/84575398442
```

All six status checks passed on the green commit. The final branch does
not contain the deliberately vulnerable dependency.

## 4. Design questions h–j

### h) How does reachability change vulnerability triage?

A module-level match proves that a vulnerable version is present, but not
that the application invokes the affected symbol. Reachability analysis
uses the call graph to distinguish an imported but unused vulnerable path
from a path the application can execute.

This reduces triage workload and allows reachable issues to receive the
highest priority. It does not make unreachable findings irrelevant:
build tags, reflection, plugins, future code changes, and alternate
execution paths can change reachability, so the dependency should still
be monitored or upgraded when practical.

### i) Why pin the scanner version?

Pinning makes CI reproducible. Scanner releases can change command-line
behavior, supported Go versions, output formats, rules, defaults, and exit
codes. Installing `@latest` allows an unrelated upstream release to break
or silently change a previously reviewed pipeline.

A pinned scanner version also provides a clear audit record and allows
upgrades to be reviewed intentionally in their own change. The
vulnerability database can continue to update independently.

### j) What will govulncheck miss that Trivy image scanning can find?

`govulncheck` focuses on Go modules, the Go standard library, and reachable
Go call paths. It does not provide full coverage for:

- Debian, Alpine, or other operating-system packages.
- The container base image and image layers.
- Non-Go binaries and libraries.
- Dockerfile, Compose, Kubernetes, or Terraform misconfigurations.
- Embedded JavaScript, Python, Java, or native dependencies.
- Secrets, licences, excessive permissions, or unsafe container runtime
  configuration.

Trivy complements it by scanning the assembled artifact and infrastructure
configuration rather than only the Go call graph.

---

# Final verification summary

- All four required Trivy operations were captured.
- All twenty HIGH image occurrences were individually triaged and fixed.
- The final image contains zero HIGH/CRITICAL findings.
- The repository filesystem and configuration scans contain zero
  HIGH/CRITICAL findings.
- A CycloneDX `1.6` SBOM is committed.
- Passive ZAP HTML and JSON reports are committed for before and after.
- Every ZAP alert instance is triaged.
- The cacheability issue was fixed through router-wrapping middleware.
- The middleware is applied to all routes and protected by a unit test.
- The red and green `govulncheck` CI states are recorded.
- The final branch has a green security gate and no deliberate vulnerable
  dependency.
