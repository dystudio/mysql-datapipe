package canal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/siddontang/go-mysql/mysql"
)

var (
	ErrHandleInterrupted = errors.New("do handler error, interrupted")
)

// add for all QueryEvent
type QueryEventHandler interface {
	// Handle RowsEvent, if return ErrHandleInterrupted, canal will
	// stop the sync
	Do(e *QueryEvent) error
	String() string
	Close()
}

func (c *Canal) RegQueryEventHandler(h QueryEventHandler) {
	c.rsLock.Lock()
	c.quHandlers = append(c.quHandlers, h)
	c.rsLock.Unlock()
}

func (c *Canal) travelQueryEventHandler(e *QueryEvent) error {
	c.rsLock.Lock()
	defer c.rsLock.Unlock()

	// 过滤表，按照DUMP区域的配置，如果符合DUMP过程中过滤原则，那么其ROWS也不会通过
	var exists bool

	if len(c.dumper.Tables) != 0 {
		if !strings.EqualFold(c.dumper.TableDB, e.Table.Schema) {
			log.Infof("Table (%s.%s) not Match Config,Ignore ...", e.Table.Schema, e.Table.Name)
			return nil
		}

		for _, tb := range c.dumper.Tables {
			if strings.EqualFold(tb, e.Table.Name) {
				exists = true
				break
			}
		}

		if !exists {
			log.Infof("Table (%s.%s) not Match Config,Ignore ...", e.Table.Schema, e.Table.Name)
			return nil
		}
	}

	if len(c.dumper.Tables) == 0 && len(c.dumper.Databases) != 0 {
		for _, db := range c.dumper.Databases {
			if strings.EqualFold(db, e.Table.Schema) {
				exists = true
				break
			}
		}

		if !exists {
			log.Infof("Table (%s.%s) not Match Config,Ignore ...", e.Table.Schema, e.Table.Name)
			return nil
		}
	}

	if len(c.cfg.Dump.IgnoreTables) != 0 {
		itbi := fmt.Sprintf("%s.%s", e.Table.Schema, e.Table.Name)
		for _, itb := range c.cfg.Dump.IgnoreTables {
			if strings.EqualFold(itbi, itb) {
				exists = true
				break
			}
		}

		if exists {
			log.Infof("Table (%s.%s) not Match Config,Ignore ...", e.Table.Schema, e.Table.Name)
			return nil
		}
	}
	/*
		var tableExist bool
		tableExist = true

		for _, tb := range c.cfg.Dump.IgnoreTables {
			if strings.EqualFold(tb, e.Table.Name) {
				tableExist = false
			}
		}

		if !tableExist {
			log.Debugf("table name(%s.%s) not match ignore...",
				e.Table.Schema, e.Table.Name)
			return nil
		}

		tableExist = false

		if len(c.cfg.Dump.Tables) == 0 {
			tableExist = true
		} else {
			for _, tb := range c.cfg.Dump.Tables {
				if strings.EqualFold(tb, e.Table.Name) {
					tableExist = true
				}
			}
		}

		if !strings.EqualFold(e.Table.Schema, c.cfg.Dump.TableDB) || !tableExist {
			log.Debugf("table name(%s.%s) not match ignore...",
				e.Table.Schema, e.Table.Name)
			return nil
		}
	*/

	var err error
	for _, h := range c.quHandlers {
		if err = h.Do(e); err != nil && !mysql.ErrorEqual(err, ErrHandleInterrupted) {
			log.Errorf("handle %v err: %v", h, err)
		} else if mysql.ErrorEqual(err, ErrHandleInterrupted) {
			log.Errorf("handle %v err, interrupted", h)
			return ErrHandleInterrupted
		}

	}
	return nil
}

type RowsEventHandler interface {
	// Handle RowsEvent, if return ErrHandleInterrupted, canal will
	// stop the sync
	Do(e *RowsEvent) error
	String() string
	Close()
}

func (c *Canal) RegRowsEventHandler(h RowsEventHandler) {
	c.rsLock.Lock()
	c.rsHandlers = append(c.rsHandlers, h)
	c.rsLock.Unlock()
}

