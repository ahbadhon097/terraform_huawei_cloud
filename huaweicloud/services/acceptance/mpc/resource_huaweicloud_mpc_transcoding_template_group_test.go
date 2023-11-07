package mpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	mpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/mpc/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getTemplateGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.HcMpcV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating MPC client: %s", err)
	}

	resp, err := c.ListTemplateGroup(&mpc.ListTemplateGroupRequest{GroupId: &[]string{state.Primary.ID}})
	if err != nil {
		return nil, fmt.Errorf("error retrieving MPC transcoding template group: %s", err)
	}

	templateGroupList := *resp.TemplateGroupList

	if len(templateGroupList) == 0 {
		return nil, fmt.Errorf("unable to retrieve MPC transcoding template group: %s", state.Primary.ID)
	}

	return templateGroupList[0], nil
}

func TestAccTranscodingTemplateGroup_basic(t *testing.T) {
	var templateGroup mpc.TemplateGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := rName + "-update"
	resourceName := "huaweicloud_mpc_transcoding_template_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&templateGroup,
		getTemplateGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplateGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "2"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.codec", "2"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.profile", "4"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.output_policy", "transcode"),
					resource.TestCheckResourceAttr(resourceName, "videos.0.width", "1920"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testTranscodingTemplateGroup_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "false"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "2"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.profile", "3"),
					resource.TestCheckResourceAttr(resourceName, "videos.0.width", "3840"),
					resource.TestCheckResourceAttr(resourceName, "videos.1.width", "2560"),
				),
			},
		},
	})
}

func testTranscodingTemplateGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_mpc_transcoding_template_group" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 1

  audio {
    bitrate       = 0
    channels      = 2
    codec         = 2
    output_policy = "transcode"
    sample_rate   = 1
  }

  video_common {
    max_consecutive_bframes = 7
    black_bar_removal       = 0
    codec                   = 2
    fps                     = 0
    level                   = 15
    max_iframes_interval    = 5
    output_policy           = "transcode"
    quality                 = 1
    profile                 = 4
  }

  videos {
    width   = 1920
    height  = 1080
    bitrate = 0
  }
}
`, rName)
}

func testTranscodingTemplateGroup_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_mpc_transcoding_template_group" "test" {
  name                  = "%s"
  low_bitrate_hd        = false
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 2

  audio {
    bitrate       = 0
    channels      = 2
    codec         = 1
    output_policy = "transcode"
    sample_rate   = 1
  }

  video_common {
    max_consecutive_bframes = 7
    black_bar_removal       = 0
    codec                   = 1
    fps                     = 0
    level                   = 15
    max_iframes_interval    = 5
    output_policy           = "transcode"
    quality                 = 1
    profile                 = 3
  }

  videos {
    width   = 3840
    height  = 2160
    bitrate = 0
  }

  videos {
    width   = 2560
    height  = 1440
    bitrate = 0
  }
}
`, rName)
}
