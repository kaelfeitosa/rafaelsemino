package domain

import (
	"fmt"
	"regexp"
	"strconv"
	"gopkg.in/yaml.v3"
)

type Agent struct {
	ID          string       `yaml:"id"`
	Name        string       `yaml:"name"`
	Kind        string       `yaml:"kind"` // person | collective
	Description string       `yaml:"description,omitempty"`
	FoundedByMe bool         `yaml:"founded_by_me,omitempty"`
	ActiveSince int          `yaml:"active_since,omitempty"`
	Links       []Link       `yaml:"links,omitempty"`
	Attachments []Attachment `yaml:"-"`
	Featured    bool         `yaml:"featured,omitempty"`
}

type Link struct {
	Type string `yaml:"type"` // site | instagram | linktree | outro
	URL  string `yaml:"url"`
}

type Work struct {
	ID            string       `yaml:"id"`
	Title         string       `yaml:"title"`
	Medium        string       `yaml:"medium"` // teatro | audiovisual | pesquisa | ensino | formacao | exposicao | outro
	Description   string       `yaml:"description,omitempty"`
	Year          interface{}  `yaml:"year,omitempty"` // interface to handle int or string
	Role          string       `yaml:"role,omitempty"` // For standalone works where I took a role
	Collaborators []string     `yaml:"collaborators,omitempty"`
	Attachments   []Attachment `yaml:"-"`
	Occurrences   []Occurrence `yaml:"occurrences,omitempty"`
	Featured      bool         `yaml:"featured,omitempty"`
}

type Attachment struct {
	Type     string `yaml:"type"` // image | video | pdf | link
	URL      string `yaml:"url"`
	Label    string `yaml:"label,omitempty"`
	Category string `yaml:"category,omitempty"` // documentation | poster | clipping | program | technical | outro
}

// Occurrence represents an event linked to a Work (e.g., season, presentation)
type Occurrence struct {
	Title         string       `yaml:"title"`
	Type          string       `yaml:"type"` // apresentacao | residencia | oficina | publicacao_ou_apresentacao | lancamento | premio | exposicao
	StartDate     string       `yaml:"start_date"`
	EndDate       string       `yaml:"end_date,omitempty"`
	Context       string       `yaml:"context,omitempty"`
	Role          string       `yaml:"role,omitempty"`
	Collaborators []string     `yaml:"collaborators,omitempty"`
	Attachments   []Attachment `yaml:"-"`
}

func extractAttachments(value *yaml.Node) ([]Attachment, error) {
	var attachments []Attachment
	attMap := make(map[int]*Attachment)

	if value.Kind != yaml.MappingNode {
		return nil, nil
	}

	re := regexp.MustCompile(`^attachment_(\d+)_(.+)$`)

	for i := 0; i < len(value.Content); i += 2 {
		keyNode := value.Content[i]
		valNode := value.Content[i+1]

		matches := re.FindStringSubmatch(keyNode.Value)
		if len(matches) == 3 {
			idx, _ := strconv.Atoi(matches[1])
			field := matches[2]

			if attMap[idx] == nil {
				attMap[idx] = &Attachment{}
			}

			switch field {
			case "type":
				attMap[idx].Type = valNode.Value
			case "url":
				attMap[idx].URL = valNode.Value
			case "label":
				attMap[idx].Label = valNode.Value
			case "category":
				attMap[idx].Category = valNode.Value
			}
		}
	}

	maxIdx := 0
	for idx := range attMap {
		if idx > maxIdx {
			maxIdx = idx
		}
	}

	for i := 1; i <= maxIdx; i++ {
		if att, ok := attMap[i]; ok {
			attachments = append(attachments, *att)
		}
	}

	return attachments, nil
}

func appendAttachmentsToNode(alias interface{}, attachments []Attachment) (*yaml.Node, error) {
	node := &yaml.Node{}
	if err := node.Encode(alias); err != nil {
		return nil, err
	}

	for i, att := range attachments {
		idx := i + 1

		if att.Type != "" {
			node.Content = append(node.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_type", idx)},
				&yaml.Node{Kind: yaml.ScalarNode, Value: att.Type},
			)
		}
		if att.URL != "" {
			node.Content = append(node.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_url", idx)},
				&yaml.Node{Kind: yaml.ScalarNode, Value: att.URL},
			)
		}
		if att.Label != "" {
			node.Content = append(node.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_label", idx)},
				&yaml.Node{Kind: yaml.ScalarNode, Value: att.Label},
			)
		}
		if att.Category != "" {
			node.Content = append(node.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_category", idx)},
				&yaml.Node{Kind: yaml.ScalarNode, Value: att.Category},
			)
		}
	}

	return node, nil
}

func (a *Agent) UnmarshalYAML(value *yaml.Node) error {
	type agentAlias Agent
	var alias agentAlias
	if err := value.Decode(&alias); err != nil {
		return err
	}
	attachments, err := extractAttachments(value)
	if err != nil {
		return err
	}
	alias.Attachments = attachments
	*a = Agent(alias)
	return nil
}

func (a Agent) MarshalYAML() (interface{}, error) {
	type agentAlias Agent
	alias := agentAlias(a)
	return appendAttachmentsToNode(&alias, a.Attachments)
}

func (w *Work) UnmarshalYAML(value *yaml.Node) error {
	type workAlias Work
	var alias workAlias
	if err := value.Decode(&alias); err != nil {
		return err
	}
	attachments, err := extractAttachments(value)
	if err != nil {
		return err
	}
	alias.Attachments = attachments
	*w = Work(alias)
	return nil
}

func (w Work) MarshalYAML() (interface{}, error) {
	type workAlias Work
	alias := workAlias(w)
	return appendAttachmentsToNode(&alias, w.Attachments)
}

func (o *Occurrence) UnmarshalYAML(value *yaml.Node) error {
	type occurrenceAlias Occurrence
	var alias occurrenceAlias
	if err := value.Decode(&alias); err != nil {
		return err
	}
	attachments, err := extractAttachments(value)
	if err != nil {
		return err
	}
	alias.Attachments = attachments
	*o = Occurrence(alias)
	return nil
}

func (o Occurrence) MarshalYAML() (interface{}, error) {
	type occurrenceAlias Occurrence
	alias := occurrenceAlias(o)
	return appendAttachmentsToNode(&alias, o.Attachments)
}
