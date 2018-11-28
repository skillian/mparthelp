// Package mparthelp is a mime/multipart helper package.
package mparthelp

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"

	"github.com/skillian/errors"
)

// Parts is a collection of parts of a multipart message.
type Parts []Part

// Into creates a multipart message into the given target from the provided
// parts.
func (p Parts) Into(target io.Writer) (formDataContentType string, err error) {
	w := multipart.NewWriter(target)
	for _, part := range p {
		err := part.Source.Add(part.Name, w)
		if err != nil {
			return "", errors.ErrorfWithCause(
				err,
				"failed to add %T part %v to %v: %v",
				part, part, w, err)
		}
	}
	formDataContentType = w.FormDataContentType()
	return formDataContentType, w.Close()
}

// Part defines a named part inside of a multipart message.
type Part struct {
	Name string
	Source
}

// Source is a data source that can add itself to a mime/multipart.Writer.
type Source interface {
	Add(name string, w *multipart.Writer) error
}

// JSON is a Source implementation that handles marshaling a value to JSON
type JSON struct {
	Value interface{}
}

// Add implements the Source interface.
func (j JSON) Add(name string, w *multipart.Writer) error {
	jsonBytes, err := json.Marshal(j.Value)
	if err != nil {
		return err
	}
	part, err := w.CreateFormField(name)
	if err != nil {
		return err
	}
	jsonBuffer := bytes.NewBuffer(jsonBytes)
	_, err = io.Copy(part, jsonBuffer)
	return err
}

// File is a Source implementation for files read from an io.Reader.
type File struct {
	// Name is the name of the file, not to be confused with the name of the
	// Part.
	Name string

	// Reader is the data source that the part is populated from.
	io.Reader

	// Closer is an optional io.Closer that is called after reading the Reader
	io.Closer
}

// Add implements the Source interface.
func (f File) Add(name string, w *multipart.Writer) error {
	part, err := w.CreateFormFile(name, f.Name)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, f.Reader)
	if err != nil {
		return err
	}
	if f.Closer != nil {
		return f.Closer.Close()
	}
	return nil
}
