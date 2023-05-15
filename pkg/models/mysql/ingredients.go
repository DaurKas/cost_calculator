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
func (m *IngredientModel) AddIngredient(name string, content string, quantity float32,
	remains float32, price float32, priceForQunatity float32,
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
        ?
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
	err := row.Scan(&i.Ing_id, &i.Ing_name, &i.Quantity, &i.QuantityType, &i.Price, &i.PriceForQuantity, &i.Tag, &i.Remains)
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
		err := rows.Scan(&i.Ing_id, &i.Ing_name, &i.Quantity, &i.QuantityType, &i.Price, &i.PriceForQuantity, &i.Tag, &i.Remains)
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

func (m *IngredientModel) AddQuantity(id int, addition float32) error {
	stmt := `SELECT id, quantity FROM ingredients WHERE id = ?;`
	updateStmt := `UPDATE ingredients SET quantity = ? WHERE id = ?;`
	row := m.DB.QueryRow(stmt, id)
	var quantity float32
	err := row.Scan(&id, &quantity)
	newQuantity := quantity + addition

	if err != nil {
		return err
	}

	_, err = m.DB.Exec(updateStmt, newQuantity, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *IngredientModel) SubstractQuantity(id int, substraction float32) error {
	stmt := `SELECT id, quantity FROM ingredients WHERE id = ?;`
	updateStmt := `UPDATE ingredients SET quantity = ? WHERE id = ?;`
	row := m.DB.QueryRow(stmt, id)
	var quantity float32
	err := row.Scan(&id, &quantity)
	newQuantity := quantity - substraction

	if err != nil {
		return err
	}

	_, err = m.DB.Exec(updateStmt, newQuantity, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *IngredientModel) DeleteIngredient(id int) error {
	stmt := `DELETE FROM ingredients WHERE id = ?;`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
