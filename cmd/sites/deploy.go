package sites

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
	"github.com/zeiss/builder/internal/adapters/client"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/console"
	"github.com/zeiss/builder/pkg/utils"

	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a site",
	RunE:  runDeploy,
}

func runDeploy(cmd *cobra.Command, _ []string) error {
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

	root, err := flags.ConfigPath()
	if err != nil {
		return err
	}

	scanDir := filepath.Join(filepath.Dir(root), spec.Sites.Path)
	files := utils.ScanDir(scanDir, spec.Sites.Ignore)

	tw, _, _ := term.GetSize(os.Stdout.Fd())

	for _, file := range files {
		file = strings.TrimPrefix(filepath.Clean(strings.TrimPrefix(file, scanDir)), "/")
		err := c.UploadFile(ctx, spec.Sites.Name, scanDir, file)
		if err != nil {
			return err
		}

		body := console.Render(tw, console.Success(fmt.Sprintf("%s file", file)))
		fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)
	}

	body := console.Render(tw, console.Success(fmt.Sprintf("%d files uploaded successfully.", len(files))))
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
