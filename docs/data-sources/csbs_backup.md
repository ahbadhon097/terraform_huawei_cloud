---
subcategory: "Deprecated"
---

# huaweicloud\_csbs\_backup

!> **Warning:** It has been deprecated.

The HuaweiCloud CSBS Backup data source allows access of backup resources.
This is an alternative to `huaweicloud_csbs_backup_v1`

## Example Usage


```hcl
variable "backup_name" {}

data "huaweicloud_csbs_backup" "csbs" {
  backup_name = var.backup_name
}
```

## Argument Reference
The following arguments are supported:

* `region` - (Optional) The region in which to obtain the CSBS backups. If omitted, the provider-level region will work as default.

* `id` - (Optional) Specifies the ID of backup.

* `backup_name` - (Optional) Specifies the backup name.

* `status` - (Optional) Specifies the backup status.

* `resource_name` - (Optional) Specifies the backup object name.

* `backup_record_id` - (Optional) Specifies the backup record ID.

* `resource_type` - (Optional) Specifies the type of backup objects.

* `resource_id` - (Optional) Specifies the backup object ID.

* `policy_id` - (Optional) Specifies the Policy Id.

* `vm_ip` - (Optional) Specifies the ip of VM.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `description` - Provides the backup description.

* `auto_trigger` - Specifies whether automatic trigger is enabled.

* `average_speed` - Specifies average speed.

* `size` - Specifies the backup capacity.

* volume_backups 

  * `space_saving_ratio` - Specifies the space saving rate.

  * `volume_backups` block supports the following arguments:

  * `status` -  Status of backup Volume.
    
  * `space_saving_ratio` -  Specifies space saving rate.

  * `name` -  It gives EVS disk backup name.

  * `bootable` -  Specifies whether the disk is bootable.

  * `average_speed` -  Specifies the average speed.

  * `source_volume_size` -  Shows source volume size in GB.

  * `source_volume_id` -  It specifies source volume ID.

  * `incremental` -  Shows whether incremental backup is used.

  * `snapshot_id` -  ID of snapshot.

  * `source_volume_name` -  Specifies source volume name.

  * `image_type` -  It specifies backup. The default value is backup.

  * `id` -  Specifies Cinder backup ID.

  * `size` -  Specifies accumulated size (MB) of backups.

* `vm_metadata` block supports the following arguments:

  * `name` - Name of backup data.

  * `eip` - Specifies elastic IP address of the ECS.

  * `cloud_service_type` - Specifies ECS type.

  * `ram` - Specifies memory size of the ECS, in MB.

  * `vcpus` - Specifies CPU cores corresponding to the ECS.

  * `private_ip` - It specifies internal IP address of the ECS.

  * `disk` - Shows system disk size corresponding to the ECS specifications.

  * `image_type` - Specifies image type.
