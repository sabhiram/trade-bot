package hub

type Hub struct {
}

func New() (*Hub, error) {
	return &Hub{}, nil
}
