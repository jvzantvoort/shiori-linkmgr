package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jvzantvoort/linkmgr/internal/config"
	"github.com/jvzantvoort/linkmgr/internal/display"
	"github.com/jvzantvoort/linkmgr/internal/repository"
	"github.com/spf13/cobra"
)

var testConnection bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long: `Initialize linkmgr configuration by creating a config file
with database connection settings. The configuration will be saved
to ~/.linkmgr.yaml with secure permissions (0600).`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&testConnection, "test-connection", false, "test existing database connection")
}

func runInit(cmd *cobra.Command, args []string) error {
	configPath := config.DefaultConfigPath()

	// If test-connection flag, just test and exit
	if testConnection {
		return testDatabaseConnection()
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("linkmgr Configuration Setup")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	cfg := &config.Config{}

	// Database configuration
	fmt.Println("Database Configuration:")
	fmt.Println()

	cfg.Database.Host = prompt(reader, "Database host", "localhost")
	cfg.Database.Port = promptInt(reader, "Database port", 3306)
	cfg.Database.User = prompt(reader, "Database user", "")
	cfg.Database.Password = promptPassword(reader, "Database password")
	cfg.Database.Database = prompt(reader, "Database name", "shiori")
	cfg.Database.MaxConnections = promptInt(reader, "Max connections", 10)

	// Display configuration
	fmt.Println()
	fmt.Println("Display Configuration:")
	fmt.Println()

	cfg.Display.DefaultLimit = promptInt(reader, "Default result limit", 50)
	cfg.Display.DateFormat = prompt(reader, "Date format", "2006-01-02 15:04")

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Test database connection
	fmt.Println()
	fmt.Println("Testing database connection...")

	dbConn, err := repository.NewDB(repository.Config{
		Host:           cfg.Database.Host,
		Port:           cfg.Database.Port,
		User:           cfg.Database.User,
		Password:       cfg.Database.Password,
		Database:       cfg.Database.Database,
		MaxConnections: cfg.Database.MaxConnections,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbConn.Close()

	ctx := context.Background()
	if err := dbConn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	display.Success("Database connection successful!")

	// Save configuration
	fmt.Println()
	fmt.Printf("Saving configuration to %s...\n", configPath)

	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	display.Success("Configuration saved successfully!")
	fmt.Println()
	fmt.Println("You can now use linkmgr to manage your bookmarks.")
	fmt.Println("Try: linkmgr list")

	return nil
}

func testDatabaseConnection() error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	fmt.Println("Testing database connection...")

	dbConn, err := repository.NewDB(repository.Config{
		Host:           cfg.Database.Host,
		Port:           cfg.Database.Port,
		User:           cfg.Database.User,
		Password:       cfg.Database.Password,
		Database:       cfg.Database.Database,
		MaxConnections: cfg.Database.MaxConnections,
	})
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer dbConn.Close()

	ctx := context.Background()
	if err := dbConn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	display.Success("Database connection successful!")
	return nil
}

func prompt(reader *bufio.Reader, label, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", label, defaultValue)
	} else {
		fmt.Printf("%s: ", label)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

func promptInt(reader *bufio.Reader, label string, defaultValue int) int {
	defaultStr := strconv.Itoa(defaultValue)
	input := prompt(reader, label, defaultStr)

	val, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}
	return val
}

func promptPassword(reader *bufio.Reader, label string) string {
	fmt.Printf("%s: ", label)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
