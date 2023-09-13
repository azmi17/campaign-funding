package campaign

import "strings"

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {

	campaignFormatter := CampaignFormatter{}

	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	// check if campaign is doesn't exist and it will return an empty array
	campaignsFormatter := []CampaignFormatter{} // zero-value is an empty array
	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)                      // single of campaigns
		campaignsFormatter = append(campaignsFormatter, campaignFormatter) // total of all campaigns w appended
	}

	return campaignsFormatter
}

type CampaignDetailFormatter struct {
	ID               int                      `json:"id"`
	Name             string                   `json:"name"`
	ShortDescription string                   `json:"short_description"`
	Description      string                   `json:"description"`
	ImageURL         string                   `json:"img_url"`
	GoalAmount       int                      `json:"goal_amount"`
	CurrentAmount    int                      `json:"current_amount"`
	BackerCount      int                      `json:"backer_count"`
	UserID           int                      `json:"user_id"`
	Slug             string                   `json:"slug"`
	Perks            []string                 `json:"perks"`
	User             CampaignUserFormatter    `json:"user"`
	Images           []CampaignImageFormatter `json:"images"`
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImageFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {

	CampaignDetailFormatter := CampaignDetailFormatter{}

	CampaignDetailFormatter.ID = campaign.ID
	CampaignDetailFormatter.Name = campaign.Name
	CampaignDetailFormatter.ShortDescription = campaign.ShortDescription
	CampaignDetailFormatter.Description = campaign.Description
	CampaignDetailFormatter.GoalAmount = campaign.GoalAmount
	CampaignDetailFormatter.CurrentAmount = campaign.CurrentAmount
	CampaignDetailFormatter.BackerCount = campaign.BackerCount
	CampaignDetailFormatter.UserID = campaign.UserID
	CampaignDetailFormatter.Slug = campaign.Slug
	CampaignDetailFormatter.ImageURL = ""
	if len(campaign.CampaignImages) > 0 {
		CampaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}
	CampaignDetailFormatter.Perks = perks

	// Init Campaign User..
	user := campaign.User
	campaignUserFormatter := CampaignUserFormatter{}
	campaignUserFormatter.Name = user.Name
	campaignUserFormatter.ImageURL = user.AvatarFileName

	CampaignDetailFormatter.User = campaignUserFormatter

	// Init Campaign Images
	images := []CampaignImageFormatter{}
	for _, image := range campaign.CampaignImages {

		campaignImageFormatter := CampaignImageFormatter{}

		campaignImageFormatter.ImageURL = image.FileName
		isPrimary := false
		if image.IsPrimary == 1 {
			isPrimary = true
		}
		campaignImageFormatter.IsPrimary = isPrimary
		images = append(images, campaignImageFormatter)
	}
	CampaignDetailFormatter.Images = images

	return CampaignDetailFormatter
}
