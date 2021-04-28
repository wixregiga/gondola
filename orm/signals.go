package orm

import (
	"gondola/signals"
)

type ormSignal struct {
	s *signals.Signal
}

func (s *ormSignal) Listen(handler func(o *Orm)) signals.Listener {
	return s.s.Listen(func(data interface{}) {
		handler(data.(*Orm))
	})
}

func (s *ormSignal) emit(o *Orm) {
	s.s.Emit(o)
}

// Signals declares the signals emitted by this package. See
// gondola/signals for more information.
var Signals = struct {
	// WillInitialize is emitted just before a gondola/orm.Orm is
	// initialized.
	WillInitialize *ormSignal
}{
	&ormSignal{signals.New("will-initialize")},
}
