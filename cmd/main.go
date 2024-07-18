package main

import (
	"fmt"
	"os"

	photoLayouts "github.com/leptobo/photo-layouts"
	"github.com/spf13/cobra"
)

var photo photoLayouts.Photo

var rootCmd = &cobra.Command{
	Use:   "pl",
	Short: "PL is a simple photo layout program",
	Long:  `Complete documentation is available at https://github.com/leptobo/photo-layouts`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := photoLayouts.Layout(&photo)
		if err != nil {
			cmd.PrintErrf("error: %v\n", err)
			os.Exit(1)
		}

		cmd.Printf("success, save to: %s\n", file)
	},
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&photo.File, "file", "f", "", "photo file")
	rootCmd.PersistentFlags().Float64VarP(&photo.ContainerWidth, "ppwidth", "W", 0, "photo paper width (mm)")
	rootCmd.PersistentFlags().Float64VarP(&photo.ContainerHeight, "ppheight", "H", 0, "photo paper height (mm)")
	rootCmd.PersistentFlags().Float64Var(&photo.PhotoWidth, "pw", 0, "photo width (mm)")
	rootCmd.PersistentFlags().Float64Var(&photo.PhotoHeight, "ph", 0, "photo height (mm)")
	rootCmd.PersistentFlags().StringVarP(&photo.Color, "color", "c", "", "set background color (hex)")
	rootCmd.PersistentFlags().Float64VarP(&photo.Dpi, "dpi", "d", 0, "dpi")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
