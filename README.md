[![CircleCI](https://circleci.com/gh/giantswarm/capi-garbage-collector.svg?style=shield)](https://circleci.com/gh/giantswarm/capi-garbage-collector)

# capi-garbage-collector

Clean up leftover resources caused by bugs or unexpected situations

| Resource Name      | Reason |
| :---        |    :----:   |
| MachinePool      | !machinePool.DeletionTimestamp.IsZero()       |
