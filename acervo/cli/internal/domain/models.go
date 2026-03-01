package domain

type Agent struct {
	ID          string       `yaml:"id"`
	Name        string       `yaml:"name"`
	Kind        string       `yaml:"kind"` // person | collective
	Description string       `yaml:"description,omitempty"`
	FoundedByMe bool         `yaml:"founded_by_me,omitempty"`
	ActiveSince int          `yaml:"active_since,omitempty"`
	Links       []Link       `yaml:"links,omitempty"`
	Attachments []Attachment `yaml:"attachments,omitempty"`
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
	Attachments   []Attachment `yaml:"attachments,omitempty"`
	Occurrences   []Occurrence `yaml:"occurrences,omitempty"`
	Featured      bool         `yaml:"featured,omitempty"`
}

type Attachment struct {
	Type  string `yaml:"type"` // image | video | pdf | link
	URL   string `yaml:"url"`
	Label string `yaml:"label,omitempty"`
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
	Attachments   []Attachment `yaml:"attachments,omitempty"`
}
