// ./cmd/lg/main.go
package main

import (
 "flag"
 "fmt"
 "io/fs"
 "os"
 "path/filepath"
 "sort"
 "strings"
)

// Config from CLI flags
type Config struct {
 Color     bool
 Directory string
 Recursive bool
 Long      bool
 All       bool
}

func main() {
 var cfg Config
 flag.BoolVar(&cfg.Color, "c", false, "Enable RGBA/truecolor (ANSI) colored output")
 flag.StringVar(&cfg.Directory, "d", "", "Directory to list (default: current directory)")
 flag.BoolVar(&cfg.Recursive, "s", false, "Show sub-directories (recursive)")
 flag.BoolVar(&cfg.Long, "l", false, "Long/verbose output (permissions, size, modified, type)")
 flag.BoolVar(&cfg.All, "a", false, "Include ALL files (including hidden files and directories)")
 flag.Usage = func() {
  fmt.Fprintf(flag.CommandLine.Output(), "usage:\n  lg [OPTIONS] [directory]\n\nOptions:\n")
  flag.PrintDefaults()
  fmt.Fprintf(flag.CommandLine.Output(), `
Examples:
  lg -d /path/to/dir
  lg -dasl ~/
`)
 }
 flag.Parse()

 // positional arg overrides -d if provided
 if cfg.Directory == "" && flag.NArg() > 0 {
  cfg.Directory = flag.Arg(0)
 }

 if cfg.Directory == "" {
  pwd, err := os.Getwd()
  if err != nil {
   fmt.Fprintf(os.Stderr, "error: cannot get current directory: %v\n", err)
   os.Exit(1)
  }
  cfg.Directory = pwd
 }

 // Clean path and expand ~ if present
 cfg.Directory = expandPath(cfg.Directory)

 // Print root and then list children
 rootInfo, err := os.Stat(cfg.Directory)
 if err != nil {
  fmt.Fprintf(os.Stderr, "error: cannot stat %s: %v\n", cfg.Directory, err)
  os.Exit(1)
 }

 // Root line
 fmt.Println(formatRootLine(cfg.Directory, rootInfo, cfg))

 // list top-level entries
 if err := printChildren(cfg.Directory, 1, cfg); err != nil {
  fmt.Fprintf(os.Stderr, "error: %v\n", err)
  os.Exit(1)
 }
}

// expandPath handles ~ expansion for unix-like shells
func expandPath(p string) string {
 if strings.HasPrefix(p, "~") {
  home, err := os.UserHomeDir()
  if err == nil {
   return filepath.Join(home, strings.TrimPrefix(p, "~"))
  }
 }
 return p
}

func formatRootLine(path string, info os.FileInfo, cfg Config) string {
 display := path
 if display == "" {
  display = "."
 }
 if cfg.Long {
  meta := longMeta(info, path)
  if cfg.Color {
   return colorWrap(display+"/", 0, 120, 255) + " " + meta
  }
  return display + "/" + " " + meta
 }
 if cfg.Color {
  return colorWrap(display+"/", 0, 120, 255)
 }
 return display + "/"
}

