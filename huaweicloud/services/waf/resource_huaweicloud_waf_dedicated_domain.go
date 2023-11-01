package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf/v1/certificates"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"
	domains "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	protectStatusEnable = 1
)

func ResourceWafDedicatedDomainV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedDomainV1Create,
		ReadContext:   resourceWafDedicatedDomainV1Read,
		UpdateContext: resourceWafDedicatedDomainV1Update,
		DeleteContext: resourceWafDedicatedDomainV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 80,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
						},
						"server_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
						},
						"address": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"port": {
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(0, 65535),
							Required:     true,
							ForceNew:     true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"keep_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"protect_status": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"tls": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TLS v1.0", "TLS v1.1", "TLS v1.2",
				}, false),
			},
			"cipher": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cipher_1",
					"cipher_2",
					"cipher_3",
					"cipher_4",
					"cipher_default",
				}, false),
			},
			"pci_3ds": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"tls", "cipher"},
			},
			"pci_dss": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"tls", "cipher"},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alarm_page": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"compliance_certification": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
			"traffic_identifier": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func getCertificateNameById(d *schema.ResourceData, meta interface{}) (string, error) {
	cfg := meta.(*config.Config)
	client, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return "", fmt.Errorf("error creating WAF client: %s", err)
	}

	if _, ok := d.GetOk("certificate_id"); ok {
		epsID := cfg.GetEnterpriseProjectID(d)
		certificateId := d.Get("certificate_id").(string)
		r, err := certificates.GetWithEpsID(client, certificateId, epsID).Extract()
		if err != nil {
			return "", fmt.Errorf("error retrieving WAF certificate name according ID: %s, error: %s", certificateId, err)
		}
		return r.Name, nil
	}
	return "", nil
}

func buildCreatePremiumHostOpts(d *schema.ResourceData, meta interface{}) (*domains.CreateOpts, error) {
	certName, err := getCertificateNameById(d, meta)
	if err != nil {
		return nil, err
	}

	servers := d.Get("server").([]interface{})
	serverOpts := make([]domains.Server, len(servers))
	for i, v := range servers {
		s := v.(map[string]interface{})
		serverOpts[i] = domains.Server{
			FrontProtocol: s["client_protocol"].(string),
			BackProtocol:  s["server_protocol"].(string),
			Address:       s["address"].(string),
			Port:          s["port"].(int),
			Type:          s["type"].(string),
			VpcId:         s["vpc_id"].(string),
		}
	}

	proxy := d.Get("proxy").(bool)
	opts := domains.CreateOpts{
		CertificateId:   d.Get("certificate_id").(string),
		CertificateName: certName,
		HostName:        d.Get("domain").(string),
		Proxy:           &proxy,
		PolicyId:        d.Get("policy_id").(string),
		Servers:         serverOpts,
	}
	return &opts, nil
}

