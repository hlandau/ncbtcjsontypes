// This package contains extensions to btcjson used by ncdns.
package ncbtcjsontypes

import "github.com/hlandauf/btcjson"
import "encoding/json"
import "fmt"

// name_show

type NameShowCmd struct {
	id   interface{}
	Name string `json:"name"`
}

func NewNameShowCmd(id interface{}, name string) (*NameShowCmd, error) {
	return &NameShowCmd{
		id:   id,
		Name: name,
	}, nil
}

func (c *NameShowCmd) Id() interface{} {
	return c.id
}

func (c *NameShowCmd) Method() string {
	return "name_show"
}

func (c *NameShowCmd) MarshalJSON() ([]byte, error) {
	params := []interface{}{
		c.Name,
	}

	raw, err := btcjson.NewRawCmd(c.id, c.Method(), params)
	if err != nil {
		return nil, err
	}

	return json.Marshal(raw)
}

func (c *NameShowCmd) UnmarshalJSON(b []byte) error {
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	if len(r.Params) != 1 {
		return btcjson.ErrWrongNumberOfParams
	}

	var name string
	if err := json.Unmarshal(r.Params[0], &name); err != nil {
		return fmt.Errorf("first argument 'name' must be a string: %v", err)
	}

	newCmd, err := NewNameShowCmd(r.Id, name)
	if err != nil {
		return err
	}

	*c = *newCmd
	return nil
}

func showCmdParser(rc *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(rc.Params) < 1 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	var name string
	if err := json.Unmarshal(rc.Params[0], &name); err != nil {
		return nil, fmt.Errorf("first argument 'name' must be a string: %v", err)
	}

	return NewNameShowCmd(rc.Id, name)
}

type NameShowReply struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Height    int    `json:"height"`
	ExpiresIn int    `json:"expires_in"`
	Expired   bool   `json:"expired"`
	Address   string `json:"address"`
	TxID      string `json:"txid"`
	VOut      int    `json:"vout"`
}

func showReplyParser(m json.RawMessage) (interface{}, error) {
	nsr := &NameShowReply{}
	err := json.Unmarshal(m, nsr)
	if err != nil {
		return nil, err
	}

	return nsr, nil
}

// name_sync

type NameSyncCmd struct {
	id        interface{}
	BlockHash string `json:"hash"`
	Count     int
	Wait      bool
}

func NewNameSyncCmd(id interface{}, blockHash string, count int, wait bool) (*NameSyncCmd, error) {
	return &NameSyncCmd{
		id:        id,
		BlockHash: blockHash,
		Count:     count,
		Wait:      wait,
	}, nil
}

func (c *NameSyncCmd) Id() interface{} {
	return c.id
}

func (c *NameSyncCmd) Method() string {
	return "name_sync"
}

func (c *NameSyncCmd) MarshalJSON() ([]byte, error) {
	params := []interface{}{
		c.BlockHash,
		c.Count,
		c.Wait,
	}

	raw, err := btcjson.NewRawCmd(c.id, c.Method(), params)
	if err != nil {
		return nil, err
	}

	return json.Marshal(raw)
}

func (c *NameSyncCmd) UnmarshalJSON(b []byte) error {
	// We don't need to implement this as we are only ever the client.
	panic("not implemented")
	return nil
}

type NameSyncReply []NameSyncEvent
type NameSyncEvent struct {
	Type string // "firstupdate" or "update" or "atblock"

	// Used for firstupdate and update.
	Name  string // "d/example"
	Value string // "..."

	// Used for atblock.
	BlockHash   string // in hex
	BlockHeight int
}

var errMalformed = fmt.Errorf("malformed name_sync event")

