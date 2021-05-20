package huaweicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cci/v1/persistentvolumeclaims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var (
	fsType = map[string]string{
		"sas":             "ext4",
		"ssd":             "ext4",
		"sata":            "ext4",
		"nfs-rw":          "nfs",
		"efs-performance": "nfs",
		"efs-standard":    "nfs",
		"obs":             "obs",
	}
	volumeTypeForList = map[string]string{
		"sas":             "bs",
		"ssd":             "bs",
		"sata":            "bs",
		"obs":             "obs",
		"nfs-rw":          "nfs",
		"efs-performance": "efs",
		"efs-standard":    "efs",
	}
)

type StateRefresh struct {
	Pending      []string
	Target       []string
	Delay        time.Duration
	Timeout      time.Duration
	PollInterval time.Duration
}

func ResourceCCIPersistentVolumeClaimV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCIPersistentVolumeClaimV1Create,
		Read:   resourceCCIPersistentVolumeClaimV1Read,
		Delete: resourceCCIPersistentVolumeClaimV1Delete,
		Importer: &schema.ResourceImporter{
			State: resourceCCIPvcImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$"),
					"The name consists of 1 to 63 characters, including lowercase letters, digits and hyphens, "+
						"and must start and end with lowercase letters and digits"),
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"device_mount_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "sas",
				ValidateFunc: validation.StringInSlice([]string{
					"sas", "ssd", "sata", "obs", "nfs-rw", "efs-standard", "efs-performance",
				}, false),
			},
			"access_modes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildCCIPersistentVolumeClaimV1CreateParams(d *schema.ResourceData) (persistentvolumeclaims.CreateOpts, error) {
	createOpts := persistentvolumeclaims.CreateOpts{
		Kind:       "PersistentVolumeClaim",
		ApiVersion: "v1",
	}
	volumeType := d.Get("volume_type").(string)
	fsType, ok := fsType[volumeType]
	if !ok {
		return createOpts, fmt.Errorf("The volume type (%s) is not available", volumeType)
	}
	createOpts.Metadata = persistentvolumeclaims.Metadata{
		Namespace: d.Get("namespace").(string),
		Name:      d.Get("name").(string),
		Annotations: &persistentvolumeclaims.Annotations{
			FsType:          fsType,
			VolumeID:        d.Get("volume_id").(string),
			DeviceMountPath: d.Get("device_mount_path").(string),
		},
	}
	createOpts.Spec = persistentvolumeclaims.Spec{
		StorageClassName: volumeType,
		Resources: persistentvolumeclaims.ResourceRequirement{
			Requests: &persistentvolumeclaims.ResourceName{
				Storage: fmt.Sprintf("%dGi", d.Get("volume_size").(int)),
			},
		},
	}

	return createOpts, nil
}

