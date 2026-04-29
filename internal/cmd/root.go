package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jvzantvoort/linkmgr/internal/config"
	"github.com/jvzantvoort/linkmgr/internal/repository"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
	cfg     *config.Config
	db      *repository.DB
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "linkmgr",
	Short: "Terminal-based bookmark manager",
	Long: `linkmgr is a fast, terminal-based bookmark management application.
It stores bookmarks in a MySQL database and provides CRUD operations
through an intuitive command-line interface.`,
	PersistentPreRunE: initializeConfig,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.linkmgr.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

// initializeConfig loads configuration and initializes database
func initializeConfig(cmd *cobra.Command, args []string) error {
	// Skip for init command
	if cmd.Name() == "init" || cmd.Name() == "version" {
		return nil
	}

	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w\nRun 'linkmgr init' to create a configuration file", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w\nRun 'linkmgr init' to reconfigure", err)
	}

	// Initialize database
	db, err = repository.NewDB(repository.Config{
		Host:           cfg.Database.Host,
		Port:           cfg.Database.Port,
		User:           cfg.Database.User,
		Password:       cfg.Database.Password,
		Database:       cfg.Database.Database,
		MaxConnections: cfg.Database.MaxConnections,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w\nCheck your configuration and database status", err)
	}

	// Test connection
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
