// Command build renders the workshop's YAML content into a static site under dist/.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/FelipeFuhr/FelipeFuhr.github.io/internal/site"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "build:", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	dataDir := fs.String("data", "data", "directory of YAML content")
	tmplDir := fs.String("templates", "templates", "directory of HTML templates")
	assetsDir := fs.String("assets", "assets", "directory of static assets (css/js)")
	staticDir := fs.String("static", "static", "directory of root static files (favicon, robots, 404)")
	outDir := fs.String("out", "dist", "output directory")
	if err := fs.Parse(args); err != nil {
		return err
	}

	content, err := site.Load(*dataDir)
	if err != nil {
		return err
	}
	if err := site.Build(site.BuildOpts{
		Content:     content,
		TemplateDir: *tmplDir,
		AssetsDir:   *assetsDir,
		StaticDir:   *staticDir,
		OutDir:      *outDir,
	}); err != nil {
		return err
	}
	fmt.Fprintf(out, "built %d projects, %d log entries, %d patterns -> %s\n",
		len(content.Projects), len(content.BuildLog), len(content.Patterns), *outDir)
	return nil
}
