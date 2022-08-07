package app

type mode int

func (m mode) IsCommandMode() bool {
	return m&modeCommand == modeCommand
}

func (m mode) IsExplorerMode() bool {
	return m&modeExplorer == modeExplorer
}

func (m mode) IsSongMode() bool {
	return m&modeSong == modeSong
}

const (
	modeSong            = 1
	modeSongCommand     = modeSong | modeCommand
	modeExplorer        = 2
	modeExplorerCommand = modeExplorer | modeCommand

	modeCommand = 4
)
