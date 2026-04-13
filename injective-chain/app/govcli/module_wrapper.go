package govcli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/x/gov"
)

// GovModuleWrapper embeds gov.AppModule and provides a custom query command
// (proposal and proposals) that uses the app codec for JSON output, avoiding
// aminojson reflection panics on proposals containing messages with LegacyDec.
// Autocli adds the remaining gov query subcommands via EnhanceCustomCommand.
type GovModuleWrapper struct {
	gov.AppModule
}

// GetQueryCmd returns the custom gov query command so autocli uses it
// instead of generating one for the gov module.
func (GovModuleWrapper) GetQueryCmd() *cobra.Command {
	return GetQueryCmd()
}
