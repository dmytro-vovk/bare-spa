package alias

type Alias struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
	Path  string `json:"path"`
}
