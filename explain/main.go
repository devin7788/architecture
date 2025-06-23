package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var dsn = "root:admin13@tcp(localhost:3306)/sqlopt_demo?parseTime=true"

func main() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(60 * time.Second)

	log.Println("连接成功，开始执行优化演示...")

	createProfileSession(db)
	insertTestData(db)
	queryWithExplain(db)
}

func createProfileSession(db *sql.DB) {
	_, _ = db.Exec("SET profiling = 1")
}

func queryWithExplain(db *sql.DB) {
	log.Println("执行 EXPLAIN 语句")

	// 不带索引的慢查询
	slowSQL := `SELECT * FROM users WHERE age + 1 = 30`

	// 推荐优化后的语句（避免函数）
	// fastSQL := `SELECT id, name FROM users WHERE age = 29`

	explainSQL := "EXPLAIN " + slowSQL

	rows, err := db.Query(explainSQL)
	if err != nil {
		log.Fatal("Explain 执行失败:", err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	vals := make([]interface{}, len(cols))
	valPtrs := make([]interface{}, len(cols))
	for i := range vals {
		valPtrs[i] = &vals[i]
	}

	fmt.Println("Explain 输出：")
	for rows.Next() {
		_ = rows.Scan(valPtrs...)
		for i, col := range cols {
			fmt.Printf("%s: %v\t", col, vals[i])
		}
		fmt.Println()
	}

	log.Println("执行实际 SQL 并获取 profile")
	_, _ = db.Exec(slowSQL)

	showProfiles(db)
}

func showProfiles(db *sql.DB) {
	rows, err := db.Query("SHOW PROFILES")
	if err != nil {
		log.Println("SHOW PROFILES 错误：", err)
		return
	}
	defer rows.Close()

	fmt.Println("执行时间分析（SHOW PROFILES）：")
	for rows.Next() {
		var queryID int
		var duration float64
		var query string
		_ = rows.Scan(&queryID, &duration, &query)
		fmt.Printf("QueryID=%d, Duration=%.6fs, SQL=%s\n", queryID, duration, query)
	}
}

func insertTestData(db *sql.DB) {
	total := 100_000
	batchSize := 1000

	log.Println("开始插入测试数据，共", total, "条...")

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("事务开始失败:", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO users (name, email, age, created_at) VALUES (?, ?, ?, ?)`)
	if err != nil {
		log.Fatal("准备语句失败:", err)
	}
	defer stmt.Close()

	start := time.Now()

	for i := 1; i <= total; i++ {
		name := fmt.Sprintf("User%d", i)
		email := fmt.Sprintf("user%d@example.com", i)
		age := 18 + (i % 50)
		createdAt := time.Now().Add(-time.Duration(i) * time.Hour)

		_, err := stmt.Exec(name, email, age, createdAt)
		if err != nil {
			log.Fatalf("插入第 %d 条失败: %v", i, err)
		}

		// 每 batchSize 提交一次事务
		if i%batchSize == 0 {
			err := tx.Commit()
			if err != nil {
				log.Fatal("事务提交失败:", err)
			}
			log.Printf("已插入 %d 条，耗时 %.2fs", i, time.Since(start).Seconds())

			// 开启下一批事务
			tx, err = db.Begin()
			if err != nil {
				log.Fatal("开启新事务失败:", err)
			}

			stmt, err = tx.Prepare(`INSERT INTO users (name, email, age, created_at) VALUES (?, ?, ?, ?)`)
			if err != nil {
				log.Fatal("新批次准备语句失败:", err)
			}
		}
	}

	// 提交最后剩余部分
	err = tx.Commit()
	if err != nil {
		log.Fatal("最后事务提交失败:", err)
	}

	log.Printf("插入完成，共 %d 条，用时 %.2fs", total, time.Since(start).Seconds())
}
