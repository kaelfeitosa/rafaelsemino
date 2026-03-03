package domain

import (
	"fmt"
	"regexp"
	"sort"
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
	Attachments []Attachment `yaml:"-"`
	Occurrences   []*Occurrence `yaml:"occurrences,omitempty"`
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
	Title         string       `yaml:"title,omitempty"`
	Type          string       `yaml:"type"` // apresentacao | residencia | oficina | publicacao_ou_apresentacao | lancamento | premio | exposicao
	StartDate     string       `yaml:"start_date"`
	EndDate       string       `yaml:"end_date,omitempty"`
	Context       string       `yaml:"context,omitempty"`
	Role          string       `yaml:"role,omitempty"`
	Collaborators []string     `yaml:"collaborators,omitempty"`
	Attachments []Attachment `yaml:"-"`
}

var attachmentRegexp = regexp.MustCompile(`^attachment_(\d+)_(.+)$`)

type Attachable interface {
	GetAttachments() []Attachment
	SetAttachments([]Attachment)
}

func UnmarshalAttachments(value *yaml.Node, a Attachable) error {
	var m map[string]interface{}
	if err := value.Decode(&m); err != nil {
		return err
	}
	a.SetAttachments(ExtractAttachments(m))
	return nil
}

func MarshalAttachments(a Attachable, alias interface{}) (interface{}, error) {
	var node yaml.Node
	if err := node.Encode(alias); err != nil {
		return nil, err
	}
	if err := injectAttachmentsToNode(&node, a.GetAttachments()); err != nil {
		return nil, err
	}
	return node, nil
}


func ExtractAttachments(m map[string]interface{}) []Attachment {
	var attachments []Attachment

	// Support old nested "attachments" list (for migration)
	if rawAtts, ok := m["attachments"].([]interface{}); ok {
		for _, rawAtt := range rawAtts {
			if attMap, ok := rawAtt.(map[string]interface{}); ok {
				var att Attachment
				if v, ok := attMap["type"].(string); ok { att.Type = v }
				if v, ok := attMap["url"].(string); ok { att.URL = v }
				if v, ok := attMap["label"].(string); ok && v != "" { att.Label = v } else if v, ok := attMap["caption"].(string); ok { att.Label = v }
				if v, ok := attMap["category"].(string); ok && v != "" { att.Category = v } else if v, ok := attMap["role"].(string); ok { att.Category = v }
				attachments = append(attachments, att)
			}
		}
	}

	// Support flattened keys
	attFieldsMap := make(map[int]map[string]string)
	for k, v := range m {
		matches := attachmentRegexp.FindStringSubmatch(k)
		if len(matches) == 3 {
			idx, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}
			field := matches[2]

			if _, ok := attFieldsMap[idx]; !ok {
				attFieldsMap[idx] = make(map[string]string)
			}

			strVal, ok := v.(string)
			if !ok {
				if v != nil {
					strVal = fmt.Sprintf("%v", v)
				}
			}
			attFieldsMap[idx][field] = strVal
		}
	}

	attMap := make(map[int]Attachment)
	for idx, fields := range attFieldsMap {
		att := Attachment{
			Type: fields["type"],
			URL:  fields["url"],
		}

		if label, ok := fields["label"]; ok && label != "" {
			att.Label = label
		} else {
			att.Label = fields["caption"]
		}

		if category, ok := fields["category"]; ok && category != "" {
			att.Category = category
		} else {
			att.Category = fields["role"]
		}
		attMap[idx] = att
	}

	if len(attMap) > 0 {
		var indices []int
		for idx := range attMap {
			indices = append(indices, idx)
		}
		sort.Ints(indices)

		for _, idx := range indices {
			attachments = append(attachments, attMap[idx])
		}
	}

	return attachments
}

func injectAttachmentsToNode(node *yaml.Node, attachments []Attachment) error {
	if node.Kind == yaml.DocumentNode {
		node = node.Content[0]
	}
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected MappingNode")
	}

	for i, att := range attachments {
		idx := i + 1
		fields := []struct {
			Name  string
			Value string
		}{
			{"type", att.Type},
			{"url", att.URL},
			{"label", att.Label},
			{"category", att.Category},
		}

		for _, field := range fields {
			if field.Value != "" {
				node.Content = append(node.Content,
					&yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("attachment_%d_%s", idx, field.Name)},
					&yaml.Node{Kind: yaml.ScalarNode, Value: field.Value},
				)
			}
		}
	}
	return nil
}


func (a *Agent) GetAttachments() []Attachment { return a.Attachments }
func (a *Agent) SetAttachments(atts []Attachment) { a.Attachments = atts }

func (w *Work) GetAttachments() []Attachment { return w.Attachments }
func (w *Work) SetAttachments(atts []Attachment) { w.Attachments = atts }

func (o *Occurrence) GetAttachments() []Attachment { return o.Attachments }
func (o *Occurrence) SetAttachments(atts []Attachment) { o.Attachments = atts }

// Agent custom marshal/unmarshal
func (a *Agent) UnmarshalYAML(value *yaml.Node) error {
	type alias Agent
	var aux alias
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*a = Agent(aux)
	return UnmarshalAttachments(value, a)
}

func (a *Agent) MarshalYAML() (interface{}, error) {
	type alias Agent
	if a == nil {
		return nil, nil
	}
	return MarshalAttachments(a, alias(*a))
}

// Work custom marshal/unmarshal
func (w *Work) UnmarshalYAML(value *yaml.Node) error {
	type alias Work
	var aux alias
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*w = Work(aux)
	return UnmarshalAttachments(value, w)
}

func (w *Work) MarshalYAML() (interface{}, error) {
	type alias Work
	if w == nil {
		return nil, nil
	}
	return MarshalAttachments(w, alias(*w))
}

// Occurrence custom marshal/unmarshal
func (o *Occurrence) UnmarshalYAML(value *yaml.Node) error {
	type alias Occurrence
	var aux alias
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*o = Occurrence(aux)
	return UnmarshalAttachments(value, o)
}

func (o *Occurrence) MarshalYAML() (interface{}, error) {
	type alias Occurrence
	if o == nil {
		return nil, nil
	}
	return MarshalAttachments(o, alias(*o))
}
