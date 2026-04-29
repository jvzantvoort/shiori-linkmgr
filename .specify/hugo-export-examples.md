# Hugo Export - Example Output

## Per-Tag Mode (Default)

When you run:
```bash
linkmgr hugo --output content/bookmarks
```

You get one markdown file per tag. For example, `golang.md`:

```markdown
---
title: "Golang Bookmarks"
date: 2026-04-29T18:42:00-00:00
author: Your Name
type: bookmark
tags:
  - golang
description: "Collection of 15 bookmarks tagged with golang"
---

# Golang Bookmarks

A curated collection of 15 links related to golang.

## [The Go Programming Language](https://golang.org)

Official Go programming language website with documentation and downloads.

**Author:** Google

**Also tagged:** `programming`, `documentation`

---

## [Effective Go](https://golang.org/doc/effective_go)

Tips for writing clear, idiomatic Go code.

**Also tagged:** `programming`, `tutorial`

---

## [Go by Example](https://gobyexample.com)

Hands-on introduction to Go using annotated example programs.

**Author:** Mark McGranaghan

**Also tagged:** `tutorial`, `learning`

### Summary

Go by Example is a hands-on introduction to Go using annotated example programs. Check out the first example or browse the full list below...

---
```

## Single File Mode

When you run:
```bash
linkmgr hugo --output content/bookmarks --single-file
```

You get a single `_index.md` file:

```markdown
---
title: "Bookmarks"
date: 2026-04-29T18:42:00-00:00
author: Your Name
type: bookmark
description: "Collection of 150 bookmarks"
---

# Bookmarks

Total bookmarks: 150

## Golang

### [The Go Programming Language](https://golang.org)

Official Go programming language website.

**Tags:** `golang`, `programming`, `documentation`

---

### [Effective Go](https://golang.org/doc/effective_go)

Tips for writing clear, idiomatic Go code.

**Tags:** `golang`, `programming`, `tutorial`

---

## Kubernetes

### [Kubernetes Documentation](https://kubernetes.io/docs/)

Official Kubernetes documentation.

**Tags:** `kubernetes`, `devops`, `cloud`

---
```

## Hugo Integration

### Directory Structure

```
my-hugo-site/
├── config.toml
├── content/
│   └── bookmarks/
│       ├── _index.md          # (if using --single-file)
│       ├── golang.md          # (if using --per-tag)
│       ├── kubernetes.md
│       ├── python.md
│       └── javascript.md
├── layouts/
│   └── bookmarks/
│       ├── list.html
│       └── single.html
└── themes/
```

### Example Hugo Config

Add to your `config.toml`:

```toml
[params]
  [params.bookmarks]
    enabled = true
    title = "My Bookmarks"
    description = "Curated collection of useful links"

[[menu.main]]
  name = "Bookmarks"
  url = "/bookmarks/"
  weight = 5
```

### Example Hugo Template

Create `layouts/bookmarks/list.html`:

```html
{{ define "main" }}
<article>
  <header>
    <h1>{{ .Title }}</h1>
    {{ if .Params.description }}
      <p class="description">{{ .Params.description }}</p>
    {{ end }}
  </header>
  
  <div class="content">
    {{ .Content }}
  </div>
  
  {{ if .Params.tags }}
  <footer>
    <div class="tags">
      {{ range .Params.tags }}
        <a href="/tags/{{ . | urlize }}" class="tag">{{ . }}</a>
      {{ end }}
    </div>
  </footer>
  {{ end }}
</article>
{{ end }}
```

## Usage Scenarios

### 1. Regular Export (Weekly Update)

```bash
# Update your Hugo site with latest bookmarks
linkmgr hugo --output ~/my-hugo-site/content/bookmarks
cd ~/my-hugo-site
hugo
```

### 2. Export Specific Tags

```bash
# Only export bookmarks with specific tag
linkmgr list --tag golang > /tmp/golang.txt
linkmgr hugo --output content/bookmarks --per-tag
```

### 3. Create Draft Pages

```bash
# Export as drafts for review
linkmgr hugo --output content/bookmarks --draft
```

### 4. Custom Content Type

```bash
# Use different content type
linkmgr hugo --output content/links --type links --section links
```

## Automation

### Cron Job

```bash
# Update bookmarks daily at midnight
0 0 * * * cd ~/my-hugo-site && /usr/local/bin/linkmgr hugo --output content/bookmarks && hugo
```

### Git Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Export bookmarks before committing
linkmgr hugo --output content/bookmarks

# Add generated files
git add content/bookmarks/*.md
```

## Customization

### Front Matter Options

The command supports these Hugo front matter fields:

- `title` - Generated from tag name
- `date` - Current timestamp
- `author` - Via `--author` flag
- `type` - Via `--type` flag (default: "bookmark")
- `tags` - Array of tags
- `draft` - Via `--draft` flag
- `description` - Auto-generated

### Content Customization

Edit `internal/cmd/hugo.go` to customize:

- Front matter fields
- Markdown formatting
- Section headings
- Link formatting
- Tag display

## Tips

1. **Use taxonomies**: Configure Hugo to use tags as a taxonomy
2. **Create list pages**: Use `_index.md` for overview pages
3. **Add shortcodes**: Create Hugo shortcodes for bookmark cards
4. **Style with CSS**: Add custom CSS for bookmark styling
5. **Enable search**: Use Hugo's search functionality for bookmarks

## Benefits

- ✅ SEO-friendly static pages
- ✅ Fast loading (static HTML)
- ✅ Version control friendly
- ✅ Full Hugo theming support
- ✅ Easy to customize
- ✅ Can be automated
- ✅ Works with existing Hugo sites
