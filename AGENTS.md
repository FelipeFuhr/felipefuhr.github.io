# Agent Context

**This repo:** `felipefuhr.github.io` — the workshop site behind ffreis.com.
A single static page generated from hand-curated YAML by a small Go tool, served
at **workshop.ffreis.com** via GitHub Pages.

## Non-obvious facts

- **Content is YAML, not HTML.** Everything visible comes from [`data/`](data/)
  (`site.yaml`, `projects.yaml`, `buildlog.yaml`, `patterns.yaml`). Edit the YAML;
  do not hand-edit generated HTML. The generator re-renders `dist/` from scratch.
- **The generator is the build.** [`cmd/build`](cmd/build) → [`internal/site`](internal/site)
  loads + validates the YAML and renders [`templates/`](templates/) (`html/template`,
  auto-escaping) into `dist/`. There is no other build step.
- **`CNAME` is generated, not committed.** `internal/site.Build` writes
  `dist/CNAME` from `site.yaml`'s `domain` — the custom domain lives in **one**
  place (`data/site.yaml`). Do not add a separate root `CNAME` file.
- **Project `href` must be a PUBLIC URL.** A card linking a private repo 404s for
  visitors. Leave `href` empty to render a non-linked card.
- **Filter keys are a contract.** Each project's `kinds:` must use a key defined in
  `site.yaml`'s `filters:`; `assets/app.js` filters on `data-kind`.
- **Pages deploy is `main`-only.** [`.github/workflows/pages.yml`](.github/workflows/pages.yml)
  builds and publishes on push to `main`. PRs (especially drafts) never deploy.
- **Coverage floor is 75%** (`COVERAGE_MIN` in the Makefile), matching the fleet Go
  floor — below the go-cli template's stricter 90% because this is a tiny generator.

## Add a project card

1. Append an entry to [`data/projects.yaml`](data/projects.yaml) (see the header
   comment for the schema; `kinds` must be filter keys from `site.yaml`).
2. `make build` (or `make serve` to preview at localhost:8080).
3. Commit; `pages.yml` deploys on merge to `main`.

## Build and run

```bash
make build            # render dist/
make serve            # build + serve at http://localhost:8080
make ci               # fmt-check + lint + coverage-gate + build (the pre-PR gate)
make lefthook-bootstrap  # install git hooks
```

## Structure

```
cmd/build/         entry point (flags → site.Load → site.Build)
internal/site/     YAML model + loader + renderer (the tested core)
data/              hand-curated content (the source of truth)
templates/         html/template partials (base, nav, hero, projects, …)
assets/            styles.css, app.js, favicon.svg
static/            root files copied verbatim (robots.txt, 404.html)
.github/workflows/ ci.yml (Go lint+test) · pages.yml (deploy) · devops-* (governance)
```

## Keeping this file current

- New content sections → update the structure list and the data contract notes.
- New non-obvious constraint → add it under "Non-obvious facts".
