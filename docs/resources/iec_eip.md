---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_eip

Manages a eip resource within HuaweiCloud IEC.

## Example Usage

```hcl
data "huaweicloud_iec_sites" "iec_sites" {}

resource "huaweicloud_iec_eip" "eip_test" {
  site_id    = data.huaweicloud_iec_sites.iec_sites.sites[0].id
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required, String, ForceNew) Specifies the id of iec sevice site. 
    Changing this parameter creates a new iec eip resource.

* `port_id` - (Optional, String) Specifies the port id which this eip will associate with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `ip_version` - The version of elastic IP address. 
    IEC services only support IPv4(4) now.
    
* `status` - The status of iec eip.

* `public_ip` - The address of elastic IP.

* `private_ip` - The address of private IP.

* `bandwitch_name` - The name of bandwidth.

* `bandwitch_size` - The size of bandwidth.

* `bandwitch_share_type` - Whether the bandwidth is shared or exclusive. 

* `site_info` - The located information of the iec site. It contains 
    area, province and city.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 3 minute.

## Import

IEC EIPs can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iec_eip.eip_test b5ad19d1-57d1-48fd-aab7-1378f9bee169
```