func (e *NameSyncEvent) UnmarshalJSON(data []byte) error {
	a := []interface{}{}
	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	if len(a) < 1 {
		return errMalformed
	}

	eventType, ok := a[0].(string)
	if !ok {
		return errMalformed
	}

	e.Type = eventType

	switch eventType {
	case "firstupdate", "update":
		if len(a) < 3 {
			return errMalformed
		}

		k, ok := a[1].(string)
		if !ok {
			return errMalformed
		}

		v, ok := a[2].(string)
		if !ok {
			return errMalformed
		}

		e.Name = k
		e.Value = v

	case "atblock":
		if len(a) < 3 {
			return errMalformed
		}

		hash, ok := a[1].(string)
		if !ok {
			return errMalformed
		}

		heightf, ok := a[2].(float64)
		if !ok {
			return errMalformed
		}

		height := int(heightf)

		e.BlockHash = hash
		e.BlockHeight = height
	}

	return nil
}

func syncReplyParser(m json.RawMessage) (interface{}, error) {
	nsr := NameSyncReply{}
	err := json.Unmarshal(m, &nsr)
	if err != nil {
		return nil, err
	}

	return nsr, nil
}

// name_scan

type NameScanCmd struct {
	id    interface{}
	From  string
	Count int
}

func NewNameScanCmd(id interface{}, from string, count int) (*NameScanCmd, error) {
	return &NameScanCmd{
		id:    id,
		From:  from,
		Count: count,
	}, nil
}

func (c *NameScanCmd) Id() interface{} {
	return c.id
}

func (c *NameScanCmd) Method() string {
	return "name_scan"
}

func (c *NameScanCmd) MarshalJSON() ([]byte, error) {
	params := []interface{}{
		c.From,
		c.Count,
	}

	raw, err := btcjson.NewRawCmd(c.id, c.Method(), params)
	if err != nil {
		return nil, err
	}

	return json.Marshal(raw)
}

func (c *NameScanCmd) UnmarshalJSON(b []byte) error {
	panic("not implemented")
	return nil
}

// name_filter

type NameFilterCmd struct {
	id     interface{}
	Regexp string
	MaxAge int
	From   int
	Count  int
}

func NewNameFilterCmd(id interface{}, regexp string, maxage, from, count int) (*NameFilterCmd, error) {
	return &NameFilterCmd{
		id:     id,
		Regexp: regexp,
		MaxAge: maxage,
		From:   from,
		Count:  count,
	}, nil
}

func (c *NameFilterCmd) Id() interface{} {
	return c.id
}

func (c *NameFilterCmd) Method() string {
	return "name_filter"
}

func (c *NameFilterCmd) MarshalJSON() ([]byte, error) {
	params := []interface{}{
		c.Regexp,
		c.MaxAge,
		c.From,
		c.Count,
	}

	raw, err := btcjson.NewRawCmd(c.id, c.Method(), params)
	if err != nil {
		return nil, err
	}

	return json.Marshal(raw)
}

func (c *NameFilterCmd) UnmarshalJSON(b []byte) error {
	// We don't need to implement this as we are only ever the client.
	panic("not implemented")
	return nil
}

type NameFilterReply []NameFilterItem
type NameFilterItem struct {
	Name      string // "d/example"
	Value     string // "..."
	TxID      string
	Address   string
	Height    int
	ExpiresIn int `json:"expires_in"`
	Expired   bool
}

func filterReplyParser(m json.RawMessage) (interface{}, error) {
	nsr := NameFilterReply{}
	err := json.Unmarshal(m, &nsr)
	if err != nil {
		return nil, err
	}

	return nsr, nil
}

func init() {
	btcjson.RegisterCustomCmd("name_show", showCmdParser, showReplyParser, "name_show <name>")
	btcjson.RegisterCustomCmd("name_sync", nil, syncReplyParser, "name_sync <block-hash> <count> <wait?>")
	btcjson.RegisterCustomCmd("name_filter", nil, filterReplyParser, "name_filter <regexp> <maxage> <from> <count>")
	btcjson.RegisterCustomCmd("name_scan", nil, filterReplyParser, "name_scan <from> <count>")
}

// © 2014 Hugo Landau <hlandau@devever.net>    GPLv3 or later
