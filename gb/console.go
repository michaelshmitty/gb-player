package gb

type Console struct {
	Cartridge *Cartridge
}

func NewConsole(path string) (*Console, error) {
	cartridge := InsertCartridge(path)

	console := Console{cartridge}

	return &console, nil
}

func (console *Console) Update(dt uint64) {
	console.StepSeconds(dt)
}
