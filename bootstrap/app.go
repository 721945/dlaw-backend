package bootstrap

//
//import (
//	"github.com/721945/dlaw-backend/cmd"
//	"github.com/spf13/cobra"
//)
//
////var rootCmd = &cobra.Command{
////	Use:   "dlaw-backend",
////	Short: "dlaw backend using gin framework",
////	Long: `
////This is a command runner or cli for api architecture in golang.
////Using this we can use underlying dependency injection container for running scripts.
////Main advantage is that, we can use same services, repositories, infrastructure present in the application itself`,
////	TraverseChildren: true,
////}
//
//type App struct {
//	*cobra.Command
//}
//
//func NewApp() App {
//	rootWrappedCmd := cmd.WrapCommand(CommonModules)
//
//	rootCmd := App{
//		Command: rootWrappedCmd,
//	}
//
//	rootCmd.AddCommand(cmd.MigrationCmd)
//
//	return rootCmd
//}
//
//var RootApp = NewApp()
