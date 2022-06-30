package campaign

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {

	// If client are not send params by userID, so get all datas of campaigns..
	if userID == 0 {
		campaigns, err := s.repository.FindAll()
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	// If client send params by userID
	campaigns, err := s.repository.FindByUserID(userID)
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}
