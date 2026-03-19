"""LitmusChaos experiment lifecycle: pod selection, YAML rendering, start, and wait."""

import logging
import random
import tempfile
import time
from pathlib import Path

from kubernetes import client as k8s_client

from .helpers import get_service_nodeport, run, run_json

CHAOS_TEMPLATE = Path(__file__).resolve().parent.parent / "charts" / "chaos" / "client-experiment.yaml"

log = logging.getLogger(__name__)


def select_target_pods(
    v1: k8s_client.CoreV1Api,
    namespace: str,
    pods_affected_perc: int,
    seed: int,
) -> list[str]:
    """Deterministically select a subset of client pods based on percentage and seed.

    :param v1: Kubernetes CoreV1Api client
    :param namespace: Kubernetes namespace
    :param pods_affected_perc: Percentage of pods to target (1-100)
    :param seed: Random seed for reproducibility
    :return: List of selected pod names
    """
    pods = v1.list_namespaced_pod(
        namespace,
        label_selector="app.kubernetes.io/component=fedn-client",
    )
    pod_names = sorted(p.metadata.name for p in pods.items if p.status.phase == "Running")

    count = max(1, len(pod_names) * pods_affected_perc // 100)
    rng = random.Random(seed)
    selected = rng.sample(pod_names, min(count, len(pod_names)))
    log.info(
        "Selected %d/%d client pods (perc=%d%%, seed=%d): %s",
        len(selected), len(pod_names), pods_affected_perc, seed, selected,
    )
    return selected


def render_chaos_yaml(target_pods: list[str]) -> Path:
    """Read the chaos experiment template and fill in TARGET_PODS.

    Returns the path to a temporary YAML file with the value substituted.
    """
    content = CHAOS_TEMPLATE.read_text()
    content = content.replace("__TARGET_PODS__", ",".join(target_pods))
    tmp = tempfile.NamedTemporaryFile(
        mode="w", suffix=".yaml", prefix="chaos-experiment-", delete=False,
    )
    tmp.write(content)
    tmp.close()
    return Path(tmp.name)


def start_chaos_experiment(
    v1: k8s_client.CoreV1Api,
    node_ip: str,
    namespace: str,
    litmus_namespace: str,
    target_pods: list[str],
) -> tuple[str, str, Path] | None:
    """Start the LitmusChaos pod-delete experiment. Returns (project_id, exp_id) or None."""
    chaos_port = get_service_nodeport(
        v1,
        "chaos-litmus-frontend-service",
        litmus_namespace,
        "http",
    )

    run(
        [
            "litmusctl",
            "config",
            "set-account",
            "-n",
            "-e",
            f"http://{node_ip}:{chaos_port}",
            "-u",
            "admin",
            "-p",
            "litmus",
        ]
    )

    projects = run_json(["litmusctl", "get", "projects", "-o", "json"])
    project_id = projects["projects"][0]["projectID"]

    infras = run_json(
        [
            "litmusctl",
            "get",
            "chaos-infra",
            "--project-id",
            project_id,
            "-o",
            "json",
        ]
    )
    infra_id = infras["listInfras"]["infras"][0]["infraID"]

    log.info("Project: %s  Infra: %s", project_id, infra_id)

    # Delete existing experiment so it can be recreated with updated parameters
    experiments = run_json(
        [
            "litmusctl",
            "get",
            "chaos-experiments",
            "--project-id",
            project_id,
            "-o",
            "json",
        ]
    )
    existing = [e["experimentID"] for e in experiments["listExperiment"]["experiments"] if "client-experiment" in e.get("name", "")]

    for old_id in existing:
        log.info("Deleting old chaos experiment %s...", old_id)
        run(
            [
                "litmusctl",
                "delete",
                "chaos-experiment",
                old_id,
                "--project-id",
                project_id,
            ],
            check=False,
        )

    # Create experiment from template with TARGET_PODS filled in
    log.info("Creating chaos experiment (target_pods=%s)...", target_pods)
    chaos_yaml = render_chaos_yaml(target_pods)
    log.debug("Rendered chaos YAML: %s", chaos_yaml)
    run(
        [
            "litmusctl",
            "create",
            "chaos-experiment",
            "-f",
            str(chaos_yaml),
            "--project-id",
            project_id,
            "--chaos-infra-id",
            infra_id,
            "-d",
            "client-experiment",
        ]
    )

    # Fetch the new experiment ID
    experiments = run_json(
        [
            "litmusctl",
            "get",
            "chaos-experiments",
            "--project-id",
            project_id,
            "-o",
            "json",
        ]
    )
    existing = [e["experimentID"] for e in experiments["listExperiment"]["experiments"] if "client-experiment" in e.get("name", "")]
    exp_id = existing[0]

    log.info("Running experiment %s...", exp_id)
    run(
        [
            "litmusctl",
            "run",
            "chaos-experiment",
            "--project-id",
            project_id,
            "--experiment-id",
            exp_id,
        ]
    )

    return project_id, exp_id, chaos_yaml


def wait_for_chaos_experiment(
    project_id: str,
    exp_id: str,
    timeout: int,
    poll_interval: int,
):
    """Poll until the chaos experiment completes."""
    log.info("Waiting for chaos experiment to complete (timeout: %ds)...", timeout)
    elapsed = 0
    while elapsed < timeout:
        runs = run_json(
            [
                "litmusctl",
                "get",
                "chaos-experiment-runs",
                "--experiment-id",
                exp_id,
                "--project-id",
                project_id,
                "-o",
                "json",
            ]
        )
        phase = runs["listExperimentRun"]["experimentRuns"][0]["phase"]

        if phase == "Completed":
            log.info("Chaos experiment completed.")
            return

        log.info("  Experiment phase: %s (%ds elapsed)...", phase, elapsed)
        time.sleep(poll_interval)
        elapsed += poll_interval

    log.warning("Chaos experiment did not complete within %ds. Continuing anyway.", timeout)
