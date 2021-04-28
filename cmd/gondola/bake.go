package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rainycape/command"

	"gondola/internal/gen/genutil"
	"gondola/log"
	"gondola/util/vfsutil"

	"github.com/rainycape/vfs"
)

type bakeOptions struct {
	Dir        string `help:"Root directory with the files to bake"`
	Name       string `help:"Variable name of the baked files. If empty, defaults to <dir>Data for data or <dir>FS for VFS."`
	VFS        bool   `help:"Wheter to generate a vfs.VFS variable or just an string which can be passed to VFS related functions" name:"vfs"`
	Out        string `name:"o" help:"Output filename. If empty, defaults to <dir>_baked.go"`
	Extensions string `name:"ext" help:"Additional extensions (besides html, css and js) to include, separated by commas"`
}

func bakeCommand(_ *command.Args, opts *bakeOptions) error {
	extensions := []string{".html", ".css", ".js"}
	if opts.Dir == "" {
		return errors.New("dir can't be empty")
	}
	if opts.Name == "" {
		base := filepath.Base(opts.Dir)
		if opts.VFS {
			opts.Name = base + "FS"
		} else {
			opts.Name = base + "Data"
		}
	}
	if opts.Out == "" {
		opts.Out = filepath.Base(opts.Dir) + "_baked.go"
	}
	// go ignores files starting with _
	opts.Out = strings.TrimLeft(opts.Out, "_")
	extensions = append(extensions, strings.Split(opts.Extensions, ",")...)
	var buf bytes.Buffer
	odir := filepath.Dir(opts.Out)
	p, err := build.ImportDir(odir, 0)
	if err == nil {
		buf.WriteString(fmt.Sprintf("package %s\n", p.Name))
	}
	buf.WriteString(genutil.AutogenString())
	if err := writeBakedFSCode(&buf, opts, extensions); err != nil {
		return err
	}
	if err := genutil.WriteAutogen(opts.Out, buf.Bytes()); err != nil {
		return err
	}
	log.Debugf("Assets written to %s (%d bytes)", opts.Out, buf.Len())
	return nil
}

func writeBakedFSCode(w io.Writer, opts *bakeOptions, extensions []string) (err error) {
	if opts.VFS {
		if _, err = io.WriteString(w, "import \"gondola/util/vfsutil\"\n"); err != nil {
			return err
		}
	}
	var buf bytes.Buffer
	if err := vfsutil.Bake(&buf, opts.Dir, extensions); err != nil {
		return err
	}
	data := buf.String()
	var files []string
	fs, err := vfs.TarGzip(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	err = vfs.Walk(fs, "/", func(fs vfs.VFS, path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, strings.TrimPrefix(path, "/"))
		}
		return err
	})
	if err != nil {
		return err
	}
	if len(files) > 0 {
		if _, err := fmt.Fprintf(w, "// %s \n", strings.Join(files, ", ")); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, "var %s = ", opts.Name); err != nil {
		return err
	}
	if opts.VFS {
		_, err = fmt.Fprintf(w, "vfsutil.MustOpenBaked(%q)\n", data)
	} else {
		_, err = fmt.Fprintf(w, "%q\n", data)
	}
	return err
}
