package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"ANJALI/config"
	"ANJALI/plugins/admins"
	"ANJALI/plugins/bot"
	"ANJALI/plugins/misc"
	"ANJALI/plugins/play"
	"ANJALI/plugins/sudo"
	"ANJALI/plugins/tools"
	"ANJALI/utils/database"

	tb "gopkg.in/tucnak/telebot.v2"
)

// ==========================================
// 🧹 CACHE SWEEPER
// ==========================================
func clearAllCaches() {
	Logger.Println("🧹 Sweeping old files for Main Bot and Clones...")
	dirs := []string{"downloads", "cache", "playback"}
	patterns := []string{"vid_*.mp4", "vid_*.m4a", "vid_*.webm", "*.webm", "*.mp4"}

	cleanedCount := 0

	// Clear floating matching files
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(pattern)
		for _, match := range matches {
			if err := os.Remove(match); err == nil {
				cleanedCount++
			}
		}
	}

	// Clear specified directories
	for _, dir := range dirs {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			files, _ := os.ReadDir(dir)
			for _, f := range files {
				path := filepath.Join(dir, f.Name())
				if f.IsDir() {
					os.RemoveAll(path)
				} else {
					if err := os.Remove(path); err == nil {
						cleanedCount++
					}
				}
			}
		}
	}

	if cleanedCount > 0 {
		Logger.Printf("✅ Successfully swept %d leftover temporary files from all bots.", cleanedCount)
	} else {
		Logger.Println("✅ Server storage is already clean.")
	}
}

func main() {
	// Initialize Logging
	InitLogger()

	// Clear Caches
	clearAllCaches()

	// Load Config
	cfg := config.LoadConfig()
	if cfg.BotToken == "" {
		Logger.Fatal("BOT_TOKEN is required. Add it to the deployment environment.")
	}

	// Connect to Database
	Logger.Println("Connecting to Database...")
	if err := database.Connect(cfg.MongoDBURI); err != nil {
		Logger.Printf("⚠️ MongoDB unavailable; continuing in in-memory mode: %v", err)
	}

	// Sync Sudoers from DB
	InitSudoers()

	Logger.Println("Starting Main Bot...")
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		Logger.Fatalf("Failed to start bot: %v", err)
	}

	Logger.Println("Registering Modules...")

	// Register Handlers
	bot.RegisterStartHandlers(b)
	bot.RegisterHelpHandlers(b)
	bot.RegisterSettingsHandlers(b)

	admins.RegisterAuthHandlers(b)
	admins.RegisterAutoplayHandlers(b)
	admins.RegisterAdminCallbacks(b)
	admins.RegisterStreamControls(b)

	misc.RegisterBroadcastHandlers(b)
	misc.RegisterCBroadcastHandlers(b)
	misc.StartSeekerLoop()
	misc.RegisterWatcherHandlers(b)

	play.RegisterChannelPlayHandlers(b)
	play.RegisterLiveStreamHandlers(b)
	play.RegisterPlayHandlers(b)
	play.RegisterPlaymodeHandlers(b)

	sudo.RegisterBlChatHandlers(b)
	sudo.RegisterBlUserHandlers(b)
	sudo.RegisterGbanHandlers(b)
	sudo.RegisterSystemSettingsHandlers(b)
	sudo.RegisterRestartHandlers(b)

	tools.RegisterActiveHandlers(b)
	tools.RegisterCloneHandlers(b)
	tools.RegisterDevHandlers(b)
	tools.RegisterLanguageHandlers(b)
	tools.RegisterPingHandlers(b)
	tools.RegisterReloadHandlers(b)
	tools.RegisterRepoUploadHandlers(b)
	tools.RegisterSecurityHandlers(b)
	tools.RegisterStatsHandlers(b)
	tools.RegisterUserIDHandlers(b)
	tools.RegisterVCLoggerHandlers(b)
	tools.RegisterWelcomeHandlers(b)

	Logger.Println("𝗔𝗹𝗹 𝗙𝗲𝗮𝘁𝘂𝗿𝗲𝘀 𝗟𝗼𝗮𝗱𝗲𝗱 𝗕𝗮𝗯𝘆🥳...")

	Logger.Println("╔═════ஜ۩۞۩ஜ════╗")
	Logger.Println("  ☠︎︎𝗠𝗔𝗗𝗘 𝗕𝗬 THE SHIV𝘀☠︎︎")
	Logger.Println("╚═════ஜ۩۞۩ஜ════╝")

	// Start Telebot
	go b.Start()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	Logger.Println("𝗦𝗧𝗢𝗣 𝗠𝗨𝗦𝗜𝗖🎻 𝗕𝗢𝗧..")
	b.Stop()
}
