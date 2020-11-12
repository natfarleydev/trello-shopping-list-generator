package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/adlio/trello"
)

func mymain() error {
	key, _ := os.LookupEnv("TRELLO_KEY")
	token, _ := os.LookupEnv("TRELLO_TOKEN")
	client := trello.NewClient(key, token)

	board, err := client.GetBoard("WEp3YmB9", trello.Defaults())
	if err != nil {
		return err
	}

	lists, err := board.GetLists(trello.Defaults())

	mealsLists := map[string]*trello.List{}
	ingredientsLists := map[string]*trello.List{}
	for _, list := range lists {
		if strings.Contains(list.Name, "Meals") {
			mealsLists[list.Name] = list
		}
		if strings.Contains(list.Name, "Ingredients") {
			ingredientsLists[list.Name] = list
		}

		if len(mealsLists) == 2 && len(ingredientsLists) == 2 {
			break
		}
	}

	if _, ok := mealsLists["Meals: to buy"]; !ok {
		return fmt.Errorf("Meals: to want list is not present")
	}

	meals, err := mealsLists["Meals: to buy"].GetCards(trello.Defaults())
	if err != nil {
		return err
	}
	if len(meals) < 1 {
		return fmt.Errorf("Meals: to buy list is empty")
	}

	ingredients := map[string]*trello.Label{}
	for _, meal := range meals {
		for _, l := range meal.Labels {
			ingredients[l.Name] = l
		}
	}

	for _, il := range ingredientsLists {
		cards, err := il.GetCards(trello.Defaults())
		if err != nil {
			return err
		}

		for _, card := range cards {
			if _, ok := ingredients[card.Name]; ok {
				delete(ingredients, card.Name)
			}
		}

	}

	// Now we create ingredient cards that have not been created yet

	for ingredient := range ingredients {
		if _, ok := ingredientsLists["Ingredients: want"]; !ok {
			return fmt.Errorf("Ingredients: want list not found")
		}
		err = ingredientsLists["Ingredients: want"].AddCard(&trello.Card{
			Name: ingredient,
		}, trello.Defaults())
		if err != nil {
			return err
		}
	}

	// TODO move the ingredients into the right column if they're need for this week

	return nil

}

func main() {
	err := mymain()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
