package client

// Color
type Color int

const (
	Red Color = iota
	Green
	Blue
	Yellow
	Purple
	Cyan
	White
)

func (c Color) String() string {
	var color string
	switch c {
	case Red:
		color = "#FF0000"
	case Green:
		color = "#00FF00"
	case Blue:
		color = "#0000FF"
	case Yellow:
		color = "#FFFF00"
	case Purple:
		color = "#800080"
	case Cyan:
		color = "#48D1CC" // mediumturquoise
	case White:
		color = "#FFFFFF"
	}
	return color
}
