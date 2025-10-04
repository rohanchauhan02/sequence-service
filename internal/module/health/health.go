package health

type Repository interface {
	Health() (map[string]any, error)
}

type Usecase interface {
	Health() (map[string]any, error)
}
