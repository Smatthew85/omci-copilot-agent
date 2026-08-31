package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"

	generated "github.com/opencord/omci-lib-go/v2/generated"
)

const (
	sourceLibrary = "github.com/opencord/omci-lib-go/v2"
)

type meCatalog struct {
	ClassID      uint             `json:"class_id"`
	Name         string           `json:"name"`
	MessageTypes []string         `json:"message_types,omitempty"`
	Attributes   []attributeEntry `json:"attributes"`
	Source       sourceInfo       `json:"source"`
}

type attributeEntry struct {
	Index     uint     `json:"index"`
	Name      string   `json:"name"`
	SizeBytes int      `json:"size_bytes"`
	Access    []string `json:"access"`
	Mandatory bool     `json:"mandatory"`
	Table     bool     `json:"table"`
}

type sourceInfo struct {
	Library string `json:"library"`
	Version string `json:"version"`
}

type meFile struct {
	ClassID uint
	Name    string
	File    string
	Data    []byte
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func main() {
	out := flag.String("out", "knowledge/me-catalog", "Output directory")
	pretty := flag.Bool("pretty", true, "Pretty-print JSON")
	writeIndex := flag.Bool("index", true, "Also write INDEX.md")
	flag.Parse()

	if err := run(*out, *pretty, *writeIndex); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(out string, pretty bool, writeIndex bool) error {
	if err := os.MkdirAll(out, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := clearGeneratedFiles(out); err != nil {
		return err
	}

	version := resolveLibraryVersion()
	classIDs := generated.GetSupportedClassIDs()
	sort.Slice(classIDs, func(i, j int) bool { return classIDs[i] < classIDs[j] })

	files := make([]meFile, 0, len(classIDs))
	for _, cid := range classIDs {
		me, err := generated.LoadManagedEntityDefinition(cid)
		if err.StatusCode() != generated.Success {
			continue
		}

		catalog := meCatalog{
			ClassID:    uint(me.GetClassID()),
			Name:       humanizeName(me.GetName()),
			Attributes: extractAttributes(me.GetAttributeDefinitions()),
			Source: sourceInfo{
				Library: sourceLibrary,
				Version: version,
			},
		}

		if msgTypes := extractMessageTypes(me.GetMessageTypes()); len(msgTypes) > 0 {
			catalog.MessageTypes = msgTypes
		}

		data, errMarshal := marshalCatalog(catalog, pretty)
		if errMarshal != nil {
			return fmt.Errorf("marshal class %d: %w", cid, errMarshal)
		}

		filename := fmt.Sprintf("%03d-%s.json", uint(cid), slugify(catalog.Name))
		files = append(files, meFile{ClassID: uint(cid), Name: catalog.Name, File: filename, Data: data})
	}

	sort.Slice(files, func(i, j int) bool { return files[i].ClassID < files[j].ClassID })
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(out, file.File), file.Data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", file.File, err)
		}
	}

	if writeIndex {
		if err := writeIndexFile(out, files, version); err != nil {
			return err
		}
	}

	return nil
}

func clearGeneratedFiles(out string) error {
	entries, err := os.ReadDir(out)
	if err != nil {
		return fmt.Errorf("read output directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".json") || name == "INDEX.md" {
			if removeErr := os.Remove(filepath.Join(out, name)); removeErr != nil {
				return fmt.Errorf("remove %s: %w", name, removeErr)
			}
		}
	}
	return nil
}

func extractAttributes(defs generated.AttributeDefinitionMap) []attributeEntry {
	attrs := make([]attributeEntry, 0, len(defs))
	for _, def := range defs {
		table := def.IsTableAttribute() || def.GetSize() == 0
		size := def.GetSize()
		if table {
			size = 0
		}

		attrs = append(attrs, attributeEntry{
			Index:     def.GetIndex(),
			Name:      humanizeName(def.GetName()),
			SizeBytes: size,
			Access:    extractAccess(def),
			Mandatory: !def.Optional,
			Table:     table,
		})
	}
	sort.Slice(attrs, func(i, j int) bool { return attrs[i].Index < attrs[j].Index })
	return attrs
}

func extractAccess(def generated.AttributeDefinition) []string {
	access := make([]string, 0, 4)
	if generated.SupportsAttributeAccess(def, generated.Read) {
		access = append(access, "Read")
	}
	if generated.SupportsAttributeAccess(def, generated.Write) {
		access = append(access, "Write")
	}
	if generated.SupportsAttributeAccess(def, generated.SetByCreate) {
		access = append(access, "SetByCreate")
	}
	return access
}

func extractMessageTypes(set interface{ ToSlice() []interface{} }) []string {
	values := set.ToSlice()
	result := make([]string, 0, len(values))
	for _, v := range values {
		msg, ok := v.(generated.MsgType)
		if !ok {
			continue
		}
		name := strings.ReplaceAll(msg.String(), " ", "")
		if name == "Unknown" {
			continue
		}
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func marshalCatalog(c meCatalog, pretty bool) ([]byte, error) {
	if pretty {
		out, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			return nil, err
		}
		return append(out, '\n'), nil
	}
	out, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}

func writeIndexFile(out string, files []meFile, version string) error {
	var b bytes.Buffer
	b.WriteString("# ME Catalog Index\n\n")
	b.WriteString(fmt.Sprintf("Extracted from `%s` (`%s`). Total: `%d` MEs.\n\n", sourceLibrary, version, len(files)))
	b.WriteString("| Class ID | Name | File |\n")
	b.WriteString("|---|---|---|\n")
	for _, file := range files {
		b.WriteString(fmt.Sprintf("| %d | %s | [%s](./%s) |\n", file.ClassID, file.Name, file.File, file.File))
	}
	return os.WriteFile(filepath.Join(out, "INDEX.md"), b.Bytes(), 0o644)
}

func resolveLibraryVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	for _, dep := range info.Deps {
		if dep.Path == sourceLibrary {
			if dep.Replace != nil {
				return dep.Replace.Version
			}
			return dep.Version
		}
	}
	return "unknown"
}

func humanizeName(name string) string {
	name = strings.ReplaceAll(name, "_", " ")
	runes := []rune(name)
	var out []rune
	for i, r := range runes {
		if i > 0 {
			next := rune(0)
			hasNext := i < len(runes)-1
			if hasNext {
				next = runes[i+1]
			}
			if isBoundary(runes[i-1], r, next, hasNext) {
				out = append(out, ' ')
			}
		}
		out = append(out, r)
	}
	parts := strings.Fields(string(out))
	for i, part := range parts {
		parts[i] = normalizeWord(part)
	}
	return strings.Join(parts, " ")
}

func isBoundary(prev, curr, next rune, hasNext bool) bool {
	if prev == ' ' || curr == ' ' {
		return false
	}
	if isLower(prev) && isUpper(curr) {
		return true
	}
	if isLetter(prev) && isDigit(curr) {
		return true
	}
	if isDigit(prev) && isLetter(curr) {
		return true
	}
	if isUpper(prev) && isUpper(curr) && hasNext && isLower(next) {
		return true
	}
	return false
}

func normalizeWord(word string) string {
	if word == "" {
		return ""
	}
	allCaps := true
	for _, r := range word {
		if isLetter(r) && !isUpper(r) {
			allCaps = false
			break
		}
	}
	if allCaps {
		return word
	}
	lower := strings.ToLower(word)
	return strings.ToUpper(lower[:1]) + lower[1:]
}

func slugify(name string) string {
	slug := strings.ToLower(name)
	slug = nonAlnum.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "unknown"
	}
	return slug
}

func isLower(r rune) bool  { return 'a' <= r && r <= 'z' }
func isUpper(r rune) bool  { return 'A' <= r && r <= 'Z' }
func isDigit(r rune) bool  { return '0' <= r && r <= '9' }
func isLetter(r rune) bool { return isLower(r) || isUpper(r) }
