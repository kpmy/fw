package rules

import (
	"math/big"
)

import (
	"fw/rt2"
	"fw/rt2/frame"
	"fw/rt2/scope"
)

func bit_range(_f scope.Value, _t scope.Value) scope.Value {
	f := scope.GoTypeFrom(_f).(int32)
	t := scope.GoTypeFrom(_t).(int32)
	ret := big.NewInt(0)
	for i := f; i <= t; i++ {
		ret = ret.SetBit(ret, int(i), 1)
	}
	return scope.TypeFromGo(ret)
}

func rangeSeq(f frame.Frame) (frame.Sequence, frame.WAIT) {
	n := rt2.NodeOf(f)
	return This(expectExpr(f, n.Left(), func(...IN) OUT {
		return expectExpr(f, n.Right(), func(...IN) OUT {
			rt2.ValueOf(f.Parent())[n.Adr()] = bit_range(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return End()
		})
	}))
}
