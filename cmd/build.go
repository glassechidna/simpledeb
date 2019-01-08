package cmd

import (
	"fmt"
	"github.com/esell/deb-simple/pkg/debsimple"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build an apt repo from a collection of deb files",
		Long: `
Build an apt repository by specifying a collection of deb files (as positional
parameters), an output directory and a GPG signing key to sign the Releases
index.
`,
		Run: func(cmd *cobra.Command, args []string) {
			key, _ := cmd.PersistentFlags().GetString("key")
			key, _ = homedir.Expand(key)

			output, _ := cmd.PersistentFlags().GetString("output")
			output, _ = homedir.Expand(output)

			debs := []string{}
			for _, arg := range args {
				path, _ := homedir.Expand(arg)
				debs = append(debs, path)
			}

			build(key, output, "stable", "main", debs)
		},
	}

	buildCmd.PersistentFlags().StringP("key", "k", "signer.key", "Path to GPG key to sign repo")
	buildCmd.PersistentFlags().StringP("output", "o", "repo", "Path to GPG key to sign repo")
	RootCmd.AddCommand(buildCmd)
}

func build(key, output, distro, section string, debs []string) {
	c := debsimple.Conf{
		RootRepoPath:  output,
		SupportArch:   []string{"i386", "amd64"},
		DistroNames:   []string{distro},
		Sections:      []string{section},
		PrivateKey:    key,
		EnableSigning: true,
	}

	debsimple.Main(c)

	for _, deb := range debs {
		arch := guessArchitecture(deb)
		newpath := debsimple.CopyDeb(deb, c, distro, section, arch)
		debsimple.CreateMetadata(newpath, c)
		fmt.Fprintf(os.Stderr, "Processed %s\n", deb)
	}
}

func guessArchitecture(name string) string {
	if strings.Contains(name, "amd64") {
		return "amd64"
	} else {
		return "i386"
	}
}
