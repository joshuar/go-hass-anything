// Code generated by "stringer -type=TextEntityMode -output textEntityModeStrings.go -linecomment"; DO NOT EDIT.

package hass

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PlainText-0]
	_ = x[Password-1]
}

const _TextEntityMode_name = "textpassword"

var _TextEntityMode_index = [...]uint8{0, 4, 12}

func (i TextEntityMode) String() string {
	if i < 0 || i >= TextEntityMode(len(_TextEntityMode_index)-1) {
		return "TextEntityMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TextEntityMode_name[_TextEntityMode_index[i]:_TextEntityMode_index[i+1]]
}
