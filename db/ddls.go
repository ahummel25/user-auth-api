package main

var ddlStmts = []string{
	`
	CREATE TABLE IF NOT EXISTS users (
		user_id BLOB PRIMARY KEY, 
		first_name TEXT, 
		last_name TEXT,
		email TEXT
	)
	`,
}
