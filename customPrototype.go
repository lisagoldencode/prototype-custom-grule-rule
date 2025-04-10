package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type UserSelections struct {
	Adventure          bool
	America            bool
	CarsTrucksTractors bool
	Goodtimes          bool
	Grit               bool
	Home               bool
	Love               bool
	HeartBreak         bool
	Lessons            bool
	Rebellion          bool
	Recommendations    map[string]int
}

func (p *UserSelections) GetField(fieldName string) (bool, error) {
	val := reflect.ValueOf(p).Elem()
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return false, fmt.Errorf("field '%s' does not exist", fieldName)
	}

	if field.Kind() == reflect.Bool {
		return field.Bool(), nil
	}

	return false, fmt.Errorf("field '%s' is not a boolean", fieldName)
}

func (p *UserSelections) IsSongThemeMatch(songId string, songThemes ...string) bool {
	fmt.Println("=========")
	fmt.Println("Checking Matches: " + songId)

	for _, theme := range songThemes {
		boolValue, err := p.GetField(theme)

		if err != nil {
			panic(err)
		}
		if boolValue {
			fmt.Println("Match found!")
			return true
		}

		fmt.Println("No Matches")
	}
	return false
}

func (p *UserSelections) SetRecommendations(songId string, songThemes ...string) int {
	fmt.Println("------")
	fmt.Println("Counting Matches... (" + songId + ")")

	matchCount := 0
	for _, theme := range songThemes {
		boolValue, err := p.GetField(theme)

		if err != nil {
			panic(err)
		}
		if err != nil {
			panic(err)
		}

		if boolValue {
			matchCount += 1
			fmt.Println(theme+" --- Match found -", matchCount)
		} else {
			fmt.Println(theme)
		}
	}

	fmt.Println("\nMatches for song '"+songId+"':", matchCount)

	p.Recommendations[songId] = matchCount
	return matchCount
}

func main() {
	log.Printf("Begin Prototype for Song Qualification")

	// Create Grule DataContext and register the function
	dataCtx := ast.NewDataContext()

	userSelections := &UserSelections{
		Adventure:          true,
		America:            true,
		CarsTrucksTractors: false,
		Goodtimes:          false,
		Grit:               false,
		Home:               false,
		Love:               false,
		HeartBreak:         false,
		Lessons:            true,
		Rebellion:          false,
		Recommendations:    make(map[string]int),
	}

	dataCtx.Add("UserSelections", userSelections)

	// Project TODO: Automatically generate this rule from lamdba
	drls := `
    rule Check10000 "Take Me Home, Country Roads" salience 10 {
        when
           UserSelections.IsSongThemeMatch("10000", "Adventure", "America", "Home", "Lessons")
        then
            UserSelections.SetRecommendations("10000", "Adventure", "America", "Home", "Lessons");
            Retract("Check10000");
    }

	rule Check10001 "All My Ex's Live In Texas" salience 10 {
        when
           UserSelections.IsSongThemeMatch("10001", "Adventure", "HeartBreak", "Rebellion")
        then
            UserSelections.SetRecommendations("10001", "Adventure", "HeartBreak", "Rebellion");
            Retract("Check10001");
    }
    `

	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	bs := pkg.NewBytesResource([]byte(drls))
	err := ruleBuilder.BuildRuleFromResource("SongRecs", "0.0.1", bs)
	if err != nil {
		panic(err)
	}

	knowledgeBase, err := knowledgeLibrary.NewKnowledgeBaseInstance("SongRecs", "0.0.1")

	engine := engine.NewGruleEngine()
	err = engine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n==========")
	fmt.Println("Song recommendations for user: ")
	fmt.Println(userSelections.Recommendations)

}