func resourceCCIPersistentVolumeClaimV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CciV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCI client: %s", err)
	}

	createOpts, err := buildCCIPersistentVolumeClaimV1CreateParams(d)
	if err != nil {
		return fmt.Errorf("Unable to build createOpts of the Persistent Volume Claim: %s", err)
	}
	ns := d.Get("namespace").(string)
	create, err := persistentvolumeclaims.Create(client, createOpts, ns).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCI PVC: %s", err)
	}
	d.SetId(create.Metadata.UID)
	stateRef := StateRefresh{
		Pending:      []string{"Pending", "Bound"},
		Target:       []string{"Bound"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	if err := waitForCCIPersistentVolumeClaimStateRefresh(d, client, ns, stateRef); err != nil {
		return fmt.Errorf("Create the specifies PVC (%s) timed out: %s", d.Id(), err)
	}

	return resourceCCIPersistentVolumeClaimV1Read(d, meta)
}

func saveCCIPersistentVolumeClaimV1State(d *schema.ResourceData, resp *persistentvolumeclaims.ListResp) error {
	specResp := &resp.PersistentVolume.Spec
	metadata := &resp.PersistentVolumeClaim.Metadata
	// The volume size format is 'xGi'.
	regex := regexp.MustCompile("^([1-9][0-9]*)Gi$")
	storages := regex.FindStringSubmatch(specResp.Capacity.Storage)
	if len(storages) < 2 {
		return fmt.Errorf("The response of volume capacity is not a valid number, need 'xGi', but got '%s'",
			specResp.Capacity.Storage)
	}
	volumeSize, err := strconv.Atoi(storages[1])
	if err != nil {
		return fmt.Errorf("The input capacity (%s) is not a number: %s", storages[1], err)
	}

	mErr := multierror.Append(nil,
		d.Set("namespace", metadata.Namespace),
		d.Set("name", metadata.Name),
		d.Set("volume_id", specResp.FlexVolume.Options.VolumeID),
		d.Set("volume_size", volumeSize),
		d.Set("volume_type", specResp.StorageClassName),
		d.Set("device_mount_path", specResp.FlexVolume.Options.DeviceMountPath),
		d.Set("access_modes", specResp.AccessModes),
		d.Set("status", resp.PersistentVolumeClaim.Status.Phase),
		d.Set("creation_timestamp", metadata.CreationTimestamp),
		d.Set("enable", metadata.Enable),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

func resourceCCIPersistentVolumeClaimV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	client, err := config.CciV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	volumeType := d.Get("volume_type").(string)
	response, err := getPvcInfoFromServer(client, ns, volumeType, d.Id())
	if err != nil {
		return CheckDeleted(d, err, "Error getting the specifies PVC form server")
	}
	if response != nil {
		d.Set("region", region)
		if err := saveCCIPersistentVolumeClaimV1State(d, response); err != nil {
			return fmt.Errorf("Error saving the specifies PVC (%s) to state: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceCCIPersistentVolumeClaimV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CciV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCI Client: %s", err)
	}

	name := d.Get("name").(string)
	ns := d.Get("namespace").(string)
	_, err = persistentvolumeclaims.Delete(client, ns, name).Extract()
	if err != nil {
		return fmt.Errorf("Error deleting the specifies PVC (%s): %s", d.Id(), err)
	}

	stateRef := StateRefresh{
		Pending:      []string{"Bound"},
		Target:       []string{"DELETED"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        3 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if err := waitForCCIPersistentVolumeClaimStateRefresh(d, client, ns, stateRef); err != nil {
		return fmt.Errorf("Delete the specifies PVC (%s) timed out: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func waitForCCIPersistentVolumeClaimStateRefresh(d *schema.ResourceData, client *golangsdk.ServiceClient,
	ns string, s StateRefresh) error {
	stateConf := &resource.StateChangeConf{
		Pending:      s.Pending,
		Target:       s.Target,
		Refresh:      pvcStateRefreshFunc(d, client, ns),
		Timeout:      s.Timeout,
		Delay:        s.Delay,
		PollInterval: s.PollInterval,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Waiting for the status of the PVC (%s) to complete timeout: %s", d.Id(), err)
	}
	return nil
}

func pvcStateRefreshFunc(d *schema.ResourceData, client *golangsdk.ServiceClient,
	ns string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		volumeType := d.Get("volume_type").(string)
		response, err := getPvcInfoFromServer(client, ns, volumeType, d.Id())
		if err != nil {
			return response, "ERROR", nil
		}
		if response != nil {
			return response, response.PersistentVolumeClaim.Status.Phase, nil
		}
		return response, "DELETED", nil
	}
}

func getPvcInfoFromServer(client *golangsdk.ServiceClient, ns, volumeType,
	pvcId string) (*persistentvolumeclaims.ListResp, error) {
	var response *persistentvolumeclaims.ListResp
	// If the storage of listOpts is not set, the list method will search for all PVCs of evs type.
	storageType, ok := volumeTypeForList[volumeType]
	if !ok {
		return response, fmt.Errorf("The volume type (%s) is not available", volumeType)
	}
	listOpts := persistentvolumeclaims.ListOpts{
		StorageType: storageType,
	}
	pages, err := persistentvolumeclaims.List(client, listOpts, ns).AllPages()
	if err != nil {
		return response, fmt.Errorf("Error finding the PVCs of type %s: %s", storageType, err)
	}
	responses, err := persistentvolumeclaims.ExtractPersistentVolumeClaims(pages)
	if err != nil {
		return response, fmt.Errorf("Error retrieving HuaweiCloud CCI PVCs: %s", err)
	}
	for _, v := range responses {
		if v.PersistentVolumeClaim.Metadata.UID == pvcId {
			response = new(persistentvolumeclaims.ListResp)
			response = &v
			break
		}
	}
	return response, nil
}

func resourceCCIPvcImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid format specified for CCI PVC, must be <namespace>/<volume type>/<pvc id>")
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("namespace", parts[0]),
		d.Set("volume_type", parts[1]),
	)
	if mErr.ErrorOrNil() != nil {
		return []*schema.ResourceData{d}, mErr
	}
	return []*schema.ResourceData{d}, nil
}
