package howto

const DEFAULT_SYSTEM_MESSAGE = "You are a CLI tool that converts user requests to shell commands or short scripts. E.g., for `bash command to tar file without compression:`, you should reply `tar -cf file.tar file`. Avoid natural language. If you have to use it, be extremely concise. Less than 5 words."
const SERVICE_NAME = "howto"

const (
	OwnerReadWrite = 0600
	AllRead        = 0644
)

var examples = []string{
	"howto tar without compression",
	"howto oneline install conda",
	"howto du -hs hidden files",
	"howto donwload from gcp bucket",
	"howto pull from upstream",
	"howto push if the only update is the tag",
	"howto get ubuntu version",
	"howto undo make",
	"howto connect to mongo running inside docker",
	"howto check if something is running on my port 27017",
	"howto get user id for user vlialin",
	"howto create user vlialin with user IDs 5030 and GID 4030 and assign them a home directory in /mnt/shared_home",
	"howto tree withiout node_modules",
	"howto 'grep my zsh history and print all examples containing howto (with trailing space)'",
}
