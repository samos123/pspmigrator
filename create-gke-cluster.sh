#!/usr/bin/env bash

set -xe

gcloud beta container clusters create test-psp \
    --enable-pod-security-policy \
    --enable-ip-alias --region us-central1 --node-locations=us-central1-a \
    --machine-type=e2-medium --enable-private-nodes --master-ipv4-cidr "172.16.0.0/28" \
    --workload-pool "sam-argolis.svc.id.goog" --enable-shielded-nodes --shielded-secure-boot \
    --shielded-integrity-monitoring
