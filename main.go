package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"git.sr.ht/~spc/go-log"
)

const AgentDir = "./agents"

func main() {
	log.SetLevel(log.LevelTrace)

	infos, err := ioutil.ReadDir(AgentDir)
	if err != nil {
		log.Errorf("cannot read file: %v", err)
		os.Exit(1)
	}

	output := make(map[string][]byte)
	for _, info := range infos {
		cmd := exec.Command(filepath.Join(AgentDir, info.Name()))
		cmd.Env = []string{
			"PATH=/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin",
		}

		buf, err := cmd.Output()
		if err != nil {
			log.Errorf("failed to run agent %v: %v", info.Name(), err)
			continue
		}

		output[info.Name()] = buf
	}

	data, err := json.Marshal(output)
	if err != nil {
		log.Errorf("cannot marshal JSON: %v", err)
		os.Exit(1)
	}
	log.Tracef("%+v", output)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3000", bytes.NewReader(data))
	if err != nil {
		log.Errorf("cannot create HTTP request: %v", err)
	}
	log.Tracef("sending HTTP request: %+v", req)
}
