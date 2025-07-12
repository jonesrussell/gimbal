package common

type HUDData struct {
	Score  int
	Lives  int
	Level  int
	Health float64
}

type MenuData struct {
	Title         string
	Options       []MenuOption
	SelectedIndex int
}

type MenuOption struct {
	Label string
	Value interface{}
}
