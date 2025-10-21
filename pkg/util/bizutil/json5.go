package bizutil

import (
	"encoding/json"

	"github.com/gookit/config/v2"
	"github.com/titanous/json5"
)

// JSON5Driver instance
var JSON5Driver = config.NewDriver("json5", json5.Unmarshal, json.Marshal)
