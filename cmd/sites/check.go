package sites

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

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks the status of a deployed site",
	RunE:  runCheck,
}

func runCheck(cmd *cobra.Command, _ []string) error {
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

	_, err = c.GetSite(ctx, spec.Sites.Name)
	if err != nil {
		return err
	}

	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := console.Render(tw, console.Error(fmt.Errorf("%s already exists", spec.Sites.Name)))
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
