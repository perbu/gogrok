package repo

func New() (*Repo, error) {
	r := &Repo{
		modules: make(map[string]*Module),
	}
	return r, nil
}
