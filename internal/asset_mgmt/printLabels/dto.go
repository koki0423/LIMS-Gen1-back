package printLabels

// ===== Requests =====

type PrintConfig struct {
	UseHalfcut       bool `json:"use_halfcut"`
	ConfirmTapeWidth bool `json:"confirm_tape_width"`
	EnablePrintLog   bool `json:"enable_print_log"`
}

type LabelData struct {
	Checked bool   `json:"checked" binding:"required"`
	ColB    string `json:"col_b"   binding:"required"`
	ColC    string `json:"col_c"   binding:"required"`
	ColD    string `json:"col_d"   binding:"required"`
	ColE    string `json:"col_e"   binding:"required"`
}

type PrintRequest struct {
	Config PrintConfig `json:"config" binding:"required"`
	Label  LabelData   `json:"label"  binding:"required"`
	Width  int         `json:"width"  binding:"required"`
	Type   string      `json:"type"   binding:"required"`
}

type PrintResponse struct {
	Success bool       `json:"success"`
	Error   *APIError `json:"error,omitempty"`
}