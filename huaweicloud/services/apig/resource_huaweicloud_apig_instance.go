package apig

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/instances"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type Edition string      // The edition of the dedicated instance.
type ProviderType string // The type of the loadbalancer provider.

const (
	// IPv4 Editions
	EditionBasic        Edition = "BASIC"        // Basic Edition instance.
	EditionProfessional Edition = "PROFESSIONAL" // Professional Edition instance.
	EditionEnterprise   Edition = "ENTERPRISE"   // Enterprise Edition instance.
	EditionPlatinum     Edition = "PLATINUM"     // Platinum Edition instance.
	// IPv6 Editions
	Ipv6EditionBasic        Edition = "BASIC_IPv6"        // IPv6 instance of the Basic Edition.
	Ipv6EditionProfessional Edition = "PROFESSIONAL_IPv6" // IPv6 instance of the Professional Edition.
	Ipv6EditionEnterprise   Edition = "ENTERPRISE_IPv6"   // IPv6 instance of the Enterprise Edition.
	Ipv6EditionPlatinum     Edition = "PLATINUM_IPv6"     // IPv6 instance of the Platinum Edition.

	ProviderTypeLvs ProviderType = "lvs" // Linux virtual server.
	ProviderTypeElb ProviderType = "elb" // Elastic load balance.

	enableFeature  bool = true
	disableFeature bool = false
)

func ResourceApigInstanceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to create the dedicated instance resource.`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^([\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z-_0-9]*)$"),
						"The name can only contain letters, digits, hyphens (-) and underscore (_), and must start "+
							"with a letter."),
					validation.StringLenBetween(3, 64),
				),
				Description: `The name of the dedicated instance.`,
			},
			"edition": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(EditionBasic),
					string(EditionProfessional),
					string(EditionEnterprise),
					string(EditionPlatinum),
					string(Ipv6EditionBasic),
					string(Ipv6EditionProfessional),
					string(Ipv6EditionEnterprise),
					string(Ipv6EditionPlatinum),
				}, false),
				Description: `The edition of the dedicated instance.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the VPC used to create the dedicated instance.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the VPC subnet used to create the dedicated instance.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the security group to which the dedicated instance belongs to.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `schema: Required; The name list of availability zones for the dedicated instance.`,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[^<>]*$"),
						"The description cannot contain the angle brackets (< and >)."),
					validation.StringLenBetween(0, 255),
				),
				Description: `The description of the dedicated instance.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The enterprise project ID to which the dedicated instance belongs.`,
			},
			"bandwidth_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 2000),
				Description:  `The egress bandwidth size of the dedicated instance.`,
			},
			"eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The EIP ID associated with the dedicated instance.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether public access with an IPv6 address is supported.`,
			},
			"loadbalancer_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ProviderTypeLvs), string(ProviderTypeElb),
				}, false),
				Description: `The type of loadbalancer provider used by the instance.`,
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(02|06|10|14|18|22):00:00$`),
					"The start-time format of maintenance window is not 'xx:00:00' or "+
						"the hour is not 02, 06, 10, 14, 18 or 22."),
				Description: `The start time of the maintenance time window.`,
			},
			"feature": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the instance feature.",
						},
						"config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The configuration detail (JSON format) of the instance feature.",
						},
					},
				},
				Description: "The custom feature configuration.",
			},
			// Attributes
			"maintain_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `End time of the maintenance time window, 4-hour difference between the start time and end time.`,
			},
			"ingress_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ingress EIP address.`,
			},
			"vpc_ingress_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ingress private IP address of the VPC.`,
			},
			"egress_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The egress (NAT) public IP address.`,
			},
			"supported_features": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The supported features of the dedicated instance.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time when the dedicated instance is created, in RFC-3339 format.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Status of the dedicated instance.`,
			},
			// Deprecated arguments
			"available_zones": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `schema: Deprecated; The name list of availability zones for the dedicated instance.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use 'created_at' instead",
				Description: `schema: Deprecated; Time when the dedicated instance is created.`,
			},
		},
	}
}

