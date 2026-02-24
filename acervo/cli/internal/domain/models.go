package domain

type Agent struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Kind        string `yaml:"kind"` // person | collective
	Description string `yaml:"description,omitempty"`
	FoundedByMe bool   `yaml:"founded_by_me,omitempty"`
	ActiveSince int    `yaml:"active_since,omitempty"`
	Links       []Link `yaml:"links,omitempty"`
	Featured    bool   `yaml:"featured,omitempty"`
}

type Link struct {
	Type string `yaml:"type"` // site | instagram | linktree | outro
	URL  string `yaml:"url"`
}

type Work struct {
	ID          string       `yaml:"id"`
	Title       string       `yaml:"title"`
	Type        string       `yaml:"type"` // teatro | jogo | filme | roteiro | performance | outro
	Description string       `yaml:"description,omitempty"`
	Year        int          `yaml:"year,omitempty"`
	Attachments []Attachment `yaml:"attachments,omitempty"`
	Featured    bool         `yaml:"featured,omitempty"`
}

type Context struct {
	Label    string `yaml:"label"`
	Kind     string `yaml:"kind"` // festival | mostra | curso | oficina | residencia | premiacao | entrevista | outro
	Location string `yaml:"location,omitempty"`
	Year     int    `yaml:"year,omitempty"`
}

type Attachment struct {
	Type    string `yaml:"type"` // image | video | pdf | link
	Role    string `yaml:"role"` // documentation | clipping | press | certificate | contract | outro
	Source  string `yaml:"source,omitempty"`
	Src     string `yaml:"src"`
	Caption string `yaml:"caption,omitempty"`
}

type Collaborator struct {
	Name string `yaml:"name"`
	Role string `yaml:"role,omitempty"`
}

type Action struct {
	ID            string         `yaml:"id"`
	Title         string         `yaml:"title"`
	Kind          string         `yaml:"kind"`         // criacao | exibicao | formacao | avaliacao | curadoria | premiacao | outro
	PerformedBy   string         `yaml:"performed_by"` // Agent.id
	MyRole        string         `yaml:"my_role"`
	WorkID        string         `yaml:"work_id,omitempty"` // Work.id
	Context       Context        `yaml:"context"`
	DateStart     string         `yaml:"date_start"`
	DateEnd       string         `yaml:"date_end,omitempty"`
	Description   string         `yaml:"description,omitempty"`
	Collaborators []Collaborator `yaml:"collaborators,omitempty"`
	Attachments   []Attachment   `yaml:"attachments,omitempty"`
	Featured      bool           `yaml:"featured,omitempty"`
}
