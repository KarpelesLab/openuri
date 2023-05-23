package openuri

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"io/fs"
	"net/url"
	"strings"
	"time"
)

type dataProtocol struct{}

var Data = &dataProtocol{}

func (d *dataProtocol) OpenURI(ctx context.Context, c *Client, u *url.URL) (io.ReadCloser, error) {
	buf, _, err := d.parse(u)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(buf)), nil
}

func (d *dataProtocol) StatURI(ctx context.Context, c *Client, u *url.URL) (fs.FileInfo, error) {
	buf, typ, err := d.parse(u)
	if err != nil {
		return nil, err
	}

	return &dataUriInfo{buf: buf, typ: typ}, nil
}

func (d *dataProtocol) parse(in *url.URL) ([]byte, string, error) {
	// we expect u.Scheme = "data" and u.Opaque="..."
	u := in.Opaque
	for len(u) > 0 && u[0] == '/' {
		u = u[1:]
	}
	p := strings.IndexByte(u, ',')

	if p == -1 {
		return nil, "", errors.New("could not locate data: uri value")
	}

	opts := strings.Split(u[:p], ";") // first opt will be mime type, last will be base64 if base64
	dat := []byte(u[p+1:])            // could be base64 encoded, we'll see this later

	mime := opts[0]
	if mime == "" {
		mime = "application/octet-stream"
	}
	opts = opts[1:]

	if len(opts) > 0 && opts[len(opts)-1] == "base64" {
		// decode base64
		dat = bytes.TrimRight(dat, "=") // trim+raw base64 so we accept both types
		res := make([]byte, base64.RawStdEncoding.DecodedLen(len(dat)))
		n, err := base64.RawStdEncoding.Decode(res, dat)

		if err != nil {
			return nil, "", err
		}

		dat = res[:n]
	} else {
		tmp, err := url.QueryUnescape(string(dat))
		if err != nil {
			return nil, "", err
		}
		dat = []byte(tmp)
	}

	return dat, mime, nil
}

type dataUriInfo struct {
	buf []byte
	typ string
}

func (d *dataUriInfo) Name() string {
	return ""
}

func (d *dataUriInfo) Size() int64 {
	return int64(len(d.buf))
}

func (d *dataUriInfo) Mode() fs.FileMode {
	return 0444
}

func (d *dataUriInfo) ModTime() time.Time {
	return time.Time{}
}

func (d *dataUriInfo) IsDir() bool {
	return false
}

func (d *dataUriInfo) Sys() any {
	return d.buf
}

func (d *dataUriInfo) MimeType() string {
	return d.typ
}
