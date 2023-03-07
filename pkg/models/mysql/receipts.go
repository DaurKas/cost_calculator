package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cost_calculator/pkg/models"
)

type ReceiptModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *ReceiptModel) AddReceipt(name string, content string, quantity float32,
	tag string, quantityType string, ingrs []string, ingrQuants map[string]float32) (int, error) {

	ingrIdString := "(" + ingrs[0]
	for _, id := range ingrs[1:] {
		ingrIdString += ", " + id
	}
	ingrIdString += ");"

	stmrForRecIngrs := `SELECT id, name, quantity, quantityType, price, priceForQunatity, tag FROM ingredients
	WHERE id in` + ingrIdString
	fmt.Println(stmrForRecIngrs)
	ingrRes, err := m.DB.Query(stmrForRecIngrs)

	var reqIngrs []*models.Ingredient
	var totalPrice float32
	var priceForQunatity float32
	var ingredients string

	for ingrRes.Next() {
		ingr := &models.Ingredient{}
		ingrRes.Scan(&ingr.Ing_id, &ingr.Ing_name, &ingr.Quantity, &ingr.QuantityType, &ingr.Price, &ingr.PriceForQuantity, &ingr.Tag)
		reqIngrs = append(reqIngrs, ingr)
		costOfIngr := ingrQuants[strconv.Itoa(ingr.Ing_id)] * ingr.PriceForQuantity
		totalPrice += float32(costOfIngr)
		ingredients += strconv.Itoa(ingr.Ing_id) + ":" + strconv.FormatFloat(float64(ingrQuants[strconv.Itoa(ingr.Ing_id)]), 'f', 3, 32) + ";"

	}
	priceForQunatity = totalPrice / quantity

	stmt := `INSERT INTO receipts (name, ingredients, content, quantity, price, priceForQunatity, tag, quantityType) VALUES (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    );`
	result, err := m.DB.Exec(stmt, name, ingredients, content, quantity, totalPrice, priceForQunatity, tag, quantityType)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *ReceiptModel) GetReceipt(id int) (*models.Receipt, []*models.Ingredient, map[string]float32, error) {
	//stmt := `SELECT id, title, content, created, expires FROM snippets
	//WHERE expires > UTC_TIMESTAMP() AND id = ?`
	stmt := `SELECT id, name, ingredients, quantity, quantityType, price, priceForQunatity, tag FROM receipts 
	WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)
	fmt.Println(stmt)
	r := &models.Receipt{}
	var ingrText string
	err := row.Scan(&r.Rec_id, &r.Name, &ingrText, &r.Quantity, &r.QuantityType, &r.Price, &r.PriceForQuantity, &r.Tag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil, models.ErrNoRecord
		} else {
			return nil, nil, nil, err
		}
	}
	parsedIngr := parseIngredients(ingrText)
	ingrIdString := "("
	for ingId := range parsedIngr {
		ingrIdString += ingId + ", "
	}
	ingrIdString += "0);"

	stmrForRecIngrs := `SELECT id, name, quantity, remains, quantityType, price, priceForQunatity, tag FROM ingredients
	WHERE id in` + ingrIdString
	fmt.Println(stmrForRecIngrs)
	ingrRes, err := m.DB.Query(stmrForRecIngrs)
	var reqIngrs []*models.Ingredient
	for ingrRes.Next() {
		ingr := &models.Ingredient{}
		err = ingrRes.Scan(&ingr.Ing_id, &ingr.Ing_name, &ingr.Quantity, &ingr.Remains, &ingr.QuantityType, &ingr.Price, &ingr.PriceForQuantity, &ingr.Tag)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil, nil, models.ErrNoRecord
			} else {
				return nil, nil, nil, err
			}
		}
		reqIngrs = append(reqIngrs, ingr)
	}

	return r, reqIngrs, parsedIngr, nil
}
func parseIngredients(ingrStr string) map[string]float32 {
	ingrSep := strings.Split(ingrStr, ";")
	fmt.Println(ingrSep)
	splitMap := make(map[string]float32)
	for _, el := range ingrSep {
		if len(el) > 3 {
			idAndQuant := strings.Split(el, ":")
			quantFloat, _ := strconv.ParseFloat(idAndQuant[1], 32)
			splitMap[idAndQuant[0]] = float32(quantFloat)
		}
	}
	return splitMap
}

func (m *ReceiptModel) GetReceipts() ([]*models.Receipt, error) {
	stmt := `SELECT id, name, quantity, quantityType, price, priceForQunatity, tag FROM receipts`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var receipts []*models.Receipt

	for rows.Next() {
		r := &models.Receipt{}
		err := rows.Scan(&r.Rec_id, &r.Name, &r.Quantity, &r.QuantityType, &r.Price, &r.PriceForQuantity, &r.Tag)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return receipts, nil
}
