package main

import (
	"cost_calculator/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	i, err := app.ingredients.GetIngredientList()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	receipts, err := app.receipts.GetReceipts()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Ingredients: i,
		Receipts:    receipts,
	})

}
func (app *application) ShowIngredients(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		price, _ := strconv.ParseFloat(r.FormValue("price"), 32)
		quantity, _ := strconv.ParseFloat(r.FormValue("quantity"), 32)
		description := r.FormValue("description")
		IngrType := r.FormValue("type")
		tag := r.FormValue("tag")
		remain := quantity
		priceForQunatity := price / quantity
		id, err := app.ingredients.AddIngredient(name, description, float32(quantity), float32(remain), float32(price), float32(priceForQunatity), tag, IngrType)
		fmt.Println(id, err)
		//w.WriteHeader(303)
		http.Redirect(w, r, "/ingredients", http.StatusSeeOther)
	}

	i, err := app.ingredients.GetIngredientList()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Println("INGR LUST HANDLER ", r)

	app.render(w, r, "ingredientList.page.tmpl", &templateData{
		Ingredients: i,
	})
}
func (app *application) ShowReceipt(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	fmt.Println("IN THE REC HANDLER ", id)
	if err != nil || id < 1 {
		app.notFound(w) // Страница не найдена.
		return
	}

	rec, ingredients, reqMap, err := app.receipts.GetReceipt(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	var requirenments []*models.Requirenment
	for _, ingr := range ingredients {
		newReq := models.Requirenment{Ingr: ingr, Required: reqMap[strconv.Itoa(ingr.Ing_id)]}

		requirenments = append(requirenments, &newReq)
	}
	app.render(w, r, "showReceipt.page.tmpl", &templateData{
		Receipt:       rec,
		Requirenments: requirenments,
	})

}
func (app *application) ShowReceipts(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		name := r.FormValue("name")
		quantity, _ := strconv.ParseFloat(r.FormValue("quantity"), 32)
		description := r.FormValue("description")
		IngrType := r.FormValue("type")
		tag := r.FormValue("tag")
		ingredientList := r.Form["ingredients"]
		ingrQuants := make(map[string]float32)
		for _, ingr := range ingredientList {
			quantFloatVal, _ := strconv.ParseFloat(r.FormValue(ingr), 32)
			ingrQuants[ingr] = float32(quantFloatVal)
		}
		fmt.Print(r.Form)
		id, err := app.receipts.AddReceipt(name, description, float32(quantity), tag, IngrType, ingredientList, ingrQuants)
		fmt.Println(id, err)
		http.Redirect(w, r, "/receipts", http.StatusSeeOther)
	}

	i, err := app.ingredients.GetIngredientList()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	receipts, err := app.receipts.GetReceipts()
	fmt.Println("___________DEBUG____________", receipts)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Println("INGR LUST HANDLER ", r)

	app.render(w, r, "receiptMain.page.tmpl", &templateData{
		Ingredients: i,
		Receipts:    receipts,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Страница не найдена.
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.
	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"

	// Передаем данные в метод SnippetModel.Insert(), получая обратно
	// ID только что созданной записи в базу данных.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Перенаправляем пользователя на соответствующую страницу заметки.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
