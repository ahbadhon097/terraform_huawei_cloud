package huaweicloud

import (
	"sort"
	"strconv"
	"time"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var iamgeValidSortKeys = []string{
	"name", "container_format", "disk_format", "status", "id", "size",
}
var imageValidVisibilities = []string{
	"public", "private", "community", "shared",
}

func DataSourceImagesImageV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImagesImageV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(imageValidVisibilities, false),
			},

			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sort_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "name",
				ValidateFunc: validation.StringInSlice(iamgeValidSortKeys, false),
			},

			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "asc",
				ValidateFunc: validation.StringInSlice([]string{
					"asc", "desc",
				}, false),
			},

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// Deprecated values
			"size_min": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "size_min is deprecated",
			},
			"size_max": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "size_max is deprecated",
			},

			// Computed values
			"container_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_ram_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceImagesImageV2Read performs the image lookup.
func dataSourceImagesImageV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	imageClient, err := config.ImageV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud image client: %s", err)
	}

	listOpts := cloudimages.ListOpts{
		Name:       d.Get("name").(string),
		Visibility: d.Get("visibility").(string),
		Owner:      d.Get("owner").(string),
		Status:     "active",
		SortKey:    d.Get("sort_key").(string),
		SortDir:    d.Get("sort_direction").(string),
		Tag:        d.Get("tag").(string),
	}

	logp.Printf("[DEBUG] List Options: %#v", listOpts)

	var image cloudimages.Image
	allPages, err := cloudimages.List(imageClient, listOpts).AllPages()
	if err != nil {
		return fmtp.Errorf("Unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allImages) > 1 {
		recent := d.Get("most_recent").(bool)
		logp.Printf("[DEBUG] Multiple results found and `most_recent` is set to: %t", recent)
		if recent {
			image = mostRecentImage(allImages)
		} else {
			logp.Printf("[DEBUG] Multiple results found: %#v", allImages)
			return fmtp.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		image = allImages[0]
	}

	logp.Printf("[DEBUG] Single Image found: %s", image.ID)
	return dataSourceImagesImageV2Attributes(d, &image)
}

// dataSourceImagesImageV2Attributes populates the fields of an Image resource.
func dataSourceImagesImageV2Attributes(d *schema.ResourceData, image *cloudimages.Image) error {
	logp.Printf("[DEBUG] huaweicloud_images_image details: %#v", image)

	d.SetId(image.ID)
	d.Set("name", image.Name)
	d.Set("container_format", image.ContainerFormat)
	d.Set("disk_format", image.DiskFormat)
	d.Set("min_disk_gb", image.MinDisk)
	d.Set("min_ram_mb", image.MinRam)
	d.Set("owner", image.Owner)
	d.Set("protected", image.Protected)
	d.Set("visibility", image.Visibility)
	d.Set("checksum", image.Checksum)
	d.Set("file", image.File)
	d.Set("schema", image.Schema)
	d.Set("status", image.Status)
	d.Set("created_at", image.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", image.UpdatedAt.Format(time.RFC3339))

	if size, err := strconv.Atoi(image.ImageSize); err == nil {
		d.Set("size_bytes", size)
	}

	return nil
}

type imageSort []cloudimages.Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	itime := a[i].CreatedAt
	jtime := a[j].CreatedAt
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent Image out of a slice of images.
func mostRecentImage(images []cloudimages.Image) cloudimages.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}