func (c *Canal) travelRowsEventHandler(e *RowsEvent) error {
	c.rsLock.Lock()
	defer c.rsLock.Unlock()

	// 过滤表，按照DUMP区域的配置，如果符合DUMP过程中过滤原则，那么其ROWS也不会通过
	var exists bool

	if len(c.dumper.Tables) != 0 {
		if !strings.EqualFold(c.dumper.TableDB, e.Table.Schema) {
			log.Infof("Table (%s.%s) not Match Config,Ignore ...", e.Table.Schema, e.Table.Name)
			return nil
		}

		for _, tb := range c.dumper.Tables {
			if strings.EqualFold(tb, e.Table.Name) {
				exists = true
				break
			}
		}

		if !exists {
			log.Infof("Table (%s.%s) not Match Config,Ignore ...", e.Table.Schema, e.Table.Name)
			return nil
		}
	}

	if len(c.dumper.Tables) == 0 && len(c.dumper.Databases) != 0 {
		for _, db := range c.dumper.Databases {
			if strings.EqualFold(db, e.Table.Schema) {
				exists = true
				break
			}
		}

		if !exists {
			log.Infof("Table (%s.%s) not Match Config,Ignore ...", e.Table.Schema, e.Table.Name)
			return nil
		}
	}

	if len(c.cfg.Dump.IgnoreTables) != 0 {
		itbi := fmt.Sprintf("%s.%s", e.Table.Schema, e.Table.Name)
		for _, itb := range c.cfg.Dump.IgnoreTables {
			if strings.EqualFold(itbi, itb) {
				exists = true
				break
			}
		}

		if exists {
			log.Infof("Table (%s.%s) not Match Config,Ignore ...", e.Table.Schema, e.Table.Name)
			return nil
		}
	}

	/*
		var tableExist bool
		tableExist = true

		for _, tb := range c.cfg.Dump.IgnoreTables {
			if strings.EqualFold(tb, e.Table.Name) {
				tableExist = false
			}
		}

		if !tableExist {
			log.Debugf("table name(%s.%s) not match ignore...",
				e.Table.Schema, e.Table.Name)
			return nil
		}

		tableExist = false

		if len(c.cfg.Dump.Tables) == 0 {
			tableExist = true
		} else {
			for _, tb := range c.cfg.Dump.Tables {
				if strings.EqualFold(tb, e.Table.Name) {
					tableExist = true
				}
			}
		}

		if !strings.EqualFold(e.Table.Schema, c.cfg.Dump.TableDB) || !tableExist {
			log.Debugf("table name(%s.%s) not match ignore...",
				e.Table.Schema, e.Table.Name)
			return nil
		}
	*/
	var err error
	for _, h := range c.rsHandlers {
		if err = h.Do(e); err != nil && !mysql.ErrorEqual(err, ErrHandleInterrupted) {
			log.Errorf("handle %v err: %v", h, err)
		} else if mysql.ErrorEqual(err, ErrHandleInterrupted) {
			log.Errorf("handle %v err, interrupted", h)
			return ErrHandleInterrupted
		}

	}
	return nil
}

// []byte, int64, float64, bool, string
func InterfaceToString(s interface{},RawType string) string {
	// Handle the most common destination types using type switches and
	// fall back to reflection for all other types.
	switch s := s.(type) {
	case nil:
		return "NULL"
	case string:
		return s
	case []byte:
		return string(s)
	case bool:
		return strconv.FormatBool(s)
	case int:
		if strings.Index(strings.ToUpper(RawType),"UNSIGNED") != -1 {
			log.Infof("int8: %v",strconv.FormatUint(uint64(uint8(s)), 10))
			return strconv.FormatUint(uint64(uint(s)), 10)
		}
		return strconv.FormatInt(int64(s), 10)
	case int8:
		if strings.Index(strings.ToUpper(RawType),"UNSIGNED") != -1 {
			log.Infof("int8: %v",strconv.FormatUint(uint64(uint8(s)), 10))
			return strconv.FormatUint(uint64(uint8(s)), 10)
		}
		return strconv.FormatInt(int64(s), 10)
	case uint8:
		return strconv.FormatUint(uint64(s), 10)
	case int16:
		if strings.Index(strings.ToUpper(RawType),"UNSIGNED") != -1 {
			log.Infof("int16: %v",strconv.FormatUint(uint64(uint16(s)), 10))
			return strconv.FormatUint(uint64(uint16(s)), 10)
		}
		return strconv.FormatInt(int64(s), 10)
	case uint16:
		return strconv.FormatUint(uint64(s), 10)
	case int32:
		if strings.Index(strings.ToUpper(RawType),"UNSIGNED") != -1 {
			log.Infof("int32: %v",strconv.FormatUint(uint64(uint32(s)), 10))
			return strconv.FormatUint(uint64(uint32(s)), 10)
		}
		return strconv.FormatInt(int64(s), 10)
	case uint32:
		return strconv.FormatUint(uint64(s), 10)
	case int64:
		if strings.Index(strings.ToUpper(RawType),"UNSIGNED") != -1 {
			log.Infof("int64: %v",strconv.FormatUint(uint64(s), 10))
			return strconv.FormatUint(uint64(s), 10)
		}
		return strconv.FormatInt(int64(s), 10)
	case uint64:
		return strconv.FormatUint(uint64(s), 10)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', 4, 32)
	case float64:
		return strconv.FormatFloat(s, 'f', 4, 64)
	case time.Time:
		return s.Format(mysql.TimeFormat)
	}

	return "nil"
}
