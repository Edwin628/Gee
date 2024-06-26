package session

import (
	"fmt"
	"geeorm/log"
	"geeorm/schema"
	"reflect"
	"strings"
)

func (s *Session) Model(val interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(val) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(val, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.refTable
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw("CREATE TABLE %s (%s)", table.Name, desc).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw("DROP TABLE IF EXISTS %s", s.RefTable().Name).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, sqlVar := s.dialect.TableExistSQL(s.RefTable().Name)
	result := s.Raw(sql, sqlVar...).QueryRow()
	var tmp string
	result.Scan(&tmp)
	return tmp == s.RefTable().Name
}
