# Git Topics

Generate git commit report grouped by topics.

## Usage

``` 
Generates a summary of file changes by topics

Usage:
  git-topics [flags]

Flags:
  -f, --filter string       topic filter (default JIRA: "\\[?[a-zA-Z0-9]_*-[0-9]+\\]?") (default "\\[?[a-zA-Z0-9]_*-[0-9]+\\]?")
  -h, --help                help for ./dist/git-topics_darwin_amd64
      --json                output JSON (default: false)
  -o, --output string       output file name (default: print to stdout)
  -r, --repository string   Git repository (default: current working directory) (default ".")
  -t, --template string     template file to generate the report, see Go's text/template package (default: internal markdown template)
      --types string        JSON file describing the types
      --verbose             be verbose (default: false)
  -v, --version             version for ./dist/git-topics_darwin_amd64

```

## Export to Word, Excel, ....

The default template generates a Markdown document. You can easily transform this into another
format with pandoc.

Install [Pandoc](https://pandoc.org).

``` shell
pandoc -s output.md -o output.docx
```

For more examples, visit https://pandoc.org/demos.html

