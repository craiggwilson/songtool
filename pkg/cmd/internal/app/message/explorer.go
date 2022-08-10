package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/songio"
)

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
	Path string
	Meta *songio.Meta
}
