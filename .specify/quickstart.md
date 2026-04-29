# Quick Start Guide - linkmgr

## Installation

```bash
cd /home/jvzantvoort/Website/localsite/linkmgr
make build
sudo make install
```

---

## First Time Setup

### 1. Initialize Configuration
```bash
linkmgr init
```

Prompts for:
- Database host (localhost)
- Port (3306)
- Username
- Password
- Database name (linkmgr)

Creates `~/.linkmgr.yaml` with permissions 0600

---

## Basic Usage

### Add Bookmarks
```bash
# Simple
linkmgr add https://golang.org

# With metadata
linkmgr add https://golang.org \
  --title "Go Programming Language" \
  --excerpt "Official Go website" \
  --tags "programming,golang,documentation"

# With author
linkmgr add https://kubernetes.io \
  --title "Kubernetes" \
  --author "CNCF" \
  --tags "devops,k8s,cloud"
```

### List & Search
```bash
# List all
linkmgr list

# Filter by tag
linkmgr list --tag golang

# Search
linkmgr search "kubernetes"

# Limit results
linkmgr list --limit 10
```

### View Details
```bash
# Show bookmark
linkmgr show 5

# Show and open in browser
linkmgr show 5 --open
```

### Update Bookmarks
```bash
# Update URL (fix broken link)
linkmgr update 5 --url "https://new-url.com"

# Update title
linkmgr update 5 --title "New Title"

# Add tags
linkmgr update 5 --add-tags "tutorial,guide"

# Replace all tags
linkmgr update 5 --tags "new,tags,only"

# Remove tags
linkmgr update 5 --remove-tags "old-tag"
```

### Delete Bookmarks
```bash
# Delete with confirmation
linkmgr delete 5

# Delete multiple
linkmgr delete 5 10 15

# Skip confirmation
linkmgr delete 5 --force
```

---

## Tag Management

### List Tags
```bash
# Show all tags with counts
linkmgr tags

# Clean up orphaned tags (preview)
linkmgr tags --cleanup --dry-run

# Clean up orphaned tags
linkmgr tags --cleanup
```

### Rename Tags
```bash
# Rename globally across all bookmarks
linkmgr tag rename "golang" "go"
linkmgr tag rename "js" "javascript"
```

---

## Link Validation

### Check Link Health
```bash
# Validate all bookmarks
linkmgr validate

# Validate specific tag
linkmgr validate --tag important

# Fast validation (10 concurrent)
linkmgr validate --concurrency 10 --timeout 3

# Verbose (show all results)
linkmgr validate -v
```

### Fix Broken Links
```bash
# 1. Find broken links
linkmgr validate

# 2. View details
linkmgr show 5

# 3. Update URL
linkmgr update 5 --url "https://working-url.com"
```

---

## Export & Import

### Backup
```bash
# Export to JSON
linkmgr export bookmarks.json

# Export to CSV
linkmgr export bookmarks.csv

# Dated backup
linkmgr export backup-$(date +%Y%m%d).json
```

### Restore
```bash
# Import from file
linkmgr import bookmarks.json

# Skip duplicates
linkmgr import backup.json --skip-duplicates

# Update duplicates
linkmgr import backup.json --update-duplicates
```

---

## Hugo CMS Export

### Export for Hugo Site
```bash
# Per-tag mode (one file per tag)
linkmgr hugo --output content/bookmarks

# Single file mode
linkmgr hugo --output content/bookmarks --single-file

# With author
linkmgr hugo --output content/bookmarks --author "Your Name"

# As drafts
linkmgr hugo --output content/bookmarks --draft

# Custom content type
linkmgr hugo --output content/links --type link
```

### Hugo Workflow
```bash
# 1. Export bookmarks
linkmgr hugo --output ~/blog/content/bookmarks

# 2. Build Hugo site
cd ~/blog
hugo

# 3. Preview
hugo server

# 4. Deploy
./deploy.sh
```

---

## Complete Workflow Example

```bash
# 1. Add bookmarks during the week
linkmgr add https://example.com --title "Example" --tags "tutorial"
linkmgr add https://another.com --title "Another" --tags "tutorial,golang"

# 2. Validate links weekly
linkmgr validate

# 3. Fix broken links
linkmgr update 5 --url "https://new-url.com"

# 4. Clean up tags monthly
linkmgr tags
linkmgr tag rename "golang" "go"
linkmgr tags --cleanup

# 5. Backup weekly
linkmgr export backup-$(date +%Y%m%d).json

# 6. Publish to Hugo (if using)
linkmgr hugo --output ~/blog/content/bookmarks
cd ~/blog && hugo && ./deploy.sh
```

---

## Tips & Tricks

### Tag Naming
- Use lowercase
- Use hyphens for multi-word tags
- Be consistent: "golang" not "go-lang" or "GoLang"

### Organization
- Tag by topic: "programming", "devops", "tutorial"
- Tag by language: "golang", "python", "javascript"
- Tag by status: "read-later", "important", "reference"

### Automation
```bash
# Cron job for daily validation
0 2 * * * /usr/local/bin/linkmgr validate > /tmp/validation.log

# Cron for daily Hugo update
0 3 * * * cd ~/blog && linkmgr hugo -o content/bookmarks && hugo
```

### Shell Aliases
```bash
# Add to ~/.bashrc or ~/.zshrc
alias lm='linkmgr'
alias lma='linkmgr add'
alias lml='linkmgr list'
alias lms='linkmgr search'
alias lmv='linkmgr validate'
```

---

## Troubleshooting

### Database Connection Issues
```bash
# Test connection
linkmgr init --test-connection

# Check config
cat ~/.linkmgr.yaml

# Verify database exists
mysql -u user -p -e "SHOW DATABASES;"
```

### Validation Issues
```bash
# Increase timeout for slow sites
linkmgr validate --timeout 10

# Reduce concurrency if getting errors
linkmgr validate --concurrency 2
```

### Hugo Export Issues
```bash
# Check output directory exists
mkdir -p content/bookmarks

# Verify with verbose
linkmgr hugo -o content/bookmarks -v
```

---

## Getting Help

```bash
# General help
linkmgr --help

# Command help
linkmgr add --help
linkmgr update --help
linkmgr hugo --help
linkmgr validate --help

# Version
linkmgr version
```

---

## All Commands Reference

1. `init` - Configure database
2. `add` - Add bookmarks
3. `list` - List bookmarks
4. `search` - Search bookmarks
5. `show` - View details
6. `update` - Update bookmark
7. `delete` - Delete bookmarks
8. `tags` - List tags
9. `tag rename` - Rename tag
10. `export` - Export to file
11. `import` - Import from file
12. `validate` - Check links
13. `hugo` - Export to Hugo
14. `version` - Show version
15. `completion` - Shell completion

---

## Next Steps

After mastering the basics:

1. Set up regular validation (cron job)
2. Configure Hugo integration if publishing
3. Create backup routine
4. Standardize your tagging system
5. Build automation workflows

---

Enjoy managing your bookmarks with linkmgr! 🚀
