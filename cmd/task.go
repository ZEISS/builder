package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/katallaxie/pkg/conv"
	"github.com/katallaxie/pkg/filex"
	"github.com/katallaxie/pkg/mapx"
	"github.com/katallaxie/pkg/slices"
	"github.com/spf13/cobra"
)

var ErrTaskNotExist = fmt.Errorf("task does not exist")

func init() {
	TaskCmd.AddCommand(TaskNewCmd)
	TaskCmd.AddCommand(TaskWorkCmd)
	TaskCmd.PersistentFlags().StringVar(&cfg.Flags.TaskFlags.Name, "name", "", "Name of the task")
}

var TaskCmd = &cobra.Command{
	Use: "task",
}

var TaskNewCmd = &cobra.Command{
	Use:   "new name",
	Short: "Create a new task [name]",
	RunE:  createNewTask,
}

func createNewTask(_ *cobra.Command, args []string) error {
	err := cfg.LoadSpec()
	if err != nil {
		return err
	}

	name := slices.First(args...)

	ok := mapx.Exists(cfg.Spec.Tasks, name)
	if !ok {
		return ErrTaskNotExist
	}

	cwd, err := cfg.Cwd()
	if err != nil {
		return err
	}

	path := filepath.Join(cwd, cfg.Spec.Root, name, cfg.Flags.TaskFlags.Name)

	err = filex.MkdirAll(path, 0o755) //nolint:mnd
	if err != nil {
		return err
	}

	task := cfg.Spec.Tasks[name]

	for _, doc := range task.Context.Documents {
		fp := filepath.Join(path, doc.Generates)
		err := os.WriteFile(fp, conv.Bytes(doc.Template), 0o644) //nolint:mnd
		if err != nil {
			return err
		}
	}

	return nil
}

var TaskWorkCmd = &cobra.Command{
	Use:   "work name",
	Short: "Work on a task [name]",
	RunE:  workOnTask,
}

func workOnTask(_ *cobra.Command, _ []string) error {
	err := cfg.LoadSpec()
	if err != nil {
		return err
	}

	// name := slices.First(args...)

	// ok := mapx.Exists(cfg.Spec.Tasks, name)
	// if !ok {
	// 	return ErrTaskNotExist
	// }

	// // cwd, err := cfg.Cwd()
	// // if err != nil {
	// // 	return err
	// // }

	// // path := filepath.Join(cwd, cfg.Spec.Root, name)

	// app, err := app.New(cmd.Context(), cfg)
	// if err != nil {
	// 	return err
	// }

	// defer app.Dispose()

	// zone.NewGlobal()
	// program := tea.NewProgram(
	// 	ui.New(app),
	// 	tea.WithAltScreen(),
	// 	tea.WithReportFocus(),
	// )

	// // clear the terminal
	// os.Stdout.WriteString("\x1b[2J\x1b[3J\x1b[H")

	// _, err = program.Run()
	// if err != nil {
	// 	return err
	// }

	return nil
}
