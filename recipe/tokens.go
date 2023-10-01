package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

const (
	SYMBOL_COMMAND_GROUP   = "+"
	SYMBOL_NEGATIVE_ACTION = "!"
	SYMBOL_SEPARATOR       = ":"

	KEYWORD_VAR     = "var"
	KEYWORD_COMMAND = "command"
	KEYWORD_PACKAGE = "pkg"

	OPTION_UNSAFE_ACTIONS    = "unsafe-actions"
	OPTION_REQUIRE_ROOT      = "require-root"
	OPTION_FAST_FINISH       = "fast-finish"
	OPTION_LOCK_WORKDIR      = "lock-workdir"
	OPTION_UNBUFFER          = "unbuffer"
	OPTION_HTTPS_SKIP_VERIFY = "https-skip-verify"
	OPTION_DELAY             = "delay"

	ACTION_EXIT = "exit"
	ACTION_WAIT = "wait"

	ACTION_EXPECT          = "expect"
	ACTION_WAIT_OUTPUT     = "wait-output"
	ACTION_OUTPUT_MATCH    = "output-match"
	ACTION_OUTPUT_CONTAINS = "output-contains"
	ACTION_OUTPUT_EMPTY    = "output-empty"
	ACTION_OUTPUT_TRIM     = "output-trim"
	ACTION_PRINT           = "print"

	ACTION_CHDIR      = "chdir"
	ACTION_MODE       = "mode"
	ACTION_OWNER      = "owner"
	ACTION_EXIST      = "exist"
	ACTION_LINK       = "link"
	ACTION_READABLE   = "readable"
	ACTION_WRITABLE   = "writable"
	ACTION_EXECUTABLE = "executable"
	ACTION_DIR        = "dir"
	ACTION_EMPTY      = "empty"
	ACTION_EMPTY_DIR  = "empty-dir"

	ACTION_CHECKSUM      = "checksum"
	ACTION_CHECKSUM_READ = "checksum-read"
	ACTION_FILE_CONTAINS = "file-contains"

	ACTION_COPY     = "copy"
	ACTION_MOVE     = "move"
	ACTION_TOUCH    = "touch"
	ACTION_MKDIR    = "mkdir"
	ACTION_REMOVE   = "remove"
	ACTION_CHMOD    = "chmod"
	ACTION_CHOWN    = "chown"
	ACTION_TRUNCATE = "truncate"
	ACTION_CLEANUP  = "cleanup"

	ACTION_BACKUP         = "backup"
	ACTION_BACKUP_RESTORE = "backup-restore"

	ACTION_PROCESS_WORKS = "process-works"
	ACTION_WAIT_PID      = "wait-pid"
	ACTION_WAIT_FS       = "wait-fs"
	ACTION_WAIT_CONNECT  = "wait-connect"
	ACTION_CONNECT       = "connect"
	ACTION_APP           = "app"
	ACTION_SIGNAL        = "signal"
	ACTION_ENV           = "env"
	ACTION_ENV_SET       = "env-set"

	ACTION_USER_EXIST  = "user-exist"
	ACTION_USER_ID     = "user-id"
	ACTION_USER_GID    = "user-gid"
	ACTION_USER_GROUP  = "user-group"
	ACTION_USER_SHELL  = "user-shell"
	ACTION_USER_HOME   = "user-home"
	ACTION_GROUP_EXIST = "group-exist"
	ACTION_GROUP_ID    = "group-id"

	ACTION_SERVICE_PRESENT = "service-present"
	ACTION_SERVICE_ENABLED = "service-enabled"
	ACTION_SERVICE_WORKS   = "service-works"
	ACTION_WAIT_SERVICE    = "wait-service"

	ACTION_HTTP_STATUS     = "http-status"
	ACTION_HTTP_HEADER     = "http-header"
	ACTION_HTTP_CONTAINS   = "http-contains"
	ACTION_HTTP_JSON       = "http-json"
	ACTION_HTTP_SET_AUTH   = "http-set-auth"
	ACTION_HTTP_SET_HEADER = "http-set-header"

	ACTION_LIB_LOADED   = "lib-loaded"
	ACTION_LIB_HEADER   = "lib-header"
	ACTION_LIB_CONFIG   = "lib-config"
	ACTION_LIB_EXIST    = "lib-exist"
	ACTION_LIB_LINKED   = "lib-linked"
	ACTION_LIB_RPATH    = "lib-rpath"
	ACTION_LIB_SONAME   = "lib-soname"
	ACTION_LIB_EXPORTED = "lib-exported"

	ACTION_PYTHON2_PACKAGE = "python2-package"
	ACTION_PYTHON3_PACKAGE = "python3-package"
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
	{KEYWORD_PACKAGE, 1, 999, true, false},

	{OPTION_UNSAFE_ACTIONS, 1, 1, true, false},
	{OPTION_REQUIRE_ROOT, 1, 1, true, false},
	{OPTION_FAST_FINISH, 1, 1, true, false},
	{OPTION_LOCK_WORKDIR, 1, 1, true, false},
	{OPTION_UNBUFFER, 1, 1, true, false},
	{OPTION_HTTPS_SKIP_VERIFY, 1, 1, true, false},
	{OPTION_DELAY, 1, 1, true, false},

	{ACTION_EXIT, 1, 2, false, true},
	{ACTION_WAIT, 1, 1, false, false},

	{ACTION_EXPECT, 1, 2, false, false},
	{ACTION_WAIT_OUTPUT, 1, 1, false, false},
	{ACTION_OUTPUT_MATCH, 1, 1, false, true},
	{ACTION_OUTPUT_CONTAINS, 1, 1, false, true},
	{ACTION_OUTPUT_EMPTY, 0, 0, false, true},
	{ACTION_OUTPUT_TRIM, 0, 0, false, false},
	{ACTION_PRINT, 1, 1, false, false},

	{ACTION_CHDIR, 1, 1, false, false},
	{ACTION_MODE, 2, 2, false, true},
	{ACTION_OWNER, 2, 2, false, true},
	{ACTION_EXIST, 1, 1, false, true},
	{ACTION_LINK, 2, 2, false, true},
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
	{ACTION_CHOWN, 2, 2, false, false},
	{ACTION_TRUNCATE, 1, 1, false, false},
	{ACTION_CLEANUP, 1, 1, false, false},

	{ACTION_BACKUP, 1, 1, false, false},
	{ACTION_BACKUP_RESTORE, 1, 1, false, false},

	{ACTION_PROCESS_WORKS, 1, 1, false, true},
	{ACTION_WAIT_PID, 1, 2, false, true},
	{ACTION_WAIT_FS, 1, 2, false, true},
	{ACTION_WAIT_CONNECT, 2, 3, false, true},
	{ACTION_CONNECT, 2, 3, false, true},
	{ACTION_APP, 1, 1, false, true},
	{ACTION_SIGNAL, 1, 2, false, false},
	{ACTION_ENV, 2, 2, false, true},
	{ACTION_ENV_SET, 2, 2, false, false},

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
	{ACTION_WAIT_SERVICE, 1, 2, false, true},

	{ACTION_HTTP_STATUS, 3, 4, false, true},
	{ACTION_HTTP_HEADER, 4, 5, false, true},
	{ACTION_HTTP_CONTAINS, 3, 4, false, true},
	{ACTION_HTTP_JSON, 4, 4, false, true},
	{ACTION_HTTP_SET_AUTH, 2, 2, false, false},
	{ACTION_HTTP_SET_HEADER, 2, 2, false, false},

	{ACTION_LIB_LOADED, 1, 1, false, true},
	{ACTION_LIB_HEADER, 1, 1, false, true},
	{ACTION_LIB_CONFIG, 1, 1, false, true},
	{ACTION_LIB_EXIST, 1, 1, false, true},
	{ACTION_LIB_LINKED, 2, 2, false, true},
	{ACTION_LIB_RPATH, 2, 2, false, true},
	{ACTION_LIB_SONAME, 2, 2, false, true},
	{ACTION_LIB_EXPORTED, 2, 2, false, true},

	{ACTION_PYTHON2_PACKAGE, 1, 1, false, false},
	{ACTION_PYTHON3_PACKAGE, 1, 1, false, false},
}