func resourceWafDedicatedDomainV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	createOpts, err := buildCreatePremiumHostOpts(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	createOpts.EnterpriseProjectID = epsID
	domain, err := domains.Create(client, *createOpts)
	if err != nil {
		return diag.Errorf("error creating WAF dedicated domain: %s", err)
	}
	d.SetId(domain.Id)

	if d.Get("protect_status").(int) != protectStatusEnable {
		protectStatus := d.Get("protect_status").(int)
		_, err = domains.UpdateProtectStatusWithWpsID(client, protectStatus, d.Id(), epsID)
		if err != nil {
			// If the protection status update fails, it will be managed by terraform, and only print log.
			log.Printf("[ERROR] error updating the protection status of WAF dedicate domain %s: %s", d.Id(), err)
		}
	}

	if d.HasChanges("tls", "cipher", "pci_3ds", "pci_dss") {
		if err := updateWafDedicatedDomain(client, meta, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceWafDedicatedDomainV1Read(ctx, d, meta)
}

func updateWafDedicatedDomain(client *golangsdk.ServiceClient, meta interface{}, d *schema.ResourceData) error {
	conf := meta.(*config.Config)
	updateOpts := domains.UpdatePremiumHostOpts{
		Tls:                 d.Get("tls").(string),
		Cipher:              d.Get("cipher").(string),
		EnterpriseProjectID: conf.GetEnterpriseProjectID(d),
	}

	if d.HasChange("proxy") && !d.IsNewResource() {
		updateOpts.Proxy = utils.Bool(d.Get("proxy").(bool))
	}

	if d.HasChange("certificate_id") && !d.IsNewResource() {
		if v, ok := d.GetOk("certificate_id"); ok {
			certName, err := getCertificateNameById(d, meta)
			if err != nil {
				return err
			}
			updateOpts.CertificateName = certName
			updateOpts.CertificateId = v.(string)
		}
	}

	if d.HasChanges("pci_3ds", "pci_dss") {
		flag, err := getHostFlag(d)
		if err != nil {
			return err
		}
		updateOpts.Flag = flag
	}
	_, err := domains.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmt.Errorf("error updating WAF dedicated domain: %s", err)
	}
	return nil
}

func getHostFlag(d *schema.ResourceData) (*domains.Flag, error) {
	pci3ds := d.Get("pci_3ds").(bool)
	pciDss := d.Get("pci_dss").(bool)
	if !pci3ds && !pciDss {
		return nil, nil
	}

	// required tls="TLS v1.2" && cipher="cipher_2"
	if d.Get("tls").(string) != "TLS v1.2" || d.Get("cipher").(string) != "cipher_2" {
		return nil, fmt.Errorf("pci_3ds and pci_dss must be used together with tls and cipher. " +
			"Tls must be set to TLS v1.2, and cipher must be set to cipher_2")
	}
	return &domains.Flag{
		Pci3ds: strconv.FormatBool(pci3ds),
		PciDss: strconv.FormatBool(pciDss),
	}, nil
}

func buildDomainServerAttribute(domain *domains.PremiumHost) []map[string]interface{} {
	servers := make([]map[string]interface{}, 0, len(domain.Servers))
	for _, s := range domain.Servers {
		servers = append(servers, map[string]interface{}{
			"client_protocol": s.FrontProtocol,
			"server_protocol": s.BackProtocol,
			"address":         s.Address,
			"port":            s.Port,
			"type":            s.Type,
			"vpc_id":          s.VpcId,
		})
	}
	return servers
}

func buildComplianceCertificationAttribute(domain *domains.PremiumHost) map[string]interface{} {
	f := domain.Flag

	pciDss, _ := strconv.ParseBool(f["pci_dss"])
	pci3ds, _ := strconv.ParseBool(f["pci_3ds"])
	return map[string]interface{}{
		"pci_dss": pciDss,
		"pci_3ds": pci3ds,
	}
}

func buildTrafficIdentifierAttribute(domain *domains.PremiumHost) map[string]interface{} {
	t := domain.TrafficMark
	return map[string]interface{}{
		"ip_tag":      strings.Join(t.Sip, ","),
		"session_tag": t.Cookie,
		"user_tag":    t.Params,
	}
}

func buildAlarmPageAttribute(domain *domains.PremiumHost) map[string]interface{} {
	t := domain.BlockPage
	return map[string]interface{}{
		"template_name": t.Template,
		"redirect_url":  t.RedirectUrl,
	}
}

func resourceWafDedicatedDomainV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	dm, err := domains.GetWithEpsID(client, d.Id(), epsID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving WAF dedicated domain")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("domain", dm.HostName),
		d.Set("server", buildDomainServerAttribute(dm)),
		d.Set("certificate_id", dm.CertificateId),
		d.Set("certificate_name", dm.CertificateName),
		d.Set("policy_id", dm.PolicyId),
		d.Set("proxy", dm.Proxy),
		d.Set("protect_status", dm.ProtectStatus),
		d.Set("access_status", dm.AccessStatus),
		d.Set("protocol", dm.Protocol),
		d.Set("tls", dm.Tls),
		d.Set("cipher", dm.Cipher),
		d.Set("compliance_certification", buildComplianceCertificationAttribute(dm)),
		d.Set("traffic_identifier", buildTrafficIdentifierAttribute(dm)),
		d.Set("alarm_page", buildAlarmPageAttribute(dm)),
	)

	if dm.Flag["pci_3ds"] != "" {
		pci3ds, err := strconv.ParseBool(dm.Flag["pci_3ds"])
		if err != nil {
			log.Printf("[WARN] error parse bool pci 3ds, %s", err)
		}
		mErr = multierror.Append(mErr, d.Set("pci_3ds", pci3ds))
	}

	if dm.Flag["pci_dss"] != "" {
		pciDss, err := strconv.ParseBool(dm.Flag["pci_dss"])
		if err != nil {
			log.Printf("[WARN] error parse bool pci dss, %s", err)
		}
		mErr = multierror.Append(mErr, d.Set("pci_dss", pciDss))
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWafDedicatedDomainV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dedicatedClient, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if d.HasChanges("tls", "cipher", "proxy", "certificate_id", "pci_3ds", "pci_dss") {
		if err := updateWafDedicatedDomain(dedicatedClient, meta, d); err != nil {
			return diag.Errorf("error updating WAF dedicated domain: %s", err)
		}
	}

	if d.HasChanges("protect_status") {
		protectStatus := d.Get("protect_status").(int)
		epsID := cfg.GetEnterpriseProjectID(d)
		_, err = domains.UpdateProtectStatusWithWpsID(dedicatedClient, protectStatus, d.Id(), epsID)
		if err != nil {
			return diag.Errorf("error updating the protection status of WAF dedicate domain %s: %s",
				d.Id(), err)
		}
	}

	if d.HasChanges("policy_id") {
		oVal, nVal := d.GetChange("policy_id")
		policyId := nVal.(string)
		epsID := cfg.GetEnterpriseProjectID(d)
		updateHostsOpts := policies.UpdateHostsOpts{
			Hosts:               []string{d.Id()},
			EnterpriseProjectId: epsID,
		}
		log.Printf("[DEBUG] Bind WAF dedicated domain %s to policy %s", d.Id(), policyId)
		_, err = policies.UpdateHosts(wafClient, policyId, updateHostsOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating WAF policy hosts: %s", err)
		}

		// delete the old policy
		err = policies.DeleteWithEpsID(wafClient, oVal.(string), epsID).ExtractErr()
		if err != nil {
			// If other domains are using this policy, the deletion will fail.
			log.Printf("[WARN] error deleting WAF policy %s: %s", oVal.(string), err)
		}
	}
	return resourceWafDedicatedDomainV1Read(ctx, d, meta)
}

func resourceWafDedicatedDomainV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	log.Printf("[DEBUG] Delete WAF dedicated domain (keep_policy: %v).", d.Get("keep_policy"))
	keepPolicy := d.Get("keep_policy").(bool)
	epsID := cfg.GetEnterpriseProjectID(d)
	_, err = domains.DeleteWithEpsID(client, keepPolicy, d.Id(), epsID)
	if err != nil {
		return diag.Errorf("error deleting WAF dedicated domain: %s", err)
	}
	return nil
}
