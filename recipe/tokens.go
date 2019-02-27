package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

const (
	KEYWORD_VAR     = "var"
	KEYWORD_COMMAND = "command"

	OPTION_UNSAFE_ACTIONS = "unsafe-actions"
	OPTION_REQUIRE_ROOT   = "require-root"
	OPTION_FAST_FINISH    = "fast-finish"
	OPTION_LOCK_WORKDIR   = "lock-workdir"

	ACTION_EXIT            = "exit"
	ACTION_WAIT            = "wait"
	ACTION_EXPECT          = "expect"
	ACTION_WAIT_OUTPUT     = "wait_output"
	ACTION_OUTPUT_MATCH    = "output_match"
	ACTION_OUTPUT_CONTAINS = "output_contains"
	ACTION_OUTPUT_TRIM     = "output_trim"
	ACTION_PRINT           = "print"
	ACTION_CHDIR           = "chdir"
	ACTION_PERMS           = "perms"
	ACTION_OWNER           = "owner"
	ACTION_EXIST           = "exist"
	ACTION_READABLE        = "readable"
	ACTION_WRITABLE        = "writable"
	ACTION_EXECUTABLE      = "executable"
	ACTION_DIR             = "dir"
	ACTION_EMPTY           = "empty"
	ACTION_EMPTY_DIR       = "empty_dir"
	ACTION_CHECKSUM        = "checksum"
	ACTION_CHECKSUM_READ   = "checksum_read"
	ACTION_FILE_CONTAINS   = "file_contains"
	ACTION_COPY            = "copy"
	ACTION_MOVE            = "move"
	ACTION_TOUCH           = "touch"
	ACTION_MKDIR           = "mkdir"
	ACTION_REMOVE          = "remove"
	ACTION_CHMOD           = "chmod"
	ACTION_BACKUP          = "backup"
	ACTION_BACKUP_RESTORE  = "backup_restore"
	ACTION_PROCESS_WORKS   = "process_works"
	ACTION_WAIT_PID        = "wait_pid"
	ACTION_WAIT_FS         = "wait_fs"
	ACTION_CONNECT         = "connect"
	ACTION_APP             = "app"
	ACTION_ENV             = "env"
	ACTION_USER_EXIST      = "user_exist"
	ACTION_USER_ID         = "user_id"
	ACTION_USER_GID        = "user_gid"
	ACTION_USER_GROUP      = "user_group"
	ACTION_USER_SHELL      = "user_shell"
	ACTION_USER_HOME       = "user_home"
	ACTION_GROUP_EXIST     = "group_exist"
	ACTION_GROUP_ID        = "group_id"
	ACTION_SERVICE_PRESENT = "service_present"
	ACTION_SERVICE_ENABLED = "service_enabled"
	ACTION_SERVICE_WORKS   = "service_works"
	ACTION_HTTP_STATUS     = "http_status"
	ACTION_HTTP_HEADER     = "http_header"
	ACTION_HTTP_CONTAINS   = "http_contains"
	ACTION_LIB_LOADED      = "lib_loaded"
	ACTION_LIB_HEADER      = "lib_header"
	ACTION_LIB_CONFIG      = "lib_config"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// TokenInfo contains info about supported token
type TokenInfo struct {
	Keyword       string
	MinArgs       int
	MaxArgs       int
	Global        bool
	AllowNegative bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Tokens is slice with tokens info
var Tokens = []TokenInfo{
	{KEYWORD_VAR, 2, 2, true, false},
	{KEYWORD_COMMAND, 1, 2, true, false},

	{OPTION_UNSAFE_ACTIONS, 1, 1, true, false},
	{OPTION_REQUIRE_ROOT, 1, 1, true, false},
	{OPTION_FAST_FINISH, 1, 1, true, false},
	{OPTION_LOCK_WORKDIR, 1, 1, true, false},

	{ACTION_EXIT, 1, 2, false, true},
	{ACTION_WAIT, 1, 1, false, false},

	{ACTION_EXPECT, 1, 2, false, false},
	{ACTION_WAIT_OUTPUT, 1, 1, false, false},
	{ACTION_OUTPUT_MATCH, 1, 1, false, true},
	{ACTION_OUTPUT_CONTAINS, 1, 1, false, true},
	{ACTION_OUTPUT_TRIM, 0, 0, false, false},
	{ACTION_PRINT, 1, 1, false, false},

	{ACTION_CHDIR, 1, 1, false, false},
	{ACTION_PERMS, 2, 2, false, true},
	{ACTION_OWNER, 2, 2, false, true},
	{ACTION_EXIST, 1, 1, false, true},
	{ACTION_READABLE, 2, 2, false, true},
	{ACTION_WRITABLE, 2, 2, false, true},
	{ACTION_EXECUTABLE, 2, 2, false, true},
	{ACTION_DIR, 1, 1, false, true},
	{ACTION_EMPTY, 1, 1, false, true},
	{ACTION_EMPTY_DIR, 1, 1, false, true},

	{ACTION_CHECKSUM, 2, 2, false, true},
	{ACTION_CHECKSUM_READ, 2, 2, false, false},
	{ACTION_FILE_CONTAINS, 2, 2, false, true},

	{ACTION_COPY, 2, 2, false, false},
	{ACTION_MOVE, 2, 2, false, false},
	{ACTION_TOUCH, 1, 1, false, false},
	{ACTION_MKDIR, 1, 1, false, false},
	{ACTION_REMOVE, 1, 1, false, false},
	{ACTION_CHMOD, 2, 2, false, false},

	{ACTION_BACKUP, 1, 1, false, false},
	{ACTION_BACKUP_RESTORE, 1, 1, false, false},

	{ACTION_PROCESS_WORKS, 1, 1, false, true},
	{ACTION_WAIT_PID, 1, 2, false, true},
	{ACTION_WAIT_FS, 1, 2, false, true},
	{ACTION_CONNECT, 2, 2, false, true},
	{ACTION_APP, 1, 1, false, true},
	{ACTION_ENV, 2, 2, false, true},

	{ACTION_USER_EXIST, 1, 1, false, true},
	{ACTION_USER_ID, 2, 2, false, true},
	{ACTION_USER_GID, 2, 2, false, true},
	{ACTION_USER_GROUP, 2, 2, false, true},
	{ACTION_USER_SHELL, 2, 2, false, true},
	{ACTION_USER_HOME, 2, 2, false, true},
	{ACTION_GROUP_EXIST, 1, 1, false, true},
	{ACTION_GROUP_ID, 2, 2, false, true},

	{ACTION_SERVICE_PRESENT, 1, 1, false, true},
	{ACTION_SERVICE_ENABLED, 1, 1, false, true},
	{ACTION_SERVICE_WORKS, 1, 1, false, true},

	{ACTION_HTTP_STATUS, 3, 3, false, true},
	{ACTION_HTTP_HEADER, 4, 4, false, true},
	{ACTION_HTTP_CONTAINS, 3, 3, false, true},

	{ACTION_LIB_LOADED, 1, 1, false, true},
	{ACTION_LIB_HEADER, 1, 1, false, true},
	{ACTION_LIB_CONFIG, 1, 1, false, true},
}
