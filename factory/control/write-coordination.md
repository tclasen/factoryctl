---
type: "Factory Control"
title: "Storage write coordination"
description: "Storage write coordination."
status: "draft"
factory:
  format: 1
  assertion:
    frame: "mixed"
    certainty: "mixed"
    basis: "mixed"
sources:
  -
    id: "operator"
    resource: "../sources/operator-brief.md"
    title: "Operator requirements and scope"
    author: "human:operator"
purpose: "native-advisory-lock-anchor"
---

# Storage write coordination

This stable, versioned knowledge node explains the proposed native writer-lock anchor. Its content is not a lock, lease, ownership registry or change log. Factoryctl processes coordinate through an operating-system advisory lock on this file while validating and atomically writing a record. They never update this file to record lock events.

The [storage contract](../specifications/knowledge-format.md) specifies inode/identity rechecks, crash release and limits with external editors/Git. Preserve this node during normal writes. Removing or replacing it externally causes active writers to fail safely when identity checks detect the change. The document CLI now uses this anchor for advisory writer coordination; [CLI evidence](../sources/cli-milestone-evidence.md) records verified behavior and remaining external-writer limits.
