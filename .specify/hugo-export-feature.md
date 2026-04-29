# Hugo CMS Export Feature

## Date: 2026-04-29

## Summary

Added `linkmgr hugo` command to export bookmarks as Hugo CMS-compatible markdown files, organized by tags.

---

## Overview

The Hugo export feature generates static markdown files with Hugo front matter, making it easy to publish your bookmarks as part of a Hugo website.

---

## Command

```bash
linkmgr hugo [flags]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `content/bookmarks` | Output directory for markdown files |
| `--section` | | `bookmarks` | Hugo section name |
| `--type` | | `bookmark` | Hugo content type |
| `--author` | | | Author name for front matter |
| `--draft` | | `false` | Mark content as draft |
| `--per-tag` | | `true` | Create one file per tag (default mode) |
| `--single-file` | | `false` | Create single file with all bookmarks |

---

## Export Modes

### 1. Per-Tag Mode (Default)

Creates one markdown file for each tag:

```bash
linkmgr hugo --output content/bookmarks
```

**Output:**
```
content/bookmarks/
├── golang.md
├── kubernetes.md
├── python.md
└── javascript.md
```

**Generated File** (`golang.md`):
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
```

### 2. Single File Mode

Creates one `_index.md` with all bookmarks:

```bash
linkmgr hugo --output content/bookmarks --single-file
```

**Output:**
```
content/bookmarks/
└── _index.md
```

Bookmarks are organized by tag within the single file.

---

## Front Matter Generated

Each markdown file includes Hugo front matter:

```yaml
---
title: "Golang Bookmarks"           # Generated from tag name
date: 2026-04-29T18:42:00-00:00    # Current timestamp
author: Your Name                   # From --author flag
type: bookmark                      # From --type flag
tags:                              # Array of tags
  - golang
draft: false                        # From --draft flag
description: "Collection of 15..."  # Auto-generated
---
```

---

## Content Format

For each bookmark:

```markdown
## [Bookmark Title](https://url.com)

Excerpt text if available.

**Author:** Author name if available

**Also tagged:** `tag1`, `tag2`  # Other tags besides the primary

### Summary

Content preview if available (truncated to 500 chars)...

---
```

---

## Hugo Integration

### Basic Setup

1. **Generate bookmarks:**
   ```bash
   linkmgr hugo --output ~/my-hugo-site/content/bookmarks
   ```

2. **Build Hugo site:**
   ```bash
   cd ~/my-hugo-site
   hugo
   ```

3. **Serve locally:**
   ```bash
   hugo server
   ```

### Directory Structure

```
my-hugo-site/
├── config.toml
├── content/
│   └── bookmarks/
│       ├── golang.md
│       ├── kubernetes.md
│       └── python.md
├── layouts/
│   └── bookmarks/
│       ├── list.html      # List template
│       └── single.html    # Single template
└── themes/
```

### Hugo Config Example

Add to `config.toml`:

```toml
[params]
  [params.bookmarks]
    enabled = true
    title = "My Bookmarks"

[[menu.main]]
  name = "Bookmarks"
  url = "/bookmarks/"
  weight = 5

[taxonomies]
  tag = "tags"
```

### Hugo Template Example

Create `layouts/bookmarks/list.html`:

```html
{{ define "main" }}
<article class="bookmarks">
  <header>
    <h1>{{ .Title }}</h1>
    <p>{{ .Params.description }}</p>
  </header>
  
  <div class="content">
    {{ .Content }}
  </div>
  
  <footer>
    <div class="tags">
      {{ range .Params.tags }}
        <a href="/tags/{{ . | urlize }}">{{ . }}</a>
      {{ end }}
    </div>
  </footer>
</article>
{{ end }}
```

---

## Use Cases

### 1. Personal Knowledge Base

```bash
# Export all bookmarks organized by topic
linkmgr hugo --output ~/blog/content/resources
cd ~/blog && hugo
```

### 2. Team Resource Library

```bash
# Export as drafts for review
linkmgr hugo --output content/team-resources --draft --author "Team"
```

### 3. Public Bookmark Collection

```bash
# Export only public bookmarks
# (Note: current version exports all; filter in future version)
linkmgr hugo --output content/public-links
```

### 4. Automated Updates

