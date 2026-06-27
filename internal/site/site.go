// Package site loads the workshop's hand-curated YAML content and renders it
// into a single static HTML page plus its supporting root files.
package site

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Content is the full data model rendered into the page.
type Content struct {
	Site     Site
	Projects []Project
	BuildLog []LogEntry
	Patterns []Pattern
}

// Site holds global identity, navigation, hero copy, and filter definitions.
type Site struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Domain      string   `yaml:"domain"`
	URL         string   `yaml:"url"`
	Author      string   `yaml:"author"`
	OGImage     string   `yaml:"og_image"`
	Nav         []Link   `yaml:"nav"`
	Hero        Hero     `yaml:"hero"`
	Filters     []Filter `yaml:"filters"`
}

// Link is a labelled hyperlink; Primary marks the highlighted call-to-action.
type Link struct {
	Label   string `yaml:"label"`
	Href    string `yaml:"href"`
	Primary bool   `yaml:"primary"`
}

// Hero is the masthead: eyebrow, heading, lead, actions, and terminal lines.
type Hero struct {
	Eyebrow  string     `yaml:"eyebrow"`
	Heading  string     `yaml:"heading"`
	Lead     string     `yaml:"lead"`
	Actions  []Link     `yaml:"actions"`
	Terminal []TermLine `yaml:"terminal"`
}

// TermLine is one line of the decorative terminal. Kind selects a colour class
// (prompt, out, comment, accent).
type TermLine struct {
	Kind string `yaml:"kind"`
	Text string `yaml:"text"`
}

// Filter is a project-grid filter chip.
type Filter struct {
	Key   string `yaml:"key"`
	Label string `yaml:"label"`
}

// Project is one card in the experiments grid.
type Project struct {
	Group  string   `yaml:"group"`
	Status string   `yaml:"status"`
	Title  string   `yaml:"title"`
	Body   string   `yaml:"body"`
	Kinds  []string `yaml:"kinds"`
	Tags   []string `yaml:"tags"`
	Href   string   `yaml:"href"`
}

// KindsAttr joins a project's kinds for the data-kind filter attribute.
func (p Project) KindsAttr() string { return strings.Join(p.Kinds, " ") }

// LogEntry is one dated build-log item.
type LogEntry struct {
	Date  string `yaml:"date"`
	Title string `yaml:"title"`
	Body  string `yaml:"body"`
}

// Pattern is a copyable snippet block.
type Pattern struct {
	Label string `yaml:"label"`
	Code  string `yaml:"code"`
}

// Load reads and validates every YAML file in dataDir into a Content model.
func Load(dataDir string) (Content, error) {
	var c Content
	if err := readYAML(filepath.Join(dataDir, "site.yaml"), &c.Site); err != nil {
		return c, err
	}

	projects := struct {
		Projects []Project `yaml:"projects"`
	}{}
	if err := readYAML(filepath.Join(dataDir, "projects.yaml"), &projects); err != nil {
		return c, err
	}
	c.Projects = projects.Projects

	buildlog := struct {
		BuildLog []LogEntry `yaml:"buildlog"`
	}{}
	if err := readYAML(filepath.Join(dataDir, "buildlog.yaml"), &buildlog); err != nil {
		return c, err
	}
	c.BuildLog = buildlog.BuildLog

	patterns := struct {
		Patterns []Pattern `yaml:"patterns"`
	}{}
	if err := readYAML(filepath.Join(dataDir, "patterns.yaml"), &patterns); err != nil {
		return c, err
	}
	c.Patterns = patterns.Patterns

	if err := c.validate(); err != nil {
		return c, err
	}
	return c, nil
}

func readYAML(path string, dst any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func (c Content) validate() error {
	if c.Site.Title == "" {
		return fmt.Errorf("site.title is required")
	}
	if c.Site.Domain == "" {
		return fmt.Errorf("site.domain is required")
	}
	if c.Site.URL == "" {
		return fmt.Errorf("site.url is required")
	}
	if len(c.Projects) == 0 {
		return fmt.Errorf("at least one project is required")
	}
	for i, p := range c.Projects {
		if p.Title == "" {
			return fmt.Errorf("projects[%d]: title is required", i)
		}
	}
	return nil
}
