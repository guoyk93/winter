package wsnowid

import (
	"context"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
	"github.com/guoyk93/winter/wresty"
	"strconv"
)

var (
	ext = wext.Simple[options]("snowid")
)

// Next return a new id
func Next(ctx context.Context, altKeys ...string) string {
	return NextN(ctx, 1, altKeys...)[0]
}

// NextN return n ids generated from snowid service
func NextN(ctx context.Context, size int, altKeys ...string) []string {
	if size < 1 {
		winter.HaltString("wsnowid: invalid argument: size", winter.HaltWithBadRequest())
	}

	o := ext.Instance(altKeys...).Get(ctx)

	var ret []string
	res := rg.Must(
		wresty.R(ctx, o.restyKeys...).
			SetQueryParam("size", strconv.Itoa(size)).
			SetResult(&ret).
			Get(o.url),
	)

	if res.IsError() {
		winter.HaltString(res.String())
	}
	if len(ret) != size {
		winter.HaltString("wsnowid: invalid returns")
	}
	return ret
}

// Installer install component
func Installer(opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
