package po

import spec "github.com/Sam-Izdat/pogo/gtspec"

var o spec.Config

func loadOptions() {
	var err error
	o, err = spec.LoadOptions()
	if err != nil {
		return // defer error reporting to cli
	}
}
