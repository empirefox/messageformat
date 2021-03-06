package messageformat

import (
	"bytes"
	"fmt"
	"strconv"
)

type Ordinal struct {
	Select
}

func newOrdinal(parent Node, varname string) *Ordinal {
	nd := &Ordinal{
		Select: Select{
			varname:    varname,
			ChoicesMap: make(map[string]Node, 5),
			choices:    make([]Choice, 0, 5),
		},
	}
	AddChild(parent, nd)
	return nd
}

func (nd *Ordinal) Type() string { return "selectordinal" }

// sort choices
func (s *Ordinal) Less(i, j int) bool {
	return pluralForms[s.choices[i].Key] < pluralForms[s.choices[j].Key]
}

func (nd *Ordinal) ToString(output *bytes.Buffer) error {
	return selectNodeToString(nd, output)
}

// It will returns an error if :
// - the associated value can't be convert to string or to an int (i.e. bool, ...)
// - the pluralFunc is not defined (MessageFormat.getNamedKey)
//
// It will falls back to the "other" choice if :
// - its key can't be found in the given map
// - the computed named key (MessageFormat.getNamedKey) is not a key of the given map
func (nd *Ordinal) Format(mf *MessageFormat, output *bytes.Buffer, data Data, _ string) error {
	key := nd.Varname()
	value, err := data.ValueStr(key)
	if err != nil {
		return err
	}

	var choice Node

	if v, ok := data[key]; ok {
		switch v.(type) {
		case int, float64:
		case string:
			_, err = strconv.ParseFloat(v.(string), 64)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Ordinal: Unsupported type for named key: %T", v)
		}

		key, err = mf.getNamedKey(v, true)
		if err != nil {
			return err
		}
		choice = nd.ChoicesMap[key]
	}

	if choice == nil {
		choice = nd.other
	}
	return choice.Format(mf, output, data, value)
}
