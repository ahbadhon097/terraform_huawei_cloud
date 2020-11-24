---
subcategory: "Deprecated"
---

# huaweicloud\_blockstorage\_volume_v2

Manages a V2 volume resource within HuaweiCloud.

!> **Warning:** It has been deprecated, use `huaweicloud_evs_volume` instead.

## Example Usage

```hcl
resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  region      = "RegionOne"
  name        = "volume_1"
  description = "first test volume"
  size        = 3
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the volume. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new volume.

* `size` - (Required) The size of the volume to create (in gigabytes).

* `availability_zone` - (Optional) The availability zone for the volume.
    Changing this creates a new volume.

* `consistency_group_id` - (Optional) The consistency group to place the volume
    in.

* `description` - (Optional) A description of the volume. Changing this updates
    the volume's description.

* `image_id` - (Optional) The image ID from which to create the volume.
    Changing this creates a new volume.

* `metadata` - (Optional) Metadata key/value pairs to associate with the volume.
    Changing this updates the existing volume metadata.

* `name` - (Optional) A unique name for the volume. Changing this updates the
    volume's name.

* `snapshot_id` - (Optional) The snapshot ID from which to create the volume.
    Changing this creates a new volume.

* `source_replica` - (Optional) The volume ID to replicate with.

* `source_vol_id` - (Optional) The volume ID from which to create the volume.
    Changing this creates a new volume.

* `volume_type` - (Optional) The type of volume to create. Available types are
    `SSD`, `SAS` and `SATA`. Changing this creates a new volume.

* `cascade` - (Optional, Default:false) Specifies to delete all snapshots associated with the EVS disk.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `size` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `image_id` - See Argument Reference above.
* `source_vol_id` - See Argument Reference above.
* `snapshot_id` - See Argument Reference above.
* `metadata` - See Argument Reference above.
* `volume_type` - See Argument Reference above.
* `attachment` - If a volume is attached to an instance, this attribute will
    display the Attachment ID, Instance ID, and the Device as the Instance
    sees it.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

Volumes can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_blockstorage_volume_v2.volume_1 ea257959-eeb1-4c10-8d33-26f0409a755d
```