// printChildren lists a directory's entries and, if requested, recurses.
func printChildren(dir string, depth int, cfg Config) error {
 entries, err := os.ReadDir(dir)
 if err != nil {
  return fmt.Errorf("reading dir %s: %w", dir, err)
 }

 // Filter hidden files if not cfg.All
 filtered := make([]fs.DirEntry, 0, len(entries))
 for _, e := range entries {
  name := e.Name()
  if !cfg.All && strings.HasPrefix(name, ".") {
   continue
  }
  filtered = append(filtered, e)
 }

 // Sort directories first then files, alphabetically (case-insensitive)
 sort.Slice(filtered, func(i, j int) bool {
  a, b := filtered[i], filtered[j]
  if a.IsDir() && !b.IsDir() {
   return true
  }
  if !a.IsDir() && b.IsDir() {
   return false
  }
  return strings.ToLower(a.Name()) < strings.ToLower(b.Name())
 })

 for _, e := range filtered {
  full := filepath.Join(dir, e.Name())
  indent := prefixForDepth(depth)
  line := fmt.Sprintf("%s%s", indent, e.Name())

  // long info if requested
  if cfg.Long {
   info, err := os.Lstat(full)
   if err == nil {
    meta := longMeta(info, full)
    line = fmt.Sprintf("%s %s", line, meta)
   }
  }

  // color
  if cfg.Color {
   colored := colorizeEntry(full, e)
   // use meta or name colored accordingly
   if cfg.Long {
    // split name from meta
    parts := strings.SplitN(line, " ", 2)
    if len(parts) == 2 {
      line = fmt.Sprintf("%s %s", colored, parts[1])
    } else {
      line = colored
    }
   } else {
    line = colored
   }
  }

  fmt.Println(line)

  // if directory and recursive flag set, recurse
  if e.IsDir() && cfg.Recursive {
   if err := printChildren(full, depth+1, cfg); err != nil {
    // continue on errors but report them
    fmt.Fprintf(os.Stderr, "warning: %v\n", err)
   }
  }
 }

 return nil
}

// prefixForDepth builds the textual tree prefix used in the sample: "| - ", "| - - ", ...
func prefixForDepth(depth int) string {
 // depth 1 -> "| - "
 // depth 2 -> "| - - "
 if depth <= 0 {
  depth = 1
 }
 return "| " + strings.Repeat("- ", depth)
}

// longMeta formats permissions, size, modtime, and type
func longMeta(info os.FileInfo, path string) string {
 mode := info.Mode()
 perm := mode.Perm().String()
 size := humanSize(info.Size())
 mod := info.ModTime().Format("2006-01-02 15:04:05")
 typ := fileType(mode)
 // For symlink show target if possible
 if mode&os.ModeSymlink != 0 {
  if target, err := os.Readlink(path); err == nil {
   typ = typ + "->" + target
  }
 }
 return fmt.Sprintf("[%s] (%s) %s %s", perm, size, typ, mod)
}

func fileType(mode os.FileMode) string {
 switch {
 case mode.IsDir():
  return "dir"
 case mode&os.ModeSymlink != 0:
  return "symlink"
 case mode.IsRegular():
  return "file"
 default:
  return "other"
 }
}

func humanSize(n int64) string {
 if n < 1024 {
  return fmt.Sprintf("%dB", n)
 }
 k := float64(n) / 1024.0
 if k < 1024 {
  return fmt.Sprintf("%.1fK", k)
 }
 m := k / 1024.0
 if m < 1024 {
  return fmt.Sprintf("%.1fM", m)
 }
 g := m / 1024.0
 return fmt.Sprintf("%.1fG", g)
}

// colorizeEntry chooses a color based on entry type and returns the colored name
func colorizeEntry(full string, e fs.DirEntry) string {
 name := e.Name()
 // Determine file mode for executable/symlink check
 info, err := e.Info()
 if err != nil {
  // fallback: directory vs file
  if e.IsDir() {
   return colorWrap(name, 0, 120, 255) // blue dir
  }
  return colorWrap(name, 200, 200, 200) // grey file
 }
 mode := info.Mode()
 // Directory
 if mode.IsDir() {
  return colorWrap(name, 0, 120, 255)
 }
 // Symlink
 if mode&os.ModeSymlink != 0 {
  return colorWrap(name, 0, 200, 200)
 }
 // Executable
 if mode&0111 != 0 {
  return colorWrap(name, 0, 180, 0)
 }
 // Regular file
 return colorWrap(name, 200, 200, 200)
}

// colorWrap returns a 24-bit ANSI colored string (alpha ignored).
// r,g,b values in 0..255. If running in a terminal that doesn't support ANSI, this still prints escape sequences.
// The user flag -c enables color.
func colorWrap(s string, r, g, b int) string {
 // 24-bit foreground: \x1b[38;2;<r>;<g>;<b>m
 return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, s)
}