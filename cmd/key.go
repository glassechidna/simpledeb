package cmd

import (
	"github.com/aidansteele/simpledeb/pkg/simpledeb"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func init() {
	var keyCmd = &cobra.Command{
		Use:   "key",
		Short: "Generate a GPG key that can be used for signing apt repositories",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.PersistentFlags().GetString("name")
			email, _ := cmd.PersistentFlags().GetString("email")
			pub, _ := cmd.PersistentFlags().GetString("public-path")
			priv, _ := cmd.PersistentFlags().GetString("private-path")

			pub, _ = homedir.Expand(pub)
			priv, _ = homedir.Expand(priv)
			key(name, email, pub, priv)
		},
	}

	keyCmd.PersistentFlags().String("name", "", "Your name - appears in key")
	keyCmd.PersistentFlags().String("email", "", "Your email - appears in key")
	keyCmd.PersistentFlags().String("public-path", "signer.pub", "")
	keyCmd.PersistentFlags().String("private-path", "signer.key", "")
	RootCmd.AddCommand(keyCmd)
}

func key(name, email, pub, priv string) {
	keys := simpledeb.CreateKey(name, email)
	ioutil.WriteFile(pub, keys.PublicKey, 0644)
	ioutil.WriteFile(priv, keys.PrivateKey, 0600)
}
