package desktop

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/un-def/manygram/internal/util"
)

var entryTemplate = template.Must(template.New("Desktop Entry").Parse(`[Desktop Entry]
Version=1.1
Type=Application
Name={{.Name}}
Icon=telegram
TryExec={{.TryExec}}
Exec={{.Exec}}
Terminal=false
Categories=Chat;Network;InstantMessaging;Qt;
Keywords=tg;chat;im;messaging;messenger;sms;tdesktop;
StartupWMClass=TelegramDesktop
X-GNOME-UsesNotifications=true
`))

type entryTemplateStruct struct {
	Name    string
	TryExec string
	Exec    string
}

// Path builds a path to the desktop entry
func Path(dir, name string) string {
	return path.Join(dir, fmt.Sprintf("telegramdesktop.%s.desktop", name))
}

// Exist checks whether the desktop entry exists
func Exist(dir, name string) (bool, error) {
	return util.Exist(Path(dir, name))
}

// Create creates a new desktop entry
func Create(dir, name, tryExec, exec string) error {
	path := Path(dir, name)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return entryTemplate.Execute(file, entryTemplateStruct{
		Name:    fmt.Sprintf("Telegram Desktop – %s", name),
		TryExec: tryExec,
		Exec:    exec,
	})
}

// Remove removes the desktop entry
func Remove(dir, name string) error {
	path := Path(dir, name)
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return os.RemoveAll(path)
}
