package launcher

import "flag"

var (
	stopService       = flag.Bool("e", false, "Stop the service")
	showStatus        = flag.Bool("s", false, "Get status of service")
	outputServiceName = flag.Bool("n", false, "Output service's name")
	configFile        = flag.String("c", "config.json", "Config file")
	startService      = flag.String("b", "", "Start service,set manager's api address and output service's api address")
	installService    = flag.Bool("i", false, "Install service")
	removeService     = flag.Bool("r", false, "Remove service")
	rawRun            = flag.Bool("p", false, "Raw run service")
)
