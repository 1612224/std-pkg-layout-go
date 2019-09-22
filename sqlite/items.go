package sqlite

import (
	"database/sql"
	"log"
	app "useritem"
)

// ItemRepo is a Sqlite specific implementation of the item repository
type ItemRepo struct {
	DB *sql.DB
}

// ByUser will look for all items that belong to an user with specific user id
// return slice of app.Item and an error
func (repo *ItemRepo) ByUser(userID int) ([]app.Item, error) {
	rows, err := repo.DB.Query("select userid,name,price from items where userid=?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []app.Item
	for rows.Next() {
		var item app.Item
		err = rows.Scan(&item.UserID, &item.Name, &item.Price)
		if err != nil {
			log.Printf("Failed to scan item: %v\n", err)
			continue
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// Create insert new item into database
// return an error
func (repo *ItemRepo) Create(item *app.Item) error {
	_, err := repo.DB.Exec("insert into items(userid,name,price) values (?,?,?)", item.UserID, item.Name, item.Price)
	return err
}
