package order

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetAllOrdersByUser(userID int) ([]Order, error) {
	return s.Repo.GetAllOrdersByUser(userID)
}

func (s *Service) GetOrderByID(id int) (*Order, error) {
	return s.Repo.GetOrderByID(id)
}

func (s *Service) CreateOrder(order *NewOrderRequest, total float64) (*Order, error) {
	return s.Repo.CreateOrder(order, total)
}

func (s *Service) UpdateOrderStatus(id int, status string) (int, error) {
	return s.Repo.UpdateOrderStatus(id, status)
}
