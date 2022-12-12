package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 21 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	ingredients, potentialAllergens, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error when parsing the input data : %w", err)
	}

	// 3) Identify the safe ingredients
	safeIngredients := findSafeIngredients(ingredients, potentialAllergens)

	// 4) Count the values in the map safeIngredients to find the answer
	count := 0
	for _, val := range safeIngredients {
		count += val
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 21 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, potentialAllergens, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error when parsing the input data : %w", err)
	}

	// Will contain the mapping allergen -> ingredient
	allergens := make(map[string]string)

	for len(potentialAllergens) > 0 {

		for allergen, candidates := range potentialAllergens {
			if len(candidates) == 1 {
				// Use a loop to fetch only one element
				for ingredient := range candidates {
					allergens[allergen] = ingredient

					// Remove the ingredient from the candidates
					for _, cand := range potentialAllergens {
						delete(cand, ingredient)
					}
				}

				// Delete the allergen
				delete(potentialAllergens, allergen)

				break
			}
		}
	}

	// Create an ordered list of allergens
	var orderedAllergens []string
	for allergen := range allergens {
		orderedAllergens = append(orderedAllergens, allergen)
	}

	sort.Strings(orderedAllergens)

	var ingredientsList string

	for idx, allergen := range orderedAllergens {
		ingredientsList += allergens[allergen]

		if idx < len(orderedAllergens)-1 {
			ingredientsList += ","
		}
	}

	_, err = fmt.Fprintf(answer, "%s", ingredientsList)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

/*
* Principle of the algorithm I will use :
* 1) generate a map[string]int containing where the key is an ingredient and the value is the number of occurence in any list
* 2) generate a map[string]set where the key is an allergen and the value is the set of the ingredients that might contains this allergen
*  -> to generate this set, we can do it iteratively : for each list of ingredients containing this allergen, we keep the intersection of the set and the list of ingredients
* 3) we remove all ingredients appearing in the map[string]set of candidates to contain allergens of the map[string]int of the ingredients
* 4) The result is the sum of all values of map[string]set
 */

// Set implementation

type set map[string]bool

// PARSING INPUT

// We use regex to separate the ingredient part and the allergens part
var re = regexp.MustCompile(`^([a-z\s]+) \(contains ([a-z\s,]+)\)$`)

func parseLines(lines []string) (ingredients map[string]int, potentialAllergens map[string]set, err error) {

	ingredients = make(map[string]int)
	potentialAllergens = make(map[string]set)

	for _, line := range lines {

		// Separate the ingredients and allergens parts
		match := re.FindStringSubmatch(line)
		if len(match) != 3 {
			return nil, nil, fmt.Errorf("error parsing line %s : could not parse ingredients and allergens part", line)
		}

		currentIngredients := strings.Fields(match[1])

		// 1) Add the ingredients to the ingredients map
		for _, ingredient := range currentIngredients {
			ingredients[ingredient]++
		}

		// 2) Update the allergens candidates
		for _, allergen := range parseAllergens(match[2]) {

			_, exists := potentialAllergens[allergen]
			// If there are no potential candidates for this allergen
			if !exists {
				// Initialise the set
				candidatesSet := make(map[string]bool)

				// Add the ingredients as potential allergens in it
				for _, ingredient := range currentIngredients {
					candidatesSet[ingredient] = true
				}

				potentialAllergens[allergen] = candidatesSet
			} else {
				potentialAllergens[allergen] = updateCandidates(potentialAllergens[allergen], currentIngredients)
			}
		}
	}

	return ingredients, potentialAllergens, nil
}

// Parse the allergens part of the line
func parseAllergens(allergens string) (parsedAllergens []string) {

	for _, al := range strings.Split(allergens, ",") {
		parsedAllergens = append(parsedAllergens, strings.ReplaceAll(al, " ", ""))
	}

	return parsedAllergens
}

// Update the set of candidates for an allergen. Keep only the intersection of the previous candidates set and the list of current ingredients
func updateCandidates(candidates set, currentIngredients []string) (updatedSet set) {
	updatedSet = make(map[string]bool)

	for _, ingredient := range currentIngredients {
		// If the ingredient is not in the set, it won't be in the intersection
		if _, exists := candidates[ingredient]; exists {
			updatedSet[ingredient] = true
		}
	}

	return updatedSet
}

// Create a map of the safe ingredients
func findSafeIngredients(ingredients map[string]int, potentialAllergens map[string]set) (safeIngredients map[string]int) {
	safeIngredients = make(map[string]int)

	// For each ingredient, we look if i's a potential allergen for at least one allergen. If not, we append it to the map of safe ingredients.
	for ingredient, count := range ingredients {
		isSafe := true
		for _, candidates := range potentialAllergens {
			if _, exists := candidates[ingredient]; exists {
				isSafe = false
			}
		}

		if isSafe {
			safeIngredients[ingredient] = count
		}
	}

	return safeIngredients
}
