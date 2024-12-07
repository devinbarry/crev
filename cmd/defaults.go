package cmd

// specificPrefixesToIgnore contains file/directory prefixes that should be ignored by default
var specificPrefixesToIgnore = []string{
	// Version control and IDE directories
	".",    // Covers .git, .idea, .vscode, etc.
	"crev", // ignore crev specific files
}

// specificExtensionsToIgnore contains file extensions that should be ignored by default
var specificExtensionsToIgnore = []string{
	// Images
	".jpeg",
	".jpg",
	".png",
	".gif",
	".svg",
	".ico",

	// Documents
	".pdf",

	// Fonts
	".woff",
	".woff2",
	".eot",
	".ttf",
	".otf",
}

// specificFilesToIgnore contains specific filenames that should be ignored by default
var specificFilesToIgnore = []string{
	"Thumbs.db",   // Windows thumbnail cache
	"poetry.lock", // Python poetry lock file
	"go.mod",      // Go module file
	"go.sum",      // Go module checksum file
}
