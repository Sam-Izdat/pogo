package cli

import (
	"bytes"
	"fmt"
	"strings"
)

// DefaultUsage is the Default Usage used by odin
func (cmd *CLI) DefaultUsage() {
	fmt.Fprintln(cmd.StdOutput(), cmd.UsageString())
}

// FlagsUsageString returns the flags usage as a string
func (cmd *CLI) FlagsUsageString() string {
	var maxBufferLen int
	flagsUsages := make(map[*Flag]*bytes.Buffer)

	// init the map for each flag
	for _, flag := range cmd.aliases {
		flagsUsages[flag] = bytes.NewBufferString("")
	}

	// Get each flags aliases
	for r, flag := range cmd.aliases {
		alias := string(r)
		buffer := flagsUsages[flag]
		var err error
		if buffer.Len() == 0 {
			_, err = buffer.WriteString(fmt.Sprintf("-%s", alias))
		} else {
			_, err = buffer.WriteString(fmt.Sprintf(", -%s", alias))
		}
		exitIfError(err)
		buffLen := len(buffer.String())
		if buffLen > maxBufferLen {
			maxBufferLen = buffLen
		}
	}

	// Get each flags names
	for name, flag := range cmd.flags {
		buffer := flagsUsages[flag]
		if buffer == nil {
			flagsUsages[flag] = new(bytes.Buffer)
			buffer = flagsUsages[flag]
		}
		var err error
		if buffer.Len() == 0 {
			_, err = buffer.WriteString(fmt.Sprintf("--%s", name))
		} else {
			_, err = buffer.WriteString(fmt.Sprintf(", --%s", name))
		}
		if _, ok := flag.value.(boolFlag); !ok {
			buffer.WriteString(fmt.Sprintf("=\"%s\"", flag.DefValue))
		}
		exitIfError(err)
		buffLen := len(buffer.String())
		if buffLen > maxBufferLen {
			maxBufferLen = buffLen
		}
	}

	// get the flag strings and append the usage info
	var outputLines []string
	for i := 0; i < len(cmd.flags); i++ {
		flag := cmd.flags.Sort()[i]
		buffer := flagsUsages[flag]
		for {
			buffLen := len(buffer.String())
			if buffLen > maxBufferLen {
				break
			}
			buffer.WriteString(" ")
		}
		outputLines = append(outputLines, fmt.Sprintf("  %s # %s", buffer.String(), flag.Usage))
	}

	return strings.Join(outputLines, "\n")
}

// ParamsUsageString returns the params usage as a string
func (cmd *CLI) ParamsUsageString() string {
	var formattednames []string
	for i := 0; i < len(cmd.params); i++ {
		param := cmd.params[i]
		formattednames = append(formattednames, fmt.Sprintf("<%s>", param.Name))
	}
	return strings.Join(formattednames, " ")
}

// SubCommandsUsageString is the usage string for sub commands
func (cmd *CLI) SubCommandsUsageString() string {
	var maxBufferLen int
	for _, cmd := range cmd.subCommands {
		buffLen := len(cmd.name)
		if buffLen > maxBufferLen {
			maxBufferLen = buffLen
		}
	}

	var outputLines []string
	for _, cmd := range cmd.subCommands {
		var whitespace bytes.Buffer
		for {
			buffLen := len(cmd.name) + len(whitespace.String())
			if buffLen == maxBufferLen+5 {
				break
			}
			whitespace.WriteString(" ")
		}
		outputLines = append(outputLines, fmt.Sprintf("  %s%s%s", cmd.name, whitespace.String(), cmd.Description()))
	}

	return strings.Join(outputLines, "\n")
}

// Usage calls the Usage method for the flag set
func (cmd *CLI) Usage() {
	if cmd.usage == nil {
		cmd.usage = cmd.DefaultUsage
	}
	cmd.usage()
}

// UsageString returns the command usage as a string
func (cmd *CLI) UsageString() string {
	hasSubCommands := len(cmd.subCommands) > 0
	hasParams := len(cmd.params) > 0
	hasDescription := len(cmd.description) > 0

	// Start the Buffer
	var buff bytes.Buffer

	buff.WriteString("Usage:\n")
	buff.WriteString(fmt.Sprintf("  %s [options...]", cmd.Name()))

	// Write Param Syntax
	if hasParams {
		buff.WriteString(fmt.Sprintf(" %s", cmd.ParamsUsageString()))
	}

	// Write Sub Command Syntax
	if hasSubCommands {
		buff.WriteString(" <command> [arg...]")
	}

	if hasDescription {
		buff.WriteString(fmt.Sprintf("\n\n%s", cmd.Description()))
	}

	// Write Flags Syntax
	buff.WriteString("\n\nOptions:\n")
	buff.WriteString(cmd.FlagsUsageString())

	// Write Sub Command List
	if hasSubCommands {
		buff.WriteString("\n\nCommands:\n")
		buff.WriteString(cmd.SubCommandsUsageString())
	}

	// Return buffer as string
	return buff.String()
}
