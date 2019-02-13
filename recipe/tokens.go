package recipe

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

// TokenInfo contains info about supported token
type TokenInfo struct {
	Keyword       string
	MinArgs       int
	MaxArgs       int
	Global        bool
	AllowNegative bool
}

// Tokens is slice with tokens info
var Tokens = []TokenInfo{
	{"var", 2, 2, true, false},
	{"command", 1, 2, true, false},

	{"dir", 1, 1, true, false},
	{"unsafe-actions", 1, 1, true, false},
	{"require-root", 1, 1, true, false},
	{"fast-finish", 1, 1, true, false},
	{"lock-workdir", 1, 1, true, false},

	{"exit", 1, 2, false, true},
	{"wait", 1, 1, false, false},

	{"expect", 1, 2, false, false},
	{"wait-output", 1, 1, false, false},
	{"output-match", 1, 1, false, true},
	{"output-prefix", 1, 1, false, true},
	{"output-suffix", 1, 1, false, true},
	{"output-length", 1, 1, false, true},
	{"output-contains", 1, 1, false, true},
	{"output-equal", 1, 1, false, true},
	{"output-trim", 0, 0, false, false},
	{"print", 1, 1, false, false},

	{"chdir", 1, 1, false, false},
	{"perms", 2, 2, false, true},
	{"owner", 2, 2, false, true},
	{"exist", 1, 1, false, true},
	{"readable", 2, 2, false, true},
	{"writable", 2, 2, false, true},
	{"executable", 2, 2, false, true},
	{"directory", 1, 1, false, true},
	{"empty", 1, 1, false, true},
	{"empty-directory", 1, 1, false, true},

	{"checksum", 2, 2, false, true},
	{"checksum-read", 2, 2, false, false},
	{"file-contains", 2, 2, false, true},

	{"copy", 2, 2, false, false},
	{"move", 2, 2, false, false},
	{"touch", 1, 1, false, false},
	{"mkdir", 1, 1, false, false},
	{"remove", 1, 1, false, false},
	{"chmod", 2, 2, false, false},

	{"process-works", 1, 1, false, true},
	{"wait-pid", 1, 2, false, true},
	{"wait-fs", 1, 2, false, true},
	{"connect", 2, 2, false, true},
	{"app", 1, 1, false, true},
	{"env", 2, 2, false, true},

	{"user-exist", 1, 1, false, true},
	{"user-id", 2, 2, false, true},
	{"user-gid", 2, 2, false, true},
	{"user-group", 2, 2, false, true},
	{"user-shell", 2, 2, false, true},
	{"user-home", 2, 2, false, true},
	{"group-exist", 1, 1, false, true},
	{"group-id", 2, 2, false, true},

	{"service-present", 1, 1, false, true},
	{"service-enabled", 1, 1, false, true},
	{"service-works", 1, 1, false, true},

	{"http-status", 3, 3, false, true},
	{"http-header", 4, 4, false, true},
	{"http-contains", 3, 3, false, true},

	{"lib-loaded", 1, 1, false, true},
	{"lib-header", 1, 1, false, true},
}
