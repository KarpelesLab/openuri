package openuri

import (
	"context"
	"io"
	"io/fs"
	"net/url"
	"os"
)

type Client struct {
	AllowLocal bool // if false, file:// and local access will be disabled. If true, even removing the file:// protocol will still allow direct file opening
}

var (
	DefaultClient = &Client{}
	AllowLocal    = &Client{AllowLocal: true}
)

func (c *Client) Open(fn string) (io.ReadCloser, error) {
	return c.OpenContext(context.Background(), fn)
}

func (c *Client) OpenContext(ctx context.Context, fn string) (io.ReadCloser, error) {
	u, err := url.Parse(fn)
	if err != nil || !u.IsAbs() {
		if c.AllowLocal {
			// attempt local open
			return os.Open(fn)
		}
		if err != nil {
			return nil, err
		}
	}

	return c.OpenURLContext(ctx, u)
}

func (c *Client) OpenURL(u *url.URL) (io.ReadCloser, error) {
	return c.OpenURLContext(context.Background(), u)
}

func (c *Client) OpenURLContext(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
	if !u.IsAbs() {
		return nil, ErrNotAbsolute
	}
	if !c.AllowLocal && u.Scheme == "file" {
		// forbid local access
		return nil, ErrProtocolNotSupported
	}
	proto, ok := Protocols[u.Scheme]
	if !ok {
		return nil, ErrProtocolNotSupported
	}
	if !c.AllowLocal && proto == Local {
		// we double check if AllowLocal=false in case file:// was renamed
		return nil, ErrProtocolNotSupported
	}

	return proto.OpenURI(ctx, c, u)
}

func (c *Client) Stat(fn string) (fs.FileInfo, error) {
	return c.StatContext(context.Background(), fn)
}

func (c *Client) StatContext(ctx context.Context, fn string) (fs.FileInfo, error) {
	u, err := url.Parse(fn)
	if err != nil || !u.IsAbs() {
		if c.AllowLocal {
			// attempt local open
			return os.Stat(fn)
		}
		if err != nil {
			return nil, err
		}
	}

	return c.StatURLContext(ctx, u)
}

func (c *Client) StatURL(u *url.URL) (fs.FileInfo, error) {
	return c.StatURLContext(context.Background(), u)
}

func (c *Client) StatURLContext(ctx context.Context, u *url.URL) (fs.FileInfo, error) {
	if !u.IsAbs() {
		return nil, ErrNotAbsolute
	}
	if !c.AllowLocal && u.Scheme == "file" {
		return nil, ErrProtocolNotSupported
	}
	proto, ok := Protocols[u.Scheme]
	if !ok {
		return nil, ErrProtocolNotSupported
	}
	if !c.AllowLocal && proto == Local {
		// we double check if AllowLocal=false in case file:// was renamed
		return nil, ErrProtocolNotSupported
	}

	return proto.StatURI(ctx, c, u)
}

// default handlers
func Open(fn string) (io.ReadCloser, error) {
	return DefaultClient.Open(fn)
}

func OpenContext(ctx context.Context, fn string) (io.ReadCloser, error) {
	return DefaultClient.OpenContext(ctx, fn)
}

func OpenURL(u *url.URL) (io.ReadCloser, error) {
	return DefaultClient.OpenURL(u)
}

func OpenURLContext(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
	return DefaultClient.OpenURLContext(ctx, u)
}

func Stat(fn string) (fs.FileInfo, error) {
	return DefaultClient.Stat(fn)
}

func StatContext(ctx context.Context, fn string) (fs.FileInfo, error) {
	return DefaultClient.StatContext(ctx, fn)
}

func StatURL(u *url.URL) (fs.FileInfo, error) {
	return DefaultClient.StatURL(u)
}

func StatURLContext(ctx context.Context, u *url.URL) (fs.FileInfo, error) {
	return DefaultClient.StatURLContext(ctx, u)
}
