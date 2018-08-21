/**
 * @package hook
 * @file flags.go
 * @author sufay
 *
 * hook flags for cli
 */

package hook

import (
	"gopkg.in/urfave/cli.v1"
)

// hook enabled flag
var HookEnabledFlag = cli.BoolFlag{
	Name:  "hook",
	Usage: "Enable hook",
}
