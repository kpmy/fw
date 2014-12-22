package rules

import (
	"rt2/frame"
)

func ifSeq(f frame.Frame) (frame.Sequence, frame.WAIT) {
	return frame.End()
}
