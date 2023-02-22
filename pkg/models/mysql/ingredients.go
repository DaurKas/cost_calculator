package mysql

import (
	"database/sql"
	"errors"

	"cost_calculator/pkg/models"
)

// SnippetModel - Определяем тип который обертывает пул подключения sql.DB
type IngredientModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *IngredientModel) AddIngredient(name string, content string, quantity int,
	remains int, price int, priceForQunatity float32,
	tag string, quantityType string) (int, error) {
	//stmt := `INSERT INTO ingredients (    name, content, quantity, remains, price, priceForQunatity, tag, quantityType	)
	//VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	stmt := `INSERT INTO ingredients(name, content, quantity, remains, price, priceForQunatity, tag, quantityType) VALUES (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
    );`
	result, err := m.DB.Exec(stmt, name, content, quantity, remains, price, priceForQunatity, tag, quantityType)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *IngredientModel) GetIngredient(id int) (*models.Ingredient, error) {
	//stmt := `SELECT id, title, content, created, expires FROM snippets
	//WHERE expires > UTC_TIMESTAMP() AND id = ?`
	stmt := `SELECT id, name, content, quantity, remains, price, priceForQunatity, tag, quantityType FROM ingredients 
	WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	i := &models.Ingredient{}
	err := row.Scan(&i.Ing_id, &i.Ing_name, &i.Quantity, &i.QuantityType, &i.Price, &i.PriceForQunatity, &i.Tag, &i.Remains)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return i, nil
}

func (m *IngredientModel) GetIngredientList() ([]*models.Ingredient, error) {
	//stmt := `SELECT id, title, content, created, expires FROM snippets
	//WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	stmt := `SELECT id, name, quantity, quantityType, price, priceForQunatity, tag, remains FROM ingredients`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ingredients []*models.Ingredient

	for rows.Next() {
		i := &models.Ingredient{}
		err := rows.Scan(&i.Ing_id, &i.Ing_name, &i.Quantity, &i.QuantityType, &i.Price, &i.PriceForQunatity, &i.Tag, &i.Remains)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ingredients, nil
}