**Cron job** (daily update):
```bash
0 2 * * * cd /path/to/hugo-site && linkmgr hugo -o content/bookmarks && hugo
```

**Git hook** (pre-commit):
```bash
#!/bin/bash
linkmgr hugo --output content/bookmarks
git add content/bookmarks/
```

---

## Advanced Usage

### Custom Content Type

```bash
# Use 'link' content type instead of 'bookmark'
linkmgr hugo --output content/links --type link --section links
```

### Multiple Sections

```bash
# Export different tags to different sections
linkmgr hugo --output content/dev-resources
linkmgr hugo --output content/design-resources
```

### Integration with CI/CD

```yaml
# .github/workflows/build.yml
name: Build Hugo Site
on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Export bookmarks
        run: linkmgr hugo --output content/bookmarks
      - name: Build Hugo
        run: hugo
      - name: Deploy
        run: ./deploy.sh
```

---

## Benefits

### SEO & Performance
- ✅ Static HTML (fast loading)
- ✅ SEO-friendly URLs
- ✅ No database queries at runtime
- ✅ CDN-friendly

### Developer Experience
- ✅ Version control friendly
- ✅ Easy to preview (hugo server)
- ✅ Full Hugo theming support
- ✅ Can be automated

### Content Management
- ✅ Organized by tags
- ✅ Clean markdown format
- ✅ Easy to customize
- ✅ Supports Hugo shortcodes

### Maintenance
- ✅ One command export
- ✅ Automated updates
- ✅ Git-tracked changes
- ✅ Easy rollback

---

## Implementation Details

### File Naming

- Tags are converted to safe filenames
- Spaces replaced with hyphens
- Lowercase conversion
- Example: "Go Lang" → "go-lang.md"

### Content Organization

- Primary tag determines filename
- Other tags listed as "Also tagged"
- Bookmarks sorted by tag
- Hierarchical heading structure

### Truncation

- Content limited to 500 characters
- Maintains readability
- Full content in database
- Add "..." for truncated content

---

## Customization Tips

### 1. Custom Shortcodes

Create `layouts/shortcodes/bookmark.html`:
```html
<div class="bookmark-card">
  <h3><a href="{{ .Get "url" }}">{{ .Get "title" }}</a></h3>
  <p>{{ .Inner }}</p>
</div>
```

### 2. Styling

Add to your theme's CSS:
```css
.bookmarks h2 a {
  color: #0066cc;
  text-decoration: none;
}

.bookmarks h2 a:hover {
  text-decoration: underline;
}

.tags a {
  background: #f0f0f0;
  padding: 2px 8px;
  border-radius: 3px;
  margin-right: 5px;
}
```

### 3. List Page Template

Create `layouts/bookmarks/list.html` for tag overview:
```html
{{ define "main" }}
<h1>Bookmark Collections</h1>
<ul>
  {{ range .Pages }}
    <li>
      <a href="{{ .Permalink }}">{{ .Title }}</a>
      <span class="count">{{ len .Content }}</span>
    </li>
  {{ end }}
</ul>
{{ end }}
```

---

## Future Enhancements

Potential improvements:

- [ ] Filter export by tag (only export specific tags)
- [ ] Filter by public/private status
- [ ] Individual bookmark files (one file per bookmark)
- [ ] Custom templates for markdown generation
- [ ] Related bookmarks section
- [ ] Date-based organization
- [ ] Multi-language support

---

## Testing

### Manual Test

```bash
# 1. Export bookmarks
linkmgr hugo --output /tmp/test-hugo

# 2. Check output
ls -la /tmp/test-hugo/
cat /tmp/test-hugo/golang.md

# 3. Test with Hugo
cd /tmp
hugo new site test-site
cp -r test-hugo test-site/content/bookmarks
cd test-site
hugo server
```

### Validation

- ✅ Valid YAML front matter
- ✅ Valid markdown syntax
- ✅ Hugo builds without errors
- ✅ Links are functional
- ✅ Tags work as taxonomy

---

## Conclusion

✅ **Hugo CMS export feature complete**  
✅ **Per-tag and single-file modes**  
✅ **Full Hugo front matter support**  
✅ **Clean, readable markdown output**  
✅ **Easy Hugo integration**

This feature enables you to publish your curated bookmark collection as a static website using Hugo, with full theming support and SEO benefits!
