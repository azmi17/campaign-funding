package transaction

import user "go-campaign-funding/user"

type GetCampaignTransactionInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
