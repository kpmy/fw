package rules

import (
	"math/big"
)

import (
	"fw/rt2"
	"fw/rt2/frame"
)

func bit_range(_f interface{}, _t interface{}) interface{} {
	f := int32Of(_f)
	t := int32Of(_t)
	ret := big.NewInt(0)
	for i := f; i <= t; i++ {
		ret = ret.SetBit(ret, int(i), 1)
	}
	return ret
}

func rangeSeq(f frame.Frame) (frame.Sequence, frame.WAIT) {
	n := rt2.NodeOf(f)
	return expectExpr(f, n.Left(), func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		return expectExpr(f, n.Right(), func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			rt2.DataOf(f.Parent())[n] = bit_range(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		})
	})
}
