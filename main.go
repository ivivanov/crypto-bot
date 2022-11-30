package main

import "github.com/ivivanov/crypto-bot/cmd"

// func main() {
// 	app, err := NewApp()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", err)
// 		os.Exit(1)
// 	}

// 	if err := app.Run(); err != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", err)
// 		os.Exit(1)
// 	}
// }

func main() {
	cmd.Execute()
}
