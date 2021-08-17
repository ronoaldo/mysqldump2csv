package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/xwb1989/sqlparser"
)

var (
	nullValueStr = ""
)

func main() {
	dump2csv(os.Stdin, os.Stdout)
}

func dump2csv(sql io.Reader, out io.Writer) error {
	t := sqlparser.NewTokenizer(sql)
	w := csv.NewWriter(out)
	for {
		stmt, err := sqlparser.ParseNext(t)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch stmt := stmt.(type) {
		case *sqlparser.Insert:
			log.Printf("Found table: '%v'", stmt.Table.Name)
			switch values := stmt.Rows.(type) {
			case sqlparser.Values:
				log.Printf("Value count: %d", len(values))
				for i := range values {
					valTuple := values[i]
					csvRow := []string{}
					for j := range valTuple {
						switch valTuple[j].(type) {
						case *sqlparser.SQLVal:
							val := valTuple[j].(*sqlparser.SQLVal)
							col := sqlparser.String(val)
							if val.Type == sqlparser.StrVal {
								col = string(val.Val)
							}
							csvRow = append(csvRow, col)
						case sqlparser.BoolVal:
							val := valTuple[j].(sqlparser.BoolVal)
							col := sqlparser.String(val)
							csvRow = append(csvRow, col)
						case *sqlparser.NullVal:
							csvRow = append(csvRow, nullValueStr)
						default:
							return fmt.Errorf("dump2csv: unexpected type %T", valTuple[j])
						}
					}
					if err := w.Write(csvRow); err != nil {
						return err
					}
					w.Flush()
				}
			}
		}
	}
	return nil
}
