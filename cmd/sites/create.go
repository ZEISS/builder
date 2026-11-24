package sites

import (
	"fmt"
	"os"

	"github.com/zeiss/builder/internal/adapters/client"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/console"
	"github.com/zeiss/builder/pkg/sbac"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new site",
	RunE:  runCreate,
}

func runCreate(cmd *cobra.Command, _ []string) error {
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

	site := &sbac.Site{Name: spec.Sites.Name}
	if err := c.Create(ctx, site); err != nil {
		return err
	}

	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := console.Render(tw, console.Success(fmt.Sprintf("%s created successfully", spec.Sites.Name)))
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
