[![GoDoc](https://godoc.org/github.com/KarpelesLab/openuri?status.svg)](https://godoc.org/github.com/KarpelesLab/openuri)

# openuri

This is a library that enables opening files and urls transparently.

It supports:

* Direct file access (no protocol)
* `file://` protocol
* `http://` and `https://` trought `net/http`
* `data:` uri
* Custom protocols

## Usage

	f, err := openuri.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	...

