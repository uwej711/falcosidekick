package outputs

import (
	"strings"

	"github.com/falcosecurity/falcosidekick/types"
)

type influxdbPayload string

func newInfluxdbPayload(falcopayload types.FalcoPayload, config *types.Configuration) influxdbPayload {
	s := "events,rule=" + strings.Replace(falcopayload.Rule, " ", "_", -1) + ",priority=" + strings.Replace(falcopayload.Priority, " ", "_", -1)

	for i, j := range falcopayload.OutputFields {
		switch j.(type) {
		case string:
			s += "," + i + "=" + strings.Replace(j.(string), " ", "_", -1)
		default:
			continue
		}
	}

	s += " value=\"" + falcopayload.Output + "\""

	return influxdbPayload(s)
}

// InfluxdbPost posts event to InfluxDB
func (c *Client) InfluxdbPost(falcopayload types.FalcoPayload) {
	err := c.Post(newInfluxdbPayload(falcopayload, c.Config))
	if err != nil {
		c.Stats.Influxdb.Add(Error, 1)
		c.PromStats.Outputs.With(map[string]string{"destination": "influxdb", "status": Error}).Inc()
	} else {
		c.Stats.Influxdb.Add(OK, 1)
		c.PromStats.Outputs.With(map[string]string{"destination": "influxdb", "status": OK}).Inc()
	}

	c.Stats.Influxdb.Add(Total, 1)
}
