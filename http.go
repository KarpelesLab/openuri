package openuri

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/KarpelesLab/webutil"
)

type HttpProtocol struct {
	Client *http.Client
}

var Http = &HttpProtocol{Client: http.DefaultClient}

func (h *HttpProtocol) OpenURI(ctx context.Context, c *Client, u *url.URL) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("while creating HTTP request: %w", err)
	}
	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while performing HTTP request: %w", err)
	}
	// TODO we should process redirects here and stop go from handling these so we can support stuff like redirect to ftp, etc
	// Note: redirect from http to local protocols is dangerous and should never be allowed ever
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return nil, fs.ErrNotExist
		}

		return nil, webutil.HttpError(resp.StatusCode)
	}

	// everything is good, just return body
	return resp.Body, nil
}

func (h *HttpProtocol) StatURI(ctx context.Context, c *Client, u *url.URL) (fs.FileInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("while creating HTTP request: %w", err)
	}
	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while performing HTTP request: %w", err)
	}
	resp.Body.Close()

	return &httpFileInfo{resp: resp}, nil
}

// httpFileInfo is a basic container for stat() to work on http urls. Not all information will be returned
type httpFileInfo struct {
	resp *http.Response
}

func (i *httpFileInfo) Name() string {
	return path.Base(i.resp.Request.URL.Path)
}

func (i *httpFileInfo) Size() int64 {
	return i.resp.ContentLength
}

func (i *httpFileInfo) Mode() fs.FileMode {
	return fs.FileMode(0444)
}

func (i *httpFileInfo) ModTime() time.Time {
	h := i.resp.Header.Get("Last-Modified")
	if h != "" {
		v, err := http.ParseTime(h)
		if err == nil {
			return v
		}
	}
	return time.Time{}
}

func (i *httpFileInfo) IsDir() bool {
	return false
}

func (i *httpFileInfo) Sys() any {
	return i.resp
}
