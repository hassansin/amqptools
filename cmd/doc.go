package cmd

import (
	"bytes"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func linkHandler(name string) string {
	base := strings.TrimSuffix(name, path.Ext(name))
	return "#" + strings.Replace(base, "_", "-", 5)
}

// consumeCmd represents the consume command
var docCmd = &cobra.Command{
	Use:    "doc",
	Short:  "Generate markdown document",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		rootDoc := new(bytes.Buffer)
		if err := doc.GenMarkdownCustom(RootCmd, rootDoc, linkHandler); err != nil {
			return err
		}

		consumeDoc := new(bytes.Buffer)
		if err := doc.GenMarkdownCustom(consumeCmd, consumeDoc, linkHandler); err != nil {
			return err
		}

		produceDoc := new(bytes.Buffer)
		if err := doc.GenMarkdownCustom(produceCmd, produceDoc, linkHandler); err != nil {
			return err
		}

		f, err := os.Create("README.md")
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err = f.Write(rootDoc.Bytes()); err != nil {
			return err
		}
		if _, err = f.Write(consumeDoc.Bytes()); err != nil {
			return err
		}
		if _, err = f.Write(produceDoc.Bytes()); err != nil {
			return err
		}
		f.Sync()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(docCmd)
}
