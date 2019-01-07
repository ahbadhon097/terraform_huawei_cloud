## 1.3.0 (Unreleased)

FEATURES:

* **New Data Source:** `huaweicloud_dms_az_v1` [GH-41]
* **New Data Source:** `huaweicloud_dms_product_v1` [GH-41]
* **New Data Source:** `huaweicloud_dms_maintainwindow_v1` [GH-41]
* **New Data Source:** `huaweicloud_vbs_backup_policy_v2` [GH-44]
* **New Data Source:** `huaweicloud_vbs_backup_v2` [GH-44]
* **New Data Source:** `huaweicloud_cce_cluster_v3` [GH-19]
* **New Data Source:** `huaweicloud_cce_node_v3` [GH-19]
* **New Data Source:** `huaweicloud_cts_tracker_v1` [GH-46]
* **New Data Source:** `huaweicloud_csbs_backup_v1` [GH-42]
* **New Data Source:** `huaweicloud_csbs_backup_policy_v1` [GH-42]
* **New Data Source:** `huaweicloud_antiddos_v1` [GH-47]
* **New Data Source:** `huaweicloud_dcs_az_v1` [GH-55]
* **New Data Source:** `huaweicloud_dcs_maintainwindow_v1` [GH-55]
* **New Data Source:** `huaweicloud_dcs_product_v1` [GH-55]
* **New Resource:** `huaweicloud_dms_queue_v1` [GH-41]
* **New Resource:** `huaweicloud_dms_group_v1` [GH-41]
* **New Resource:** `huaweicloud_dms_instance_v1` [GH-41]
* **New Resource:** `huaweicloud_vbs_backup_policy_v2` [GH-44]
* **New Resource:** `huaweicloud_vbs_backup_v2` [GH-44]
* **New Resource:** `huaweicloud_cce_cluster_v3` [GH-19]
* **New Resource:** `huaweicloud_cce_node_v3` [GH-19]
* **New Resource:** `huaweicloud_cts_tracker_v1` [GH-46]
* **New Resource:** `huaweicloud_csbs_backup_v1` [GH-42]
* **New Resource:** `huaweicloud_csbs_backup_policy_v1` [GH-42]
* **New Resource:** `huaweicloud_mrs_cluster_v1` [GH-56]
* **New Resource:** `huaweicloud_mrs_job_v1` [GH-56]
* **New Resource:** `huaweicloud_dcs_instance_v1` [GH-55]
* **New Resource:** `huaweicloud_maas_task_v1` [GH-65]
* **New Resource:** `huaweicloud_networking_floatingip_associate_v2` [GH-68]
* **New Resource:** `huaweicloud_dws_cluster` [GH-69]
* **New Resource:** `huaweicloud_mls_instance` [GH-69]

BUG FIXES:

* `resource/huaweicloud_elb_listener`: Fix certificate_id check [GH-45]
* `resource/huaweicloud_smn_topic_v2`: Fix smn topic update error [GH-48]
* `resource/huaweicloud_kms_key_v1`: Add default value of pending_days [GH-62]
* `all resources`: Expose real error message of BadRequest error [GH-63]
* `resource/huaweicloud_sfs_file_system_v2`: Suppress sfs system metadata [GH-64]

## 1.2.0 (September 21, 2018)

FEATURES:

