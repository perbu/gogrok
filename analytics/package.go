package analytics

import "log/slog"

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

// Generated returns the ratio of generated files in the package, 1.0 means all files are generated, 0.0 means none are generated
func (p *Package) Generated() float32 {
	generated := 0
	total := 0
	for _, file := range p.files {
		if file.generated {
			generated += file.Lines()
		}
		total += file.Lines()
	}
	return float32(generated) / float32(total)
}

func (p *Package) CalculateComplexity() float32 {
	complexity := float32(0)
	for _, file := range p.files {
		complexity += file.CalculateComplexity()
	}
	slog.Info("complexity", "package", p.Name, "complexity", complexity, "files", len(p.files))
	return complexity / float32(len(p.files))
}
