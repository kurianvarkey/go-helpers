# Go Helpers

This is a helper package. 

### DB log

This package correctly handles type formatting, single-quote escaping for strings, NULL values, and high-precision timestamptz values

### 🚀 Installation
Assuming your module path is github.com/kurianvarkey/go-helpers:

```
go get github.com/kurianvarkey/go-helpers
```

### 💡 Usage
Use GetBoundSQL right before calling your database execution function (db.QueryContext, etc.).

Example
```
package main

import (
	"context"
	"database/sql"
	"log"
	"time"	

	"github.com/kurianvarkey/go-helpers/logs" 
)

func main() {
    // Assume db is an initialized *sql.DB and ctx is context.Context
    db := getDBConnection() 
    ctx := context.Background()

    // Query components
    sql := "SELECT id, created_at FROM users WHERE id = $1 AND last_login < $2 AND status = $3"
    userID := 42
    loginThreshold := time.Now().Add(-24 * time.Hour)
    status := "ACTIVE"

    // 1. Log the query
    boundSQL := logs.GetBoundSQL(sql, userID, loginThreshold, status)
    log.Printf("[DEBUG] Executing SQL: %s", boundSQL)

    // 2. Execute the query (SECURELY)
    rows, err := db.QueryContext(ctx, sql, userID, loginThreshold, status)
    
    // ... handle results and errors
}
```

#### Example Log Output
```
[DEBUG] Executing SQL: SELECT id, created_at FROM users WHERE id = 42 AND last_login
```