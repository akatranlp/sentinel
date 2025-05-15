package enums

// ENUM(page, popup, touch, wap)
type Display string

func (d *Display) WithDefault() {
	if *d == "" {
		*d = DisplayPage
	}
}
