// IN DEV...
package tasks

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Task struct {
	DB  *sql.DB
	Ctx context.Context
}

func (t *Task) ProcessFiles(dir string) {
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		fname := file.Name()
		switch {
		case strings.Contains(fname, "admin"):
			createTableADMIN(t.Ctx, t.DB)
		case strings.Contains(fname, "order_h"):
			createTableORDERH(t.Ctx, t.DB)
		case strings.Contains(fname, "order_i"):
			createTableORDERI(t.Ctx, t.DB)
		}
		tx, _ := t.DB.Begin()
		name := file.Name()
		f, _ := excelize.OpenFile("./files/" + fname)
		defer f.Close()
		stmt, _ := tx.PrepareContext(t.Ctx, "INSERT INTO znbol_admin VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		rows, _ := f.Rows("znbol_admin_tratada")

		for rows.Next() {
			if rows.CurrentRow() == 1 {
				continue
			}
			rowSlice, _ := rows.Columns()
			insertInto(t.Ctx, stmt, rowSlice...)
		}
		log.Println("Commiting file " + name + "...")
		tx.Commit()
	}
}

func insertInto[T string](ctx context.Context, stmt *sql.Stmt, row ...T) {
	stmt.Exec(row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9], row[10], row[11], row[12], row[13], row[14])

}

func createTableADMIN(ctx context.Context, db *sql.DB) {
	db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS admin ();
	`)
}

func createTableORDERH(ctx context.Context, db *sql.DB) {
	db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS order_h ();
	`)
}
func createTableORDERI(ctx context.Context, db *sql.DB) {
	db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS order_i ();
	`)
}