func buildMaintainEndTime(maintainStart string) (string, error) {
	result := regexp.MustCompile("^(02|06|10|14|18|22):00:00$").FindStringSubmatch(maintainStart)
	if len(result) < 2 {
		return "", fmt.Errorf("the hour is missing")
	}
	num, err := strconv.Atoi(result[1])
	if err != nil {
		return "", fmt.Errorf("the number (%s) cannot be converted to string", result[1])
	}
	return fmt.Sprintf("%02d:00:00", (num+4)%24), nil
}

func buildInstanceAvailabilityZones(d *schema.ResourceData) ([]string, error) {
	if v, ok := d.GetOk("availability_zones"); ok {
		return utils.ExpandToStringList(v.([]interface{})), nil
	}

	// When 'availability_zones' is omitted, the deprecated parameter 'available_zones' is used.
	if v, ok := d.GetOk("available_zones"); ok {
		return utils.ExpandToStringList(v.([]interface{})), nil
	}

	return nil, fmt.Errorf("the parameter 'availability_zones' must be specified.")
}

func buildInstanceCreateOpts(d *schema.ResourceData, config *config.Config) (instances.CreateOpts, error) {
	result := instances.CreateOpts{
		Name:                 d.Get("name").(string),
		Edition:              d.Get("edition").(string),
		VpcId:                d.Get("vpc_id").(string),
		SubnetId:             d.Get("subnet_id").(string),
		SecurityGroupId:      d.Get("security_group_id").(string),
		Description:          d.Get("description").(string),
		EipId:                d.Get("eip_id").(string),
		BandwidthSize:        d.Get("bandwidth_size").(int), // Bandwidth 0 means turn off the egress access.
		EnterpriseProjectId:  common.GetEnterpriseProjectID(d, config),
		Ipv6Enable:           d.Get("ipv6_enable").(bool),
		LoadbalancerProvider: d.Get("loadbalancer_provider").(string),
	}

	azList, err := buildInstanceAvailabilityZones(d)
	if err != nil {
		return result, err
	}
	result.AvailableZoneIds = azList

	if v, ok := d.GetOk("maintain_begin"); ok {
		startTime := v.(string)
		result.MaintainBegin = startTime
		endTime, err := buildMaintainEndTime(startTime)
		if err != nil {
			return result, err
		}
		result.MaintainEnd = endTime
	}

	log.Printf("[DEBUG] Create options of the dedicated instance is: %#v", result)
	return result, nil
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opts, err := buildInstanceCreateOpts(d, config)
	if err != nil {
		return diag.Errorf("error creating the dedicated instance options: %s", err)
	}
	log.Printf("[DEBUG] The CreateOpts of the dedicated instance is: %#v", opts)

	resp, err := instances.Create(client, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating the dedicated instance: %s", err)
	}
	d.SetId(resp.Id)

	instanceId := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      InstanceStateRefreshFunc(client, instanceId, []string{"Running"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the dedicated instance (%s) to become running: %s", instanceId, err)
	}

	err = updateInstanceFeatures(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceInstanceRead(ctx, d, meta)
}

// parseInstanceAvailabilityZones is a method that used to convert the string returned by the API which contains
// brackets ([ and ]) and space into a list of strings (available_zone code) and save to state.
func parseInstanceAvailabilityZones(azStr string) []string {
	codesStr := strings.TrimLeft(azStr, "[")
	codesStr = strings.TrimRight(codesStr, "]")
	codesStr = strings.ReplaceAll(codesStr, " ", "")

	return strings.Split(codesStr, ",")
}

// The response of ingress acess does not contain EIP ID, just the IP address.
func parseInstanceIngressAccess(config *config.Config, region, publicAddress string) (*string, error) {
	if publicAddress == "" {
		return nil, nil
	}

	client, err := config.NetworkingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
	}

	opt := eips.ListOpts{
		PublicIp:            []string{publicAddress},
		EnterpriseProjectId: "all_granted_eps",
	}
	allPages, err := eips.List(client, opt).AllPages()
	if err != nil {
		return nil, err
	}
	publicIps, err := eips.ExtractPublicIPs(allPages)
	if err != nil {
		return nil, err
	}
	if len(publicIps) > 0 {
		return &publicIps[0].ID, nil
	}

	log.Printf("[WARN] The instance does not synchronize EIP information, got (%s), but not found on the server",
		publicAddress)
	return nil, nil
}

func parseInstanceIpv6Enable(ipv6Address string) bool {
	return ipv6Address != ""
}

func parseInstanceFeatures(d *schema.ResourceData, features []instances.Feature) []map[string]interface{} {
	featuresConfig := d.Get("feature").(*schema.Set)
	if featuresConfig.Len() < 1 {
		return nil
	}

	result := make([]map[string]interface{}, featuresConfig.Len())
	for i, val := range featuresConfig.List() {
		config := val.(map[string]interface{})
		for _, feature := range features {
			if config["name"] == feature.Name {
				result[i] = map[string]interface{}{
					"name":   feature.Name,
					"config": feature.Config,
				}
			}
		}
	}

	return result
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Id()
	resp, err := instances.Get(client, instanceId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting instance (%s) details form server", instanceId))
	}
	log.Printf("[DEBUG] Retrieved the dedicated instance (%s): %#v", instanceId, resp)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("edition", resp.Edition),
		d.Set("vpc_id", resp.VpcId),
		d.Set("subnet_id", resp.SubnetId),
		d.Set("security_group_id", resp.SecurityGroupId),
		d.Set("description", resp.Description),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("bandwidth_size", resp.BandwidthSize),
		d.Set("ipv6_enable", parseInstanceIpv6Enable(resp.Ipv6IngressEipAddress)),
		d.Set("loadbalancer_provider", resp.LoadbalancerProvider),
		d.Set("availability_zones", parseInstanceAvailabilityZones(resp.AvailableZoneIds)),
		d.Set("maintain_begin", resp.MaintainBegin),
		// Attributes
		d.Set("maintain_end", resp.MaintainEnd),
		d.Set("ingress_address", resp.Ipv4IngressEipAddress),
		d.Set("vpc_ingress_address", resp.Ipv4VpcIngressAddress),
		d.Set("egress_address", resp.Ipv4EgressAddress),
		d.Set("supported_features", resp.SupportedFeatures),
		d.Set("status", resp.Status),
		d.Set("created_at", utils.FormatTimeStampRFC3339(resp.CreateTimestamp, false)),
		// Deprecated
		d.Set("create_time", utils.FormatTimeStampRFC3339(resp.CreateTimestamp, false)),
	)

	if eipId, err := parseInstanceIngressAccess(config, region, resp.Ipv4IngressEipAddress); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		mErr = multierror.Append(d.Set("eip_id", eipId))
	}

	// Expand the limit value, because of the parameter offset does not take effect.
	listOpts := instances.ListFeaturesOpts{Limit: 500}
	if features, err := instances.ListFeatures(client, instanceId, listOpts); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("feature", parseInstanceFeatures(d, features)))
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving resource fields of the dedicated instance: %s", mErr)
	}
	return nil
}

