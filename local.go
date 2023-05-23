package openuri

import (
	"context"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
)

// localProtocol simply implements opening local urls
//
// TODO: file:// protocol has actually quite a few quirks, see:
// https://en.wikipedia.org/wiki/File_URI_scheme
type localProtocol struct{}

var Local = &localProtocol{}

func (l *localProtocol) OpenURI(ctx context.Context, c *Client, u *url.URL) (io.ReadCloser, error) {
	if !c.AllowLocal {
		return nil, ErrProtocolNotSupported
	}

	fn, err := l.parseUrl(u)
	if err != nil {
		return nil, err
	}
	return os.Open(fn)
}

func (l *localProtocol) StatURI(ctx context.Context, c *Client, u *url.URL) (fs.FileInfo, error) {
	if !c.AllowLocal {
		return nil, ErrProtocolNotSupported
	}

	fn, err := l.parseUrl(u)
	if err != nil {
		return nil, err
	}
	return os.Stat(fn)
}

func (l *localProtocol) parseUrl(u *url.URL) (string, error) {
	if !u.IsAbs() {
		return "", ErrNotAbsolute
	}
	if u.Host != "" && u.Host != "localhost" {
		return "", ErrLocalInvalidHost
	}

	fn := filepath.FromSlash(u.Path)
	return fn, nil
}
