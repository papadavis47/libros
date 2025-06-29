package database

// Note: Backup functionality has been moved to the services/backup.go package
// for better separation of concerns. The Database struct now focuses solely
// on data persistence operations.
//
// The following methods have been moved:
// - BackupToJSON -> services.BackupService.ExportToJSON
// - BackupToMarkdown -> services.BackupService.ExportToMarkdown
// - Database file backup -> services.BackupService.BackupDatabase
//
// This refactoring improves code organization by separating backup logic
// from core database operations, making the codebase more modular and testable.