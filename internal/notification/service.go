package notification

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) CreateNotification(notification *NewNotificationRequest) (*Notification, error) {
	return s.Repo.CreateNotification(notification)
}

func (s *Service) MarkNotificationRead(notificationID int, userID int) error {
	return s.Repo.MarkNotificationRead(notificationID, userID)
}

func (s *Service) GetUnreadNotifications(userID int) ([]Notification, error) {
	return s.Repo.GetUnreadNotifications(userID)
}
