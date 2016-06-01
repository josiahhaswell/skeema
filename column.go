package tengo

import (
	"fmt"
	"strings"
)

type ColumnDefault struct {
	Null   bool
	Quoted bool
	Value  string
}

var ColumnDefaultNull = ColumnDefault{Null: true}
var ColumnDefaultCurrentTimestamp = ColumnDefault{Value: "CURRENT_TIMESTAMP"}

func ColumnDefaultValue(Value string) ColumnDefault {
	return ColumnDefault{
		Quoted: true,
		Value:  Value,
	}
}

func (cd ColumnDefault) Clause() string {
	if cd.Null {
		return "DEFAULT NULL"
	} else if cd.Quoted {
		// TODO: need to escape any quotes already here!!!
		return fmt.Sprintf("DEFAULT '%s'", cd.Value)
	} else {
		return fmt.Sprintf("DEFAULT %s", cd.Value)
	}
}

type Column struct {
	Name          string
	TypeInDB      string
	Nullable      bool
	AutoIncrement bool
	Default       ColumnDefault
	Extra         string
	//Comment       string
}

func (c Column) Definition() string {
	var notNull, autoIncrement, defaultValue, extraModifiers string
	emitDefault := c.CanHaveDefault()
	if !c.Nullable {
		notNull = " NOT NULL"
		if c.Default.Null {
			emitDefault = false
		}
	}
	if c.AutoIncrement {
		autoIncrement = " AUTO_INCREMENT"
	}
	if emitDefault {
		defaultValue = fmt.Sprintf(" %s", c.Default.Clause())
	}
	if c.Extra != "" {
		extraModifiers = fmt.Sprintf(" %s", c.Extra)
	}
	return fmt.Sprintf("%s %s%s%s%s%s", EscapeIdentifier(c.Name), c.TypeInDB, notNull, autoIncrement, defaultValue, extraModifiers)
}

func (c *Column) Equals(other *Column) bool {
	// shortcut if both nil pointers, or both pointing to same underlying struct
	if c == other {
		return true
	}
	// if one is nil, but we already know the two aren't equal, then we know the other is non-nil
	if c == nil || other == nil {
		return false
	}
	return (c.Name == other.Name &&
		c.TypeInDB == other.TypeInDB &&
		c.Nullable == other.Nullable &&
		c.AutoIncrement == other.AutoIncrement)
}

// Returns true if the column is allowed to have a DEFAULT clause
func (c Column) CanHaveDefault() bool {
	if c.AutoIncrement {
		return false
	}
	// MySQL does not permit defaults for these types
	if strings.HasSuffix(c.TypeInDB, "blob") || strings.HasSuffix(c.TypeInDB, "text") {
		return false
	}
	return true
}
