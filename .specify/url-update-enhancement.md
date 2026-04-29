# URL Update Enhancement

## Date: 2026-04-29

## Summary

Added the ability to update bookmark URLs using the `update` command.

---

## Changes Made

### Modified File
- `internal/cmd/update.go` - Added URL update capability

### New Functionality
- Added `--url` / `-u` flag to update command
- URL validation before updating
- Shows updated URL in success message

---

## Usage

### Basic URL Update
```bash
linkmgr update 5 --url "https://new-url.com"
```

### Common Use Cases

#### Fix Broken Links
```bash
# 1. Validate to find broken links
linkmgr validate

# 2. Update broken link to working URL
linkmgr update 10 --url "https://working-url.com"
```

#### Fix Redirects
```bash
# 1. Validate shows HTTP 301/302 redirects
linkmgr validate

# Output shows:
# ⚠ [#5] https://old-url.com - HTTP 301 (0.23s)

# 2. Update to final destination URL
linkmgr update 5 --url "https://final-url.com"
```

#### Update URL and Other Fields
```bash
# Update URL and title together
linkmgr update 5 --url "https://new-url.com" --title "Updated Title"

# Update URL and tags
linkmgr update 5 --url "https://new-url.com" --add-tags "updated"
```

---

## Implementation Details

### URL Validation
- Uses existing `validator.ValidateURL()` function
- Ensures URL has proper scheme (http:// or https://)
- Validates URL format before updating database

### Error Handling
- Invalid URL format → Clear error message with suggestion
- Non-existent bookmark → Error with suggestion to list bookmarks
- Database errors → Proper error reporting

### Success Feedback
- Shows bookmark ID
- Displays updated title
- **Displays updated URL** (new)
- Shows updated tags if modified

---

## Example Output

```bash
$ linkmgr update 5 --url "https://golang.org"

✓ Bookmark #5 updated
Title: Go Programming Language
URL:   https://golang.org
```

---

## Benefits

### Link Maintenance
- Fix broken links without deleting and re-adding
- Update redirected URLs to final destinations
- Preserve all metadata (tags, excerpts, timestamps)
- Preserve bookmark ID (important for references)

### Workflow Integration
```bash
# Complete link maintenance workflow
linkmgr validate                          # Find issues
linkmgr show 5                           # Review bookmark
linkmgr update 5 --url "https://new.com" # Fix URL
linkmgr validate --tag important         # Verify fix
```

### Data Integrity
- Maintains referential integrity (tags stay linked)
- Preserves created_at timestamp
- Updates modified_at timestamp automatically
- No need to delete and recreate bookmark

---

## Testing

### Validation Tests
- ✅ Valid URL update works
- ✅ Invalid URL rejected with clear error
- ✅ URL validation happens before database update
- ✅ Help text updated with URL flag
- ✅ Success message shows new URL

### Quality Checks
- ✅ `go fmt` - Code formatted
- ✅ `go vet` - No issues
- ✅ Builds successfully
- ✅ Help text comprehensive

---

## Updated Documentation

### README.md
- ✅ Added URL update example in Quick Start
- ✅ Updated update command documentation
- ✅ Added workflow example (validate + update)
- ✅ Updated command reference with --url flag

### Help Text
- ✅ Added URL example to long help
- ✅ Flag documented with description
- ✅ Shows in flag list with `-u` shorthand

---

## Related Features

This enhancement complements existing features:

### Works with validate command
```bash
linkmgr validate           # Find broken/redirected URLs
linkmgr update 5 --url ... # Fix them
```

### Works with show command
```bash
linkmgr show 5            # Review current URL
linkmgr update 5 --url ... # Update it
```

### Works with tag operations
```bash
# Update URL and tag it as fixed
linkmgr update 5 --url "https://new.com" --add-tags "verified"
```

---

## Conclusion

✅ **URL update capability added**  
✅ **Integrates seamlessly with validation workflow**  
✅ **Preserves data integrity**  
✅ **User-friendly with validation and feedback**

This completes the bookmark management feature set by allowing users to maintain their bookmark URLs without losing metadata or references.
