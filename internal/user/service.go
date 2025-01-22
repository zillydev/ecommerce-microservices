package user

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) RegisterUser(user *NewUserRequest) (*User, error) {
	return s.Repo.CreateUser(user)
}

func (s *Service) UpdateUserPreferences(updatePreferencesRequest UpdatePreferencesRequest) error {
	return s.Repo.UpdatePreferences(updatePreferencesRequest)
}

func (s *Service) GetUserByID(userID int) (*User, error) {
	return s.Repo.GetUser(userID)
}

func (s *Service) GetUsersByPreferences(preferences GetUsersByPreferencesRequest) ([]User, error) {
	return s.Repo.GetUsersByPreferences(preferences)
}
