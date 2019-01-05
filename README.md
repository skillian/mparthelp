# mparthelp

## What is it?

mparthelp is a (hopefully?) rather boring wrapper around Go's mime/multipart
package for more easily generating multipart messages in Go.

## Example(s)

There's only one situation where I've used this that I can think of, within my
[ShareBase integration](https://github.com/skillian/sharebase/blob/bfc6f6443a17a802e4edae3b031b82186f87bcee/web/models.go#L248)
package (permalink above).

The idea is you create a slice of `mparthelp.Parts`.  Parts are structs but
their sources all implement the `mparthelp.Source` interface that adds the
data from the source into the actual mime/multipart.Writer.

Some example `mparthelp.Source` implementations are:

  - `mparthelp.JSON`:  When a part should be raw JSON.  The given value to the
    `mparthelp.JSON` is automatically marshaled to JSON with the
    `encoding/json` package.

  - `mparthelp.File`:  When a part should be read from an `io.Reader`.  Readers
    aren't closed automatically, though.  The current implementation requires
    a separate `io.Closer` be passed in the `mparthelp.File` struct to close
    the `io.Reader` after reading.  In retrospect, this should have implemented
    `io.ReadCloser` instead and recommended using `ioutil.NopCloser` if the
    source should not be closed after reading.
