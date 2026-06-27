package site

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// BuildOpts groups the inputs and output location for a site build.
type BuildOpts struct {
	Content     Content
	TemplateDir string
	AssetsDir   string
	StaticDir   string
	OutDir      string
}

// Build renders the page and writes the full static site into OutDir: index.html,
// the assets/ tree, the root static files, a generated CNAME, and sitemap.xml.
func Build(o BuildOpts) error {
	if err := os.RemoveAll(o.OutDir); err != nil {
		return fmt.Errorf("clean %s: %w", o.OutDir, err)
	}
	if err := os.MkdirAll(o.OutDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", o.OutDir, err)
	}

	if err := renderToFile(filepath.Join(o.OutDir, "index.html"), o.TemplateDir, o.Content); err != nil {
		return err
	}
	if err := copyDir(o.AssetsDir, filepath.Join(o.OutDir, "assets")); err != nil {
		return fmt.Errorf("copy assets: %w", err)
	}
	if err := copyDir(o.StaticDir, o.OutDir); err != nil {
		return fmt.Errorf("copy static: %w", err)
	}

	// CNAME and sitemap derive from site.yaml, keeping the domain in one place.
	cname := o.Content.Site.Domain + "\n"
	if err := os.WriteFile(filepath.Join(o.OutDir, "CNAME"), []byte(cname), 0o644); err != nil {
		return fmt.Errorf("write CNAME: %w", err)
	}
	if err := writeSitemap(filepath.Join(o.OutDir, "sitemap.xml"), o.Content.Site.URL); err != nil {
		return fmt.Errorf("write sitemap: %w", err)
	}
	return nil
}

// Render writes the rendered HTML page to w using the templates in templateDir.
func Render(w io.Writer, templateDir string, c Content) error {
	tmpl, err := template.New("page").ParseGlob(filepath.Join(templateDir, "*.tmpl"))
	if err != nil {
		return fmt.Errorf("parse templates: %w", err)
	}
	if err := tmpl.ExecuteTemplate(w, "base", c); err != nil {
		return fmt.Errorf("render page: %w", err)
	}
	return nil
}

func renderToFile(path, templateDir string, c Content) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()
	return Render(f, templateDir, c)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func writeSitemap(path, url string) error {
	const tmpl = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>%s</loc></url>
</urlset>
`
	return os.WriteFile(path, []byte(fmt.Sprintf(tmpl, url)), 0o644)
}
