Rig
===

_A rig is an arrangement of items used for fishing. It can be assembled of one or more lines, hooks, sinkers, bobbers, swivels, lures, beads, and other fishing tackle. A rig might be held by a rod, by hand, or attached to a boat or pier._

## Local Development Environment

### Docker Compose

The entire stack is managed by Docker Compose.

To spin it up run the following:

`docker-compose up`

### Registry

Rig builds and pushes images to a registry. To make this all possible we require a consistent DNS entry between the developers local and the cluster.

This is typically a non issue when using cloud registries eg. hub.docker.com.

However, on local this is an issue as the registry is hosted on our test cluster.

To smooth this over we portforward the registry (port 5000) and setup a consistent static DNS entry.

`127.0.0.1	registry.rig.svc.cluster.local:5000`
