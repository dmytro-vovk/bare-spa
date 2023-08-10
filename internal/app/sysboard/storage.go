package sysboard

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/ifaces"
	log "github.com/sirupsen/logrus"
)

//func save(board Board) (err error) {
//	out, err := json.MarshalIndent(&board, "", "\t")
//	if err != nil {
//		return err
//	}
//
//	return ioutil.WriteFile(configFileName(board), out, 0600)
//}

func load(s ifaces.Sysboard) ifaces.Sysboard {
	src, err := os.ReadFile(configFileName(s))
	if err != nil {
		log.Printf("Can't read config file %q, using default", s.Type())
		return s
	}

	if err = json.Unmarshal(src, s); err != nil {
		log.Printf("Can't unmarshal data from config file %q, using default", s.Type())
		return s
	}

	log.Printf("Loaded config for board %q from %q", s.Type(), configFileName(s))
	return s
}

func deleteConfig(s ifaces.Sysboard) error { return os.Remove(configFileName(s)) }

func configFileName(s ifaces.Sysboard) string {
	return strings.ReplaceAll(strings.ToLower(s.Type()), " ", "-") + ".config.json"
}
