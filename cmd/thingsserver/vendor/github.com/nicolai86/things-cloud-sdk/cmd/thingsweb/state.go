package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	thingscloud "github.com/nicolai86/things-cloud-sdk"
	"github.com/nicolai86/things-cloud-sdk/state/memory"
)

type state struct {
	*thingscloud.History
	*memory.State
}

func load(file string) (state, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return state{}, err
	}
	defer f.Close()
	var s state
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return state{}, err
	}
	err = json.Unmarshal(bs, &s)
	return s, err
}

func save(file string, s state) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	bs, err := json.Marshal(s)
	if err != nil {
		return err
	}
	f.Write(bs)
	return nil
}
