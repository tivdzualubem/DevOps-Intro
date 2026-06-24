# QuickNotes High Error Rate Runbook

## What this alert means

More than 5% of QuickNotes requests have returned HTTP 4xx or 5xx responses continuously for at least five minutes, indicating sustained user-visible failures.

## Triage steps

1. Confirm that the alert is still firing and inspect the current error ratio in Prometheus or the Grafana Golden Signals dashboard.

2. Check whether all containers are running:

        docker compose ps

3. Inspect recent QuickNotes logs:

        docker compose logs --tail=200 quicknotes

4. Identify which response codes dominate using PromQL:

        sum by (code) (
          rate(quicknotes_http_responses_by_code_total[5m])
        )

5. Test the service directly:

        curl -i http://localhost:8080/health
        curl -i http://localhost:8080/notes

## Mitigations

1. Stop or rate-limit any client generating malformed or excessive requests.

2. Restart QuickNotes if it is unhealthy or stuck:

        docker compose restart quicknotes

3. Roll back the latest application or configuration change if errors began after deployment.

4. Temporarily reduce non-essential traffic while preserving health checks and normal reads.

## Post-incident

After recovery, preserve logs and monitoring evidence, identify the root cause, and complete a blameless postmortem using the Lecture 1 postmortem structure. Assign an owner and deadline to every preventive action.
