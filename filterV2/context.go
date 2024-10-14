package filterV2

import (
	"context"

	"github.com/kyaxcorp/go-helper/_context"
)

func (f *Input) SetContext(ctx context.Context) *Input {
	if ctx == nil {
		ctx = _context.GetDefaultContext()
	}
	f.ctx = ctx
	return f
}

func (f *Input) checkContext() *Input {
	if f.ctx == nil {
		f.ctx = _context.GetDefaultContext()
	}
	return f
}
