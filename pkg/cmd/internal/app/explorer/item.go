package explorer

type item struct {
	Path string
	Text string
}

type items []item

func (is items) String(i int) string {
	return is[i].Text
}

func (is items) Len() int {
	return len(is)
}

type filteredItem struct {
	*item
	itemIdx int
}
