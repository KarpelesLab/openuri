# openuri

This is a library that enables opening files and urls transparently.

It supports:

* Direct file access (no protocol)
* `file://` protocol
* `http://` and `https://` trought `net/http`
* Custom protocols

## Usage

	f, err := openuri.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	...

