package db

//func InitDatabase(driver string, username string, password string, host string, database string) (*sqlx.DB, error) {
//	dsn := fmt.Sprintf("%v://%v:%v@%v/%v?sslmode=disable",
//		driver, username, password, host, database,
//	)
//	db, err := sqlx.Open(driver, dsn)
//	if err != nil {
//		return nil, err
//	}
//	err = db.Ping()
//	if err != nil {
//		return nil, err
//	}
//	fmt.Println("Database Connected!!")
//
//	db.SetConnMaxLifetime(5 * time.Hour)
//	db.SetMaxOpenConns(10)
//	db.SetMaxIdleConns(10)
//
//	return db, nil
//}
