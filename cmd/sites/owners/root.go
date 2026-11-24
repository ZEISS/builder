package owners

import (
	"fmt"
	"os"

	"github.com/zeiss/builder/internal/adapters/client"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/console"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(AddCmd)
	RootCmd.AddCommand(RemoveCmd)

	AddCmd.Flags().StringVar(&addCmdFlags.Email, "email", "", "Email of the owner to add")
	RemoveCmd.Flags().StringVar(&removeCmdFlags.Email, "email", "", "Email of the owner to remove")
}

var RootCmd = &cobra.Command{
	Use:   "owners",
	Short: "Lists the owners of a deployed site",
	RunE:  runRoot,
}

func runRoot(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	flags := config.DefaultFlags

	path, err := config.AppDirPath()
	if err != nil {
		return err
	}

	spec, err := flags.Spec()
	if err != nil {
		return err
	}

	if err := spec.Validate(); err != nil {
		return err
	}

	srv, err := services.SetupServices(ctx, path)
	if err != nil {
		return err
	}

	c, err := client.SetupClient(ctx, srv, flags.URL)
	if err != nil {
		return err
	}

	owners, err := c.ListOwners(ctx, spec.Sites.Name)
	if err != nil {
		return err
	}

	tw, _, _ := term.GetSize(os.Stdout.Fd())

	for _, owner := range owners {
		body := console.Render(tw, console.Success(owner.Email))
		fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)
	}

	return nil
}
