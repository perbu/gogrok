package analytics

func (p *Package) Lines() int {
	lines := 0
	for _, file := range p.files {
		lines += len(file.source)
	}
	return lines
}

func (p *Package) Files() int {
	return len(p.files)
}

func (p *Package) GetFiles() []*File {
	return p.files
}

func (p *Package) GetFile(name string) (*File, bool) {
	for _, file := range p.files {
		if file.Name == name {
			return file, true
		}
	}
	return nil, false
}
