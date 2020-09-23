---
subcategory: "Anti-DDoS"
---

# huaweicloud\_antiddos

Query the Anti-DDos resource.
This is an alternative to `huaweicloud_antiddos_v1`

## Example Usage

```hcl
variable "eip_id" {}

data "huaweicloud_antiddos" "antiddos" {
  floating_ip_id = var.eip_id
}

```

## Argument Reference
The following arguments are supported:

* `floating_ip_id` - (Optional) The Elastic IP ID.

* `floating_ip_address` - (Optional) The Elastic IP address.

* `status` - (Optional) The defense status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_type` - The EIP type.

* `period_start` - The Start time.

* `bps_attack` - The Attack traffic in (bit/s).

* `bps_in` - The inbound traffic in (bit/s).

* `total_bps` - The total traffic.

* `pps_in` - The inbound packet rate (number of packets per second).

* `pps_attack` - The attack packet rate (number of packets per second).

* `total_pps` - The total packet rate.

* `start_time` - The start time of cleaning and blackhole event.

* `end_time` - The end time of cleaning and blackhole event.

* `traffic_cleaning_status` - The traffic cleaning status.

* `trigger_bps` - The traffic at the triggering point.

* `trigger_pps` - The packet rate at the triggering point.

* `trigger_http_pps` - The HTTP request rate at the triggering point.

