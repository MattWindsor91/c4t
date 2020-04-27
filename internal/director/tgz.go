// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

/// TGZWriter is a wrapper over a gzipped tar writer.
type TGZWriter struct {
	tw *tar.Writer
	gw *gzip.Writer
	bw io.WriteCloser
}

// NewTGZWriter creates a new tarball writer on top of the basic writer w.
func NewTGZWriter(w io.WriteCloser) *TGZWriter {
	gw := gzip.NewWriter(w)
	tw := tar.NewWriter(gw)
	return &TGZWriter{
		tw: tw,
		gw: gw,
		bw: w,
	}
}

// Close closes all of this writer's inner writers in order, returning the first error.
func (t *TGZWriter) Close() error {
	terr := t.tw.Close()
	gerr := t.gw.Close()
	berr := t.bw.Close()

	if terr != nil {
		return fmt.Errorf("closing tar: %w", terr)
	}
	if gerr != nil {
		return fmt.Errorf("closing gzip: %w", gerr)
	}
	if berr != nil {
		return fmt.Errorf("closing file: %w", gerr)
	}
	return nil
}

// TarFile tars the file at rpath to wpath within the tar archive represented by tarw.
// If rpath is empty, no tarring occurs.
// If rpath doesn't exist, an error occurs unless NotFoundCb is set and handles the error in a different way.
func (t *TGZWriter) TarFile(rpath, wpath string) error {
	if rpath == "" {
		return nil
	}

	if err := t.tarFileHeader(rpath, wpath); err != nil {
		return err
	}
	return t.tarFileContents(rpath)
}

func (t *TGZWriter) tarFileHeader(rpath string, wpath string) error {
	hdr, err := makeTarFileHeader(rpath, wpath)
	if err != nil {
		return fmt.Errorf("making tar header for %s: %w", rpath, err)
	}
	if err := t.tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("writing tar header for %s: %w", rpath, err)
	}
	return nil
}

func (t *TGZWriter) tarFileContents(rpath string) error {
	f, err := os.Open(rpath)
	if err != nil {
		return fmt.Errorf("opening %s: %w", rpath, err)
	}
	if _, err = iohelp.CopyCloseSrc(t.tw, f); err != nil {
		return fmt.Errorf("archiving %s: %w", rpath, err)
	}
	return nil
}

func makeTarFileHeader(rpath, wpath string) (*tar.Header, error) {
	info, err := os.Stat(rpath)
	if err != nil {
		return nil, fmt.Errorf("can't stat %s: %w", rpath, err)
	}
	hdr, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return nil, fmt.Errorf("can't get header for %s: %w", rpath, err)
	}
	hdr.Name = wpath
	return hdr, nil
}