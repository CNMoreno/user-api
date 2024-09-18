package constants

// Map errors users api.
var (
	ErrUserNotFound           = "user not found"
	ErrInvalidUserInput       = "Invalid user input"
	ErrFailedToCreateUser     = "Failed to create user"
	ErrFailedToGetUser        = "Failed to get user"
	ErrFailedToUpdateUser     = "Failed to update user"
	ErrFailedToDeleteUser     = "Failed to delete user"
	ErrFailedToGetFile        = "Failed to get file"
	ErrOnlyAcceptCSVFile      = "Only accept CSV file"
	ErrOpenFile               = "Failed to open file"
	ErrProcessCSVFile         = "Failed to process CSV file"
	ErrInsertUsers            = "Failed to create users"
	ErrClosingMongoConnection = "Failed closing MongoDB connection"
	ErrClosingFile            = "Failed closing File"
	ErrMongoUrlIsNotSet       = "MONGO_URL is not set"
	ErrMongoDatabaseIsNotSet  = "MONGO_DATABASE is not set"
	ErrCloseMongoConnection   = "Failed closing MongoDB connection"
	ErrSetUpDependencies      = "Failed to set up dependencies"
)