func buildInstanceUpdateOpts(d *schema.ResourceData) (instances.UpdateOpts, error) {
	result := instances.UpdateOpts{}
	if d.HasChange("name") {
		result.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		result.Description = utils.String(d.Get("description").(string))
	}
	if d.HasChange("security_group_id") {
		result.SecurityGroupId = d.Get("security_group_id").(string)
	}
	if d.HasChange("maintain_begin") {
		startTime := d.Get("maintain_begin").(string)
		result.MaintainBegin = startTime
		endTime, err := buildMaintainEndTime(startTime)
		if err != nil {
			return result, err
		}
		result.MaintainEnd = endTime
	}

	log.Printf("[DEBUG] Update options of the dedicated instance is: %#v", result)
	return result, nil
}

func updateApigInstanceEgressAccess(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldVal, newVal := d.GetChange("bandwidth_size")
	// Enable the egress access.
	if oldVal.(int) == 0 {
		size := d.Get("bandwidth_size").(int)
		opts := instances.EgressAccessOpts{
			BandwidthSize: strconv.Itoa(size),
		}
		egress, err := instances.EnableEgressAccess(client, d.Id(), opts).Extract()
		if err != nil {
			return fmt.Errorf("unable to enable egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
		}
		if egress.BandwidthSize != size {
			return fmt.Errorf("the egress bandwidth size change failed, want '%d', but '%d'", size, egress.BandwidthSize)
		}
	}
	// Disable the egress access.
	if newVal.(int) == 0 {
		err := instances.DisableEgressAccess(client, d.Id()).ExtractErr()
		if err != nil {
			return fmt.Errorf("unable to disable egress bandwidth of the dedicated instance (%s)", d.Id())
		}
		return nil
	}
	// Update the egress nat.
	size := d.Get("bandwidth_size").(int)
	opts := instances.EgressAccessOpts{
		BandwidthSize: strconv.Itoa(size),
	}
	egress, err := instances.UpdateEgressBandwidth(client, d.Id(), opts).Extract()
	if err != nil {
		return fmt.Errorf("unable to update egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
	}
	if egress.BandwidthSize != size {
		return fmt.Errorf("the egress bandwidth size change failed, want '%d', but '%d'", size, egress.BandwidthSize)
	}
	return nil
}

func updateInstanceIngressAccess(d *schema.ResourceData, client *golangsdk.ServiceClient) (err error) {
	oldVal, newVal := d.GetChange("eip_id")
	// Disable the ingress access.
	// The update logic is to disable first and then enable. Update means thar both oldVal and newVal exist.
	if oldVal.(string) != "" {
		err = instances.DisableIngressAccess(client, d.Id()).ExtractErr()
		if err != nil || newVal.(string) == "" {
			return
		}
	}
	// Enable the ingress access.
	updateOpts := instances.IngressAccessOpts{
		EipId: d.Get("eip_id").(string),
	}
	_, err = instances.EnableIngressAccess(client, d.Id(), updateOpts).Extract()
	return
}

func updateInstanceFeatures(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		oldVal, newVal = d.GetChange("feature")
		addRaws        = newVal.(*schema.Set).Difference(oldVal.(*schema.Set))
		removeRaws     = oldVal.(*schema.Set).Difference(newVal.(*schema.Set))
		instanceId     = d.Id()
	)
	// Disable the ingress access.
	// The update logic is to disable first and then enable. Update means thar both oldVal and newVal exist.
	for _, val := range removeRaws.List() {
		feature := val.(map[string]interface{})
		opts := instances.FeatureOpts{
			Name:   feature["name"].(string),
			Enable: utils.Bool(disableFeature),
			Config: feature["config"].(string),
		}
		_, err := instances.UpdateFeature(client, instanceId, opts)
		if err != nil {
			return err
		}
	}
	// Disable the ingress access.
	// The update logic is to disable first and then enable. Update means thar both oldVal and newVal exist.
	for _, val := range addRaws.List() {
		feature := val.(map[string]interface{})
		opts := instances.FeatureOpts{
			Name:   feature["name"].(string),
			Enable: utils.Bool(enableFeature),
			Config: feature["config"].(string),
		}
		_, err := instances.UpdateFeature(client, instanceId, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	// Update egress access
	if d.HasChange("bandwidth_size") {
		if err = updateApigInstanceEgressAccess(d, client); err != nil {
			return diag.Errorf("update egress access failed: %s", err)
		}
	}
	// Update ingerss access
	if d.HasChange("eip_id") {
		if err = updateInstanceIngressAccess(d, client); err != nil {
			return diag.Errorf("update ingress access failed: %s", err)
		}
	}
	// Update feature configuration
	if d.HasChange("feature") {
		if err = updateInstanceFeatures(d, client); err != nil {
			return diag.Errorf("update feature configuration failed: %s", err)
		}
	}
	// Update instance name, maintain window, description and security group ID.
	updateOpts, err := buildInstanceUpdateOpts(d)
	if err != nil {
		return diag.Errorf("unable to get the update options of the dedicated instance: %s", err)
	}
	if updateOpts != (instances.UpdateOpts{}) {
		_, err = instances.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating the dedicated instance: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      InstanceStateRefreshFunc(client, d.Id(), []string{"Running"}),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        20 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceInstanceRead(ctx, d, meta)
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	if err = instances.Delete(client, d.Id()).ExtractErr(); err != nil {
		return diag.Errorf("error deleting the dedicated instance (%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      InstanceStateRefreshFunc(client, d.Id(), nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func InstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := instances.Get(client, instanceId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}
			return resp, "", err
		}

		if utils.StrSliceContains([]string{"CreateFail", "InitingFailed", "RegisterFailed", "InstallFailed",
			"UpdateFailed", "RollbackFailed", "UnRegisterFailed", "DeleteFailed"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpect status (%s)", resp.Status)
		}

		if utils.StrSliceContains(targets, resp.Status) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}
