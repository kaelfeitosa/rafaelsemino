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
	Attachments []Attachment `yaml:"-"`
}



func extractAttachments(m map[string]interface{}) []Attachment {
	var attachments []Attachment

	// Support old nested "attachments" list (for migration)
	if rawAtts, ok := m["attachments"].([]interface{}); ok {
		for _, rawAtt := range rawAtts {
			if attMap, ok := rawAtt.(map[string]interface{}); ok {
				var att Attachment
				if v, ok := attMap["type"].(string); ok { att.Type = v }
				if v, ok := attMap["url"].(string); ok { att.URL = v }
				if v, ok := attMap["label"].(string); ok { att.Label = v }
				if v, ok := attMap["category"].(string); ok { att.Category = v }
				attachments = append(attachments, att)
			}
		}
	}

	// Support flattened keys
	attMap := make(map[int]Attachment)
	re := regexp.MustCompile(`^attachment_(\d+)_(.+)$`)

	for k, v := range m {
		matches := re.FindStringSubmatch(k)
		if len(matches) == 3 {
			idx, _ := strconv.Atoi(matches[1])
			field := matches[2]

			att := attMap[idx]
			strVal, _ := v.(string)
			switch field {
			case "type":
				att.Type = strVal
			case "url":
				att.URL = strVal
			case "label":
				att.Label = strVal
			case "category":
				att.Category = strVal
			}
			attMap[idx] = att
		}
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

func injectAttachments(m map[string]interface{}, attachments []Attachment) {
	for i, att := range attachments {
		idx := i + 1
		if att.Type != "" {
			m[fmt.Sprintf("attachment_%d_type", idx)] = att.Type
		}
		if att.URL != "" {
			m[fmt.Sprintf("attachment_%d_url", idx)] = att.URL
		}
		if att.Label != "" {
			m[fmt.Sprintf("attachment_%d_label", idx)] = att.Label
		}
		if att.Category != "" {
			m[fmt.Sprintf("attachment_%d_category", idx)] = att.Category
		}
	}
}

// Agent custom marshal/unmarshal
func (a *Agent) UnmarshalYAML(value *yaml.Node) error {
	type alias Agent
	var aux alias
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*a = Agent(aux)

	var m map[string]interface{}
	if err := value.Decode(&m); err != nil {
		return err
	}
	a.Attachments = extractAttachments(m)
	return nil
}

func (a Agent) MarshalYAML() (interface{}, error) {
	type alias Agent
	bytes, err := yaml.Marshal(alias(a))
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := yaml.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	injectAttachments(m, a.Attachments)
	return m, nil
}

// Work custom marshal/unmarshal
func (w *Work) UnmarshalYAML(value *yaml.Node) error {
	type alias Work
	var aux alias
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*w = Work(aux)

	var m map[string]interface{}
	if err := value.Decode(&m); err != nil {
		return err
	}
	w.Attachments = extractAttachments(m)
	return nil
}

func (w Work) MarshalYAML() (interface{}, error) {
	type alias Work
	bytes, err := yaml.Marshal(alias(w))
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := yaml.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	injectAttachments(m, w.Attachments)
	return m, nil
}

// Occurrence custom marshal/unmarshal
func (o *Occurrence) UnmarshalYAML(value *yaml.Node) error {
	type alias Occurrence
	var aux alias
	if err := value.Decode(&aux); err != nil {
		return err
	}
	*o = Occurrence(aux)

	var m map[string]interface{}
	if err := value.Decode(&m); err != nil {
		return err
	}
	o.Attachments = extractAttachments(m)
	return nil
}

func (o Occurrence) MarshalYAML() (interface{}, error) {
	type alias Occurrence
	bytes, err := yaml.Marshal(alias(o))
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := yaml.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	injectAttachments(m, o.Attachments)
	return m, nil
}
