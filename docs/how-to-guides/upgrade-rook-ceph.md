# Upgrading Rook Ceph

## Contents

- [Introduction](#introduction)
- [Steps](#steps)
  - [Step 1: Ensure `AUTOSCALE` is `on`](#step-1-ensure-autoscale-is-on)
  - [Step 2: Watch](#step-2-watch)
    - [Step 2.1: Ceph status](#step-21-ceph-status)
    - [Step 2.2: Pods in rook namespace](#step-22-pods-in-rook-namespace)
    - [Step 2.3: Rook version update](#step-23-rook-version-update)
    - [Step 2.4: Ceph version update](#step-24-ceph-version-update)
    - [Step 2.5: Events in rook namespace](#step-25-events-in-rook-namespace)
  - [Step 3: Dashboards](#step-3-dashboards)
    - [Step 3.1: Ceph](#step-31-ceph)
    - [Step 3.2: Grafana](#step-32-grafana)
  - [Step 4: Make a note of existing image versions](#step-4-make-a-note-of-existing-image-versions)
  - [Step 5: Perform updates](#step-5-perform-updates)
  - [Step 6: Verify that the CSI images are updated](#step-6-verify-that-the-csi-images-are-updated)
  - [Step 7: Final checks](#step-7-final-checks)
- [Additional resources](#additional-resources)

## Introduction

Rook Ceph is one of the storage providers of Lokomotive. With a distributed system as complex as
Ceph, the upgrade process is not trivial. This document enlists steps on how to perform the upgrade
and how to monitor this process.

## Steps

Following steps are inspired by [`rook`](https://rook.io/docs/rook/master/ceph-upgrade.html) docs.

### Step 1: Ensure `AUTOSCALE` is `on`

Exec into the toolbox pod as specified in [this
doc](rook-ceph-storage.md#enable-and-access-toolbox). Once you get the shell access, run the
following command:

```console
# ceph osd pool autoscale-status | grep replicapool
POOL           SIZE  TARGET SIZE  RATE  RAW CAPACITY   RATIO  TARGET RATIO  EFFECTIVE RATIO  BIAS  PG_NUM  NEW PG_NUM  AUTOSCALE
replicapool      0                 3.0         3241G  0.0000                                  1.0      32              on
```

Ensure that the `AUTOSCALE` column outputs `on` and not `warn`. If the output of the `AUTOSCALE`
column says `warn`, then run the following command, to make sure that pool autoscaling is enabled.
This is required to ensure that placement groups scale up as the data in the cluster increases.

```bash
ceph osd pool set replicapool pg_autoscale_mode on
```

### Step 2: Watch

Watch events, updates and pods.

#### Step 2.1: Ceph status

Leave the following running in the toolbox pod:

```bash
watch ceph status
```

Ensure that the output says that `health:` is `HEALTH_OK`. Match the output such that everything
looks fine as explained in the [rook upgrade
docs](https://rook.io/docs/rook/master/ceph-upgrade.html#status-output).

> **IMPORTANT**: Don't proceed further if the output is anything other than `HEALTH_OK`.

#### Step 2.2: Pods in rook namespace

Keep an eye on the pods status in another terminal window from the `rook` namespace. Leave the
following command running:

```bash
watch kubectl -n rook get pods -o wide
```

#### Step 2.3: Rook version update

Run the following command in a new terminal window to keep an eye on the rook version update as it
is upgrades for all the sub-components:

```bash
watch --exec kubectl -n rook get deployments -l rook_cluster=rook -o \
  jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \trook-version="}{.metadata.labels.rook-version}{"\n"}{end}'
```

```bash
watch --exec kubectl -n rook get jobs -o \
  jsonpath='{range .items[*]}{.metadata.name}{"  \tsucceeded: "}{.status.succeeded}{"      \trook-version="}{.metadata.labels.rook-version}{"\n"}{end}'
```

You should see that `rook-version` slowly changes to `v1.4.6`.

#### Step 2.4: Ceph version update

Run the following command to keep an eye on the Ceph version update as the new pods come up in a new
terminal window:

```bash
watch --exec kubectl -n rook get deployments -l rook_cluster=rook -o \
  jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \tceph-version="}{.metadata.labels.ceph-version}{"\n"}{end}'
```

You should see that `ceph-version` slowly changes to `15.2.5`.

#### Step 2.5: Events in rook namespace

In a new terminal leave the following command running, to keep track of the events happening in the
rook namespace:

```bash
kubectl -n rook get events -w
```

### Step 3: Dashboards

Monitor various dashboards.

#### Step 3.1: Ceph

Open the Ceph dashboard in a browser window. Instructions to access the dashboard can be found
[here](rook-ceph-storage.md#access-the-ceph-dashboard).

> **NOTE**: Accessing dashboard can be a hassle because while the components are upgrading you may
> lose access to it multiple times.

#### Step 3.2: Grafana

Gain access to the Grafana dashboard as instructed
[here](monitoring-with-prometheus-operator.md#access-grafana). And keep an eye on the dashboard
named `Ceph - Cluster`.

> **NOTE**: The data in the Grafana dashboard will always be outdated compared to the `watch ceph
> status` running inside the toolbox pod.

### Step 4: Make a note of existing image versions

Make a note of the images of the pods in the rook namespace:

```bash
kubectl -n rook get pod -o \
  jsonpath='{range .items[*]}{.metadata.name}{"\n\t"}{.status.phase}{"\t\t"}{.spec.containers[0].image}{"\t"}{.spec.initContainers[0].image}{"\n\n"}{end}'
```

### Step 5: Perform updates

With everything monitored, you can start the update process now by executing the following commands:

```bash
kubectl apply -f https://raw.githubusercontent.com/kinvolk/lokomotive/master/assets/charts/components/rook/templates/resources.yaml
lokoctl component apply rook rook-ceph
```

### Step 6: Verify that the CSI images are updated

Verify if the images were updated, comparing it with output of the [Step
4](#step-4-make-a-note-of-existing-image-versions).

```bash
kubectl -n rook get pod -o \
  jsonpath='{range .items[*]}{.metadata.name}{"\n\t"}{.status.phase}{"\t\t"}{.spec.containers[0].image}{"\t"}{.spec.initContainers[0].image}{"\n\n"}{end}'
```

### Step 7: Final checks

Once everything is up to date then run the following commands in the toolbox pod, to verify if all
the OSDs are in `up` state:

```bash
ceph osd status
```

## Additional resources

- [Rook Upgrade docs](https://rook.io/docs/rook/v1.4/ceph-upgrade.html).
- [General Troubleshooting](https://rook.io/docs/rook/v1.5/common-issues.html).
- [Ceph Troubleshooting](https://rook.io/docs/rook/v1.4/ceph-common-issues.html).
