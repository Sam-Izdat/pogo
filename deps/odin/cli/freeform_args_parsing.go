package cli

import "github.com/Sam-Izdat/pogo/deps/odin/cli/values"

func (cmd *CLI) assignUnparsedArgs(args []string) {
	for _, arg := range args {
		str := ""
		cmd.unparsedArgs = append(cmd.unparsedArgs, values.NewString(arg, &str))
	}
}
