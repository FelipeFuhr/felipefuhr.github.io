package site

import (
	"html"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoDirs returns the real content directories relative to this package.
func repoDirs() (data, templates, assets, static string) {
	return "../../data", "../../templates", "../../assets", "../../static"
}

func TestLoad_RealData_ReturnsPopulatedContent(t *testing.T) {
	data, _, _, _ := repoDirs()
	c, err := Load(data)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if c.Site.Title == "" {
		t.Error("expected a site title")
	}
	if c.Site.Domain != "workshop.ffreis.com" {
		t.Errorf("domain = %q, want workshop.ffreis.com", c.Site.Domain)
	}
	if len(c.Projects) == 0 {
		t.Error("expected at least one project")
	}
	if len(c.Patterns) == 0 {
		t.Error("expected at least one pattern")
	}
}

func TestLoad_MissingDataDir_ReturnsError(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Fatal("expected error for missing data dir")
	}
}

func TestValidate_RejectsInvalidContent(t *testing.T) {
	cases := map[string]Content{
		"no title":   {Site: Site{Domain: "d", URL: "u"}, Projects: []Project{{Title: "x"}}},
		"no domain":  {Site: Site{Title: "t", URL: "u"}, Projects: []Project{{Title: "x"}}},
		"no url":     {Site: Site{Title: "t", Domain: "d"}, Projects: []Project{{Title: "x"}}},
		"no project": {Site: Site{Title: "t", Domain: "d", URL: "u"}},
		"empty title": {
			Site:     Site{Title: "t", Domain: "d", URL: "u"},
			Projects: []Project{{Title: ""}},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := c.validate(); err == nil {
				t.Errorf("validate(%s) = nil, want error", name)
			}
		})
	}
}

func TestProject_KindsAttr(t *testing.T) {
	p := Project{Kinds: []string{"ml", "infra"}}
	if got := p.KindsAttr(); got != "ml infra" {
		t.Errorf("KindsAttr() = %q, want %q", got, "ml infra")
	}
}

func TestBuild_WritesStaticSite(t *testing.T) {
	data, templates, assets, static := repoDirs()
	c, err := Load(data)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	out := t.TempDir()
	if err := Build(BuildOpts{
		Content:     c,
		TemplateDir: templates,
		AssetsDir:   assets,
		StaticDir:   static,
		OutDir:      out,
	}); err != nil {
		t.Fatalf("Build: %v", err)
	}

	index := readFile(t, filepath.Join(out, "index.html"))
	// Title contains an apostrophe; html/template escapes it (e.g. &#39;).
	for _, want := range []string{html.EscapeString(c.Site.Title), "id=\"projects\"", "data-kind", c.Projects[0].Title} {
		if !strings.Contains(index, want) {
			t.Errorf("index.html missing %q", want)
		}
	}

	if cname := strings.TrimSpace(readFile(t, filepath.Join(out, "CNAME"))); cname != c.Site.Domain {
		t.Errorf("CNAME = %q, want %q", cname, c.Site.Domain)
	}
	if sm := readFile(t, filepath.Join(out, "sitemap.xml")); !strings.Contains(sm, c.Site.URL) {
		t.Errorf("sitemap.xml missing url %q", c.Site.URL)
	}
	for _, asset := range []string{"assets/styles.css", "assets/app.js", "assets/favicon.svg", "robots.txt", "404.html"} {
		if _, err := os.Stat(filepath.Join(out, asset)); err != nil {
			t.Errorf("expected %s in output: %v", asset, err)
		}
	}
}

func TestBuild_EscapesContent(t *testing.T) {
	out := t.TempDir()
	_, templates, _, _ := repoDirs()
	c := Content{
		Site:     Site{Title: "t", Domain: "d", URL: "u"},
		Projects: []Project{{Title: "<script>alert(1)</script>"}},
	}
	if err := Build(BuildOpts{Content: c, TemplateDir: templates, AssetsDir: t.TempDir(), StaticDir: t.TempDir(), OutDir: out}); err != nil {
		t.Fatalf("Build: %v", err)
	}
	index := readFile(t, filepath.Join(out, "index.html"))
	if strings.Contains(index, "<script>alert(1)</script>") {
		t.Error("project title was not HTML-escaped")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}
