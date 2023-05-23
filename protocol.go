package openuri

import (
	"context"
	"io"
	"io/fs"
	"net/url"
)

type Handler interface {
	OpenURI(ctx context.Context, c *Client, u *url.URL) (io.ReadCloser, error)
	StatURI(ctx context.Context, c *Client, u *url.URL) (fs.FileInfo, error)
}

var Protocols = map[string]Handler{
	"file":  Local,
	"http":  Http,
	"https": Http,
	"data":  Data,
}
