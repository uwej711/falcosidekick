package outputs

import (
	"strings"

	"github.com/falcosecurity/falcosidekick/types"
)

type teamsFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type teamsSection struct {
	ActivityTitle    string      `json:"activityTitle"`
	ActivitySubTitle string      `json:"activitySubtitle"`
	ActivityImage    string      `json:"activityImage,omitempty"`
	Text             string      `json:"text"`
	Facts            []teamsFact `json:"facts,omitempty"`
}

// Payload
type teamsPayload struct {
	Type       string         `json:"@type"`
	Summary    string         `json:"summary,omitempty"`
	ThemeColor string         `json:"themeColor,omitempty"`
	Sections   []teamsSection `json:"sections"`
}

func newTeamsPayload(falcopayload types.FalcoPayload, config *types.Configuration) teamsPayload {
	var sections []teamsSection
	var section teamsSection
	var facts []teamsFact
	var fact teamsFact

	section.ActivityTitle = "Falco Sidekick"
	section.ActivitySubTitle = falcopayload.Time.String()

	if config.Teams.OutputFormat == All || config.Teams.OutputFormat == "text" || config.Teams.OutputFormat == "" {
		section.Text = falcopayload.Output
	}

	if config.Teams.ActivityImage != "" {
		section.ActivityImage = config.Teams.ActivityImage
	}

	if config.Teams.OutputFormat == All || config.Teams.OutputFormat == "facts" || config.Teams.OutputFormat == "" {
		for i, j := range falcopayload.OutputFields {
			switch j.(type) {
			case string:
				fact.Name = i
				fact.Value = j.(string)
			default:
				continue
			}
			facts = append(facts, fact)
		}

		fact.Name = Rule
		fact.Value = falcopayload.Rule
		facts = append(facts, fact)
		fact.Name = Priority
		fact.Value = falcopayload.Priority
		facts = append(facts, fact)
	}

	section.Facts = facts

	var color string
	switch strings.ToLower(falcopayload.Priority) {
	case Emergency:
		color = "e20b0b"
	case Alert:
		color = "ff5400"
	case Critical:
		color = "ff9000"
	case Error:
		color = "ffc700"
	case Warning:
		color = "ffff00"
	case Notice:
		color = "5bffb5"
	case Informational:
		color = "68c2ff"
	case Debug:
		color = "ccfff2"
	}

	sections = append(sections, section)

	t := teamsPayload{
		Type:       "MessageCard",
		Summary:    falcopayload.Output,
		ThemeColor: color,
		Sections:   sections,
	}

	return t
}

// TeamsPost posts event to Teams
func (c *Client) TeamsPost(falcopayload types.FalcoPayload) {
	err := c.Post(newTeamsPayload(falcopayload, c.Config))
	if err != nil {
		c.Stats.Teams.Add(Error, 1)
		c.PromStats.Outputs.With(map[string]string{"destination": "teams", "status": Error}).Inc()
	} else {
		c.Stats.Teams.Add(OK, 1)
		c.PromStats.Outputs.With(map[string]string{"destination": "teams", "status": OK}).Inc()
	}

	c.Stats.Teams.Add(Total, 1)
}
