# felipefuhr.github.io

<!-- ffreis-badges:start -->
[![CI](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/FelipeFuhr/ffreis-badges/main/badges/felipefuhr.github.io/ci.json)](https://github.com/FelipeFuhr/felipefuhr.github.io/actions) [![License](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/FelipeFuhr/ffreis-badges/main/badges/felipefuhr.github.io/license.json)](https://github.com/FelipeFuhr/felipefuhr.github.io/blob/main/LICENSE)
<!-- ffreis-badges:end -->

The workshop behind [ffreis.com](https://ffreis.com/en/) — experiments, build
logs, repos, and reusable patterns. A single static page, served at
**[workshop.ffreis.com](https://workshop.ffreis.com)** via GitHub Pages.

## How it works

The site is **hand-curated YAML rendered to static HTML by a small Go generator** —
HTML + CSS first, with a sprinkle of JavaScript only for the project filter and
copy buttons. No framework, no build step beyond the generator.

- **Content** lives in [`data/`](data/) (`site.yaml`, `projects.yaml`,
  `buildlog.yaml`, `patterns.yaml`) — edit YAML, never the HTML.
- **Templates** in [`templates/`](templates/) (`html/template`), styles in
  [`assets/styles.css`](assets/styles.css), behaviour in
  [`assets/app.js`](assets/app.js).
- **Generator** [`cmd/build`](cmd/build) + [`internal/site`](internal/site) reads
  the YAML and writes `dist/` (index, assets, `CNAME`, `sitemap.xml`, `404.html`).
- **Deploy**: [`.github/workflows/pages.yml`](.github/workflows/pages.yml) builds
  and publishes `dist/` to GitHub Pages on every push to `main`.

To add a project card, append an entry to [`data/projects.yaml`](data/projects.yaml)
and run `make build`.

## Development

```bash
make build   # render the site into dist/
make serve    # build + serve at http://localhost:8080
make ci       # full local gate: fmt-check + lint + coverage-gate + build
```

This repo follows the fleet standards — git hooks (lefthook) and CI are consumed
from `ffreis-platform-standards` / `ffreis-workflows-*` by pinned reference.
`make lefthook-bootstrap` installs the hooks. See [`AGENTS.md`](AGENTS.md) for
this repo's conventions and the full set of targets.

## Badges

CI / version / license badges above are served from the public
[`ffreis-badges`](https://github.com/FelipeFuhr/ffreis-badges) mirror, so they
render even while this repo is private. They populate once this repo is in the
mirror's manifest (the poller refreshes on a schedule).

## License

See [`LICENSE`](LICENSE).
