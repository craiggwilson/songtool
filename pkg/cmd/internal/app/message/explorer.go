package message

import tea "github.com/charmbracelet/bubbletea"

func LoadDirectory(path string) tea.Cmd {
	return func() tea.Msg {
		return LoadDirectoryMsg{Path: path}
	}
}

type LoadDirectoryMsg struct {
	Path string
}

func UpdateFiles(files []FileItem) tea.Cmd {
	return func() tea.Msg {
		return UpdateFilesMsg{Files: files}
	}
}

type UpdateFilesMsg struct {
	Files []FileItem
}

type FileItem struct {
	Path  string
	Title string
}