* **New Data Source:** `huaweicloud_vpc_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_peering_connection_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_route_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_route_ids_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_subnet_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_vpc_subnet_ids_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Data Source:** `huaweicloud_rts_software_config_v1` ([#20](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/20))
* **New Data Source:** `huaweicloud_images_image_v2` ([#25](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/25))
* **New Resource:** `huaweicloud_vpc_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_peering_connection_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_peering_connection_accepter_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_route_v2` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_vpc_subnet_v1` ([#14](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/14))
* **New Resource:** `huaweicloud_rts_software_config_v1` ([#20](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/20))
* **New Resource:** `huaweicloud_images_image_v2` ([#25](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/25))
* **New Resource:** `huaweicloud_ces_alarmrule` ([#27](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/27))
* **New Resource:** `huaweicloud_as_configuration_v1` ([#29](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/29))
* **New Resource:** `huaweicloud_as_group_v1` ([#30](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/30))
* **New Resource:** `huaweicloud_as_policy_v1` ([#31](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/31))
* **New Resource:** `huaweicloud_cce_cluster_v3` [GH-19]
* **New Resource:** `huaweicloud_cce_node_v3` [GH-19]

ENHANCEMENTS:

* provider: Add AK/SK authentication support ([#33](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/33))

## 1.1.0 (July 20, 2018)

FEATURES:

* **New Data Source:** `huaweicloud_sfs_file_system_v2` ([#9](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/9))
* **New Data Source:** `huaweicloud_rts_stack_v1` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/10))
* **New Data Source:** `huaweicloud_rts_stack_resource_v1` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/10))
* **New Resource:** `huaweicloud_iam_agency_v3` ([#7](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/7))
* **New Resource:** `huaweicloud_sfs_file_system_v2` ([#9](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/9))
* **New Resource:** `huaweicloud_rts_stack_v1` ([#10](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/10))
* **New Resource:** `huaweicloud_iam_agency_v3` ([#16](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/16))

ENHANCEMENTS:

* resource/huaweicloud_dns_recordset_v2: Add `PTR` type ([#12](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/12))

BUG FIXES:

* provider: Create only one token ([#5](https://github.com/terraform-providers/terraform-provider-huaweicloud/issues/5))

## 1.0.0 (July 13, 2018)

FEATURES:

* **New Data Source:** `huaweicloud_networking_network_v2`
* **New Data Source:** `huaweicloud_networking_subnet_v2`
* **New Data Source:** `huaweicloud_networking_secgroup_v2`
* **New Data Source:** `huaweicloud_s3_bucket_object`
* **New Data Source:** `huaweicloud_kms_key_v1`
* **New Data Source:** `huaweicloud_kms_data_key_v1`
* **New Data Source:** `huaweicloud_rds_flavors_v1`
* **New Resource:** `huaweicloud_blockstorage_volume_v2`
* **New Resource:** `huaweicloud_compute_instance_v2`
* **New Resource:** `huaweicloud_compute_keypair_v2`
* **New Resource:** `huaweicloud_compute_secgroup_v2`
* **New Resource:** `huaweicloud_compute_servergroup_v2`
* **New Resource:** `huaweicloud_compute_floatingip_v2`
* **New Resource:** `huaweicloud_compute_floatingip_associate_v2`
* **New Resource:** `huaweicloud_compute_volume_attach_v2`
* **New Resource:** `huaweicloud_dns_recordset_v2`
* **New Resource:** `huaweicloud_dns_zone_v2`
* **New Resource:** `huaweicloud_fw_firewall_group_v2`
* **New Resource:** `huaweicloud_fw_policy_v2`
* **New Resource:** `huaweicloud_fw_rule_v2`
* **New Resource:** `huaweicloud_kms_key_v1`
* **New Resource:** `huaweicloud_elb_loadbalancer`
* **New Resource:** `huaweicloud_elb_listener`
* **New Resource:** `huaweicloud_elb_healthcheck`
* **New Resource:** `huaweicloud_lb_loadbalancer_v2`
* **New Resource:** `huaweicloud_lb_listener_v2`
* **New Resource:** `huaweicloud_lb_pool_v2`
* **New Resource:** `huaweicloud_lb_member_v2`
* **New Resource:** `huaweicloud_lb_monitor_v2`
* **New Resource:** `huaweicloud_networking_network_v2`
* **New Resource:** `huaweicloud_networking_subnet_v2`
* **New Resource:** `huaweicloud_networking_floatingip_v2`
* **New Resource:** `huaweicloud_networking_port_v2`
* **New Resource:** `huaweicloud_networking_router_v2`
* **New Resource:** `huaweicloud_networking_router_interface_v2`
* **New Resource:** `huaweicloud_networking_router_route_v2`
* **New Resource:** `huaweicloud_networking_secgroup_v2`
* **New Resource:** `huaweicloud_networking_secgroup_rule_v2`
* **New Resource:** `huaweicloud_s3_bucket`
* **New Resource:** `huaweicloud_s3_bucket_policy`
* **New Resource:** `huaweicloud_s3_bucket_object`
* **New Resource:** `huaweicloud_smn_topic_v2`
* **New Resource:** `huaweicloud_smn_subscription_v2`
* **New Resource:** `huaweicloud_rds_instance_v1`
* **New Resource:** `huaweicloud_nat_gateway_v2`
* **New Resource:** `huaweicloud_nat_snat_rule_v2`
