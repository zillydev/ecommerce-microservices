package product

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetAllProducts() ([]Product, error) {
	return s.Repo.GetAllProducts()
}

func (s *Service) GetProductByID(id int) (*Product, error) {
	return s.Repo.GetProductByID(id)
}

func (s *Service) AddProduct(product *NewProductRequest) (*Product, error) {
	return s.Repo.AddProduct(product)
}
