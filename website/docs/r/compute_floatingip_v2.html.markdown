---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_floatingip_v2"
sidebar_current: "docs-huaweicloud-resource-compute-floatingip-v2"
description: |-
  Manages a V2 floating IP resource within HuaweiCloud Nova (compute).
---

# huaweicloud\_compute\_floatingip_v2

Manages a V2 floating IP resource within HuaweiCloud Nova (compute)
that can be used for compute instances.

Please note that managing floating IPs through the HuaweiCloud Compute API has
been deprecated. Unless you are using an older HuaweiCloud environment, it is
recommended to use the [`huaweicloud_networking_floatingip_v2`](networking_floatingip_v2.html)
resource instead, which uses the HuaweiCloud Networking API.

## Example Usage

```hcl
resource "huaweicloud_compute_floatingip_v2" "floatip_1" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    A Compute client is needed to create a floating IP that can be used with
    a compute instance. If omitted, the `region` argument of the provider
    is used. Changing this creates a new floating IP (which may or may not
    have a different address).

* `pool` - (Optional) The name of the pool from which to obtain the floating
    IP. Only admin_external_net is valid. Changing this creates a new floating IP.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `pool` - See Argument Reference above.
* `address` - The actual floating IP address itself.
* `fixed_ip` - The fixed IP address corresponding to the floating IP.
* `instance_id` - UUID of the compute instance associated with the floating IP.

## Import

Floating IPs can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_compute_floatingip_v2.floatip_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
