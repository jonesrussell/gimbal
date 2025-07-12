package ui

// HUDData contains data for HUD updates
type HUDData struct {
	Score  int
	Lives  int
	Level  int
	Health float64
}

// MenuData contains data for menu systems
type MenuData struct {
	Title         string
	Options       []MenuOption
	SelectedIndex int
}

// MenuOption represents a menu option
type MenuOption struct {
	Label string
	Value interface{}
}
