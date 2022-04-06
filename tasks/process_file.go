// IN DEV...
package tasks

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/ericklima-ca/uperg/db"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

type TableMap map[string]string
var wg sync.WaitGroup
func ProcessFiles(ch chan<- bool, dir string) {
	files, _ := os.ReadDir(dir)
	
	for _, file := range files {
		fname := file.Name()
		go func(){
			wg.Add(1)
			defer wg.Done()
			readAndSaveData(fname)
		}()
	}
	wg.Wait()
	ch <- true
	defer close(ch)
}

func readAndSaveData(filename string) {
	

	f, _ := excelize.OpenFile("./tmp/uploads/" + filename)
	rows, err := f.Rows("Sheet1")
	if err != nil {
		log.Println(string(debug.Stack()))
	}
	defer rows.Close()
	rows.Next()
	firstRow, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	var mapTable = TableMap{}
	mapPattern := map[string]string{
		`(?i)^pedido`:    "TEXT",
		`(?i)quantidade`: "INTEGER",
		`(?i)data`:       "TEXT",
		`(?i)valor`:      "TEXT",
		`(?i)montante`:   "TEXT",
	}

	for _, col := range firstRow {
		for pattern, colType := range mapPattern {
			match, _ := regexp.MatchString(pattern, col)
			if match {
				mapTable[col] = colType
				break
			} else {
				mapTable[col] = "TEXT"
			}
		}
	}

	sqlTemplate := "CREATE TABLE IF NOT EXISTS %s (\n%s);"
	s := ""
	c := 0
	for k, v := range mapTable {
		c++
		s += fmt.Sprintf(`"%s" %s`, k, v)
		if c != len(mapTable) {
			s += ","
		}
	}
	fname := strings.Split(filename, ".")[0]
	createTableSql := fmt.Sprintf(sqlTemplate, fname, s)

	valuesTemplate := ""
	for i := 1; i <= len(firstRow); i++ {
		t := strconv.Itoa(i)
		valuesTemplate += "$"+t
		if i != len(firstRow) {
			valuesTemplate += ","
		}
	}
	insertTemplate := "INSERT INTO %s VALUES ( %s )"
	insertIntoSql := fmt.Sprintf(insertTemplate, fname, valuesTemplate)

	database := db.NewConnection(db.Options{
		DNS: os.Getenv("DATABASE_URL"),
	})
	tx, _ := database.Begin()
	if _, err := database.Exec(createTableSql); err != nil {
		log.Panic(string(debug.Stack()), err, createTableSql)
	}

	stmt, err := tx.Prepare(insertIntoSql)
	if err != nil {
		log.Panic(err, insertIntoSql)
	}

	for rows.Next() {
		if rows.CurrentRow() == 1 {
			continue
		}
		rowSlice, err := rows.Columns()
		if err != nil {
			log.Fatal(err)
		}
		execStatement(rowSlice, stmt)
	}
	log.Println("Commiting file " + filename + "...")

	tx.Commit()
}

func execStatement(args []string, stmt *sql.Stmt) {
	var listInterface = []interface{}{}
	for _, v := range args {
		listInterface = append(listInterface, v)
	}
	_, err := stmt.Exec(listInterface...)
	if err != nil {
		log.Fatal(err)
	}
}
