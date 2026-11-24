package owners

import (
	"fmt"
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
	"github.com/zeiss/builder/internal/adapters/client"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/console"

	"github.com/spf13/cobra"
)

type AddCmdFlags struct {
	// Email of the owners to add
	Email string `json:"email" yaml:"email" example:"user@example.com"`
}

var addCmdFlags = AddCmdFlags{}

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds an owner to a deployed site",
	RunE:  runAdd,
}

func runAdd(cmd *cobra.Command, _ []string) error {
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

	if err := c.AddOwner(ctx, spec.Sites.Name, addCmdFlags.Email); err != nil {
		return err
	}

	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := console.Render(tw, console.Success(fmt.Sprintf("Added %s successfully to %s", addCmdFlags.Email, spec.Sites.Name)))
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
