package rules

import (
	"rt2/frame"
)

func returnSeq(f frame.Frame) (frame.Sequence, frame.WAIT) {
	return frame.End()
}
