package topics

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

const (
	defaultTopicFilter = `\[?[a-zA-Z0-9]_*-[0-9]+\]?`
	defautlRepoPath    = "."
)

var (
	version          = "dev"
	topicFilter      string
	outFilename      string
	repositoryPath   string
	templateFilename string
	outputAsJSON     bool
	typesFilename    string
	typeFilters      TypeFilters
	verbose          bool

	rootCmd = &cobra.Command{
		Use:     os.Args[0],
		Short:   "Generates a summary of file changes by topics",
		Long:    `Generates a summary of file changes by topics`,
		RunE:    run,
		Version: version,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&outFilename, "output", "o", "", "output file name (default: print to stdout)")
	rootCmd.Flags().StringVarP(&topicFilter, "filter", "f", defaultTopicFilter, fmt.Sprintf("topic filter (default JIRA: %q)", defaultTopicFilter))
	rootCmd.Flags().StringVarP(&repositoryPath, "repository", "r", ".", "Git repository (default: current working directory)")
	rootCmd.Flags().StringVarP(&templateFilename, "template", "t", "", "template file to generate the report, see Go's text/template package (default: internal markdown template)")
	rootCmd.Flags().BoolVar(&outputAsJSON, "json", false, "output JSON (default: false)")
	rootCmd.Flags().StringVar(&typesFilename, "types", "", "JSON file describing the types")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "be verbose (default: false)")
}

func Execute() error {
	return rootCmd.Execute()
}

func Info(msg string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, msg, args...)
	}
}

func run(cmd *cobra.Command, args []string) error {
	topicRegex, err := regexp.Compile(topicFilter)
	if err != nil {
		return fmt.Errorf("could not compile filter regexp: %s", err)
	}

	typeFilters, err = getTypeFilters()
	if err != nil {
		return err
	}

	r, err := git.PlainOpen(repositoryPath)
	if err != nil {
		return fmt.Errorf("could not open repository in %q: %s", repositoryPath, err)
	}

	ref, err := r.Head()
	if err != nil {
		return fmt.Errorf("could not get HEAD reference: %s", err)
	}

	commits, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return fmt.Errorf("could not get commit logs: %s", err)
	}

	out, err := getOutputStream()
	if err != nil {
		return fmt.Errorf("could not open output stream: %s", err)
	}
	defer out.Close()

	topics, err := extractTopics(commits, *topicRegex)
	if err != nil {
		return fmt.Errorf("could not extract data from commits: %s", err)
	}

	if outputAsJSON {
		return json.NewEncoder(out).Encode(topics)
	}

	templ, err := getTemplate()
	if err != nil {
		return err
	}
	return templ.Execute(out, topics)
}

func getOutputStream() (io.WriteCloser, error) {
	if outFilename == "" {
		return os.Stdout, nil
	}

	_, err := os.Stat(outFilename)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("output file already exists, aborting.")
	}

	outFile, err := os.Create(outFilename)
	if err != nil {
		return nil, fmt.Errorf("could not create file %q: %s", outFilename, err)
	}

	return outFile, nil
}

func getTemplate() (*template.Template, error) {
	if templateFilename == "" {
		// use default template
		return template.New("default").Parse(defaultTemplate)
	}

	f, err := os.Stat(templateFilename)
	if err != nil {
		return nil, err
	}

	if f.IsDir() {
		return nil, fmt.Errorf("%q is a directory, not a file", templateFilename)
	}

	return template.ParseFiles(templateFilename)
}

func getTypeFilters() (TypeFilters, error) {
	var filters TypeFilters
	if typesFilename == "" {
		return filters, nil
	}

	f, err := os.Open(typesFilename)
	if err != nil {
		return filters, err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&filters); err != nil {
		return nil, err
	}

	for _, tf := range filters {
		if err = tf.Compile(); err != nil {
			return filters, fmt.Errorf("could not compile regexp %q for type %q: %s", tf.RegexString, tf.TypeName, err)
		}
	}

	return filters, nil
}

func getComponentType(filename string) string {
	for _, t := range typeFilters {
		if t.Skip {
			continue
		}

		if t.Match(filename) {
			return t.TypeName
		}
	}

	return ""
}

func ignore(filename string) bool {
	for _, t := range typeFilters {
		fmt.Printf("checking %q with %q\n", filename, t.RegexString)
		if t.Match(filename) {
			return t.Skip
		}
	}
	return false
}

func toModifiedFile(f *object.File) *ModifiedFile {
	return &ModifiedFile{
		FullName: f.Name,
		Name:     filepath.Base(f.Name),
		Path:     filepath.Dir(f.Name),
		Type:     getComponentType(f.Name),
	}
}

const defaultTemplate = `# Changes summary per topic

{{range .}}
## {{.Name}}:

| Manifest Item | Type | Description |
|:--------------|:-----|:------------|
{{- range .ModifiedFiles}}
| {{.Name}} | {{.Type}} | {{range .Commits}}{{.ShortHash}} {{.Summary}} {{end}} |
{{- end}}
{{end}}
`
