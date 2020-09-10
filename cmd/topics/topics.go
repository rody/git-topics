package topics

import (
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type Topic struct {
	Name          string
	ModifiedFiles map[string]*ModifiedFile
}

type ModifiedFile struct {
	Name     string
	Path     string
	FullName string
	Type     string
	Commits  []*CommitInfo
}

type CommitInfo struct {
	Author    object.Signature
	Committer object.Signature
	Hash      string
	ShortHash string
	Summary   string
	Message   string
}

type TypeFilters []*TypeEntry

type TypeEntry struct {
	RegexString string `json:"regex"`
	TypeName    string `json:"type"`
	Skip        bool   `json:"skip"`

	regexp *regexp.Regexp
}

func (t *TypeEntry) Compile() error {
	rg, err := regexp.Compile(t.RegexString)
	if err != nil {
		return err
	}
	t.regexp = rg
	return nil
}

func (t *TypeEntry) Match(filename string) bool {
	return t.regexp.MatchString(filename)
}

// toCommitInfo transform commit data into a CommitInfo
func toCommitInfo(c *object.Commit) *CommitInfo {
	return &CommitInfo{
		Author:    c.Author,
		Committer: c.Committer,
		Hash:      c.Hash.String(),
		ShortHash: c.Hash.String()[:7],
		Summary:   strings.Split(c.Message, "\n")[0],
		Message:   c.Message,
	}
}

// extractTopics reads the git repository and return the topic info
func extractTopics(commits object.CommitIter, filter regexp.Regexp) ([]Topic, error) {
	topicMap := make(map[string]Topic)

	err := commits.ForEach(func(c *object.Commit) error {
		Info("looking at commit %s\n", c.Hash)

		if filter.MatchString(c.Message) {
			topicName := filter.FindString(c.Message)
			topicKey := strings.ToUpper(topicName)
			topic, topicExists := topicMap[topicKey]
			if !topicExists {
				topic = Topic{
					Name:          topicName,
					ModifiedFiles: make(map[string]*ModifiedFile),
				}
				topicMap[topicKey] = topic
			}

			fIter, err := c.Files()
			if err != nil {
				return err
			}

			return fIter.ForEach(func(f *object.File) error {
				if ignore(f.Name) {
					return nil
				}

				mf, ok := topic.ModifiedFiles[f.Name]
				if !ok {
					mf = toModifiedFile(f)
					topic.ModifiedFiles[f.Name] = mf
				}
				mf.Commits = append(mf.Commits, toCommitInfo(c))
				topic.ModifiedFiles[f.Name] = mf
				return nil
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := make([]Topic, len(topicMap))
	i := 0
	for _, t := range topicMap {
		result[i] = t
		i++
	}
	return result, nil
}
