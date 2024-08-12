---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_baseline_check_results"
description: |-
  Use this data source to get the list of SecMaster baseline check results.
---

# huaweicloud_secmaster_baseline_check_results

Use this data source to get the list of SecMaster baseline check results.

## Example Usage

```hcl
variable "from_date" {}
variable "to_date" {}

data "huaweicloud_secmaster_baseline_check_results" "test" {
  from_date = var.from_date
  to_date   = var.to_date
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `from_date` - (Optional, String) Specifies the start time.
  The format is ISO 8601: YYYY-MM-DDTHH:mm:ss.ms+Time zone. Time zone refers to where the incident occurred.
  If this parameter cannot be parsed, the default time zone GMT+8 is used.

* `to_date` - (Optional, String) Specifies the end time.
  The format is ISO 8601: YYYY-MM-DDTHH:mm:ss.ms+Time zone. Time zone refers to where the incident occurred.
  If this parameter cannot be parsed, the default time zone GMT+8 is used.

* `condition` - (Optional, Map) Specifies the condition expression.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `baseline_check_results` - The list of baseline check result.
