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
	Recommendations    []string
}

func (p *UserSelections) GetField(fieldName string) (bool, error) {
	val := reflect.ValueOf(p).Elem()
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return false, fmt.Errorf("field %s does not exist", fieldName)
	}

	if field.Kind() == reflect.Bool {
		return field.Bool(), nil
	}

	return false, fmt.Errorf("field %s is not a boolean", fieldName)
}

func (p *UserSelections) Append(songId string) {
	p.Recommendations = append(p.Recommendations, songId)
}

// Custom function to check if adventure and good times are true
func (p *UserSelections) IsSongThemeMatch(songThemes ...string) bool {
	fmt.Println("=========")
	fmt.Println("UserSelections")
	fmt.Println(p)

	matchCount := 0
	for _, theme := range songThemes {
		boolValue, err := p.GetField(theme)

		if err != nil {
			panic(err)
		}

		if boolValue {
			fmt.Printf(theme + " --- Match found!")
			matchCount += 1
		} else {
			fmt.Printf(theme)
		}
	}

	fmt.Printf("\nMatches: %d", matchCount)
	if matchCount > 0 {
		return true
	}

	return false
}

func main() {
	log.Printf("Begin Prototype for Song Qualification")

	// Create Grule DataContext and register the function
	dataCtx := ast.NewDataContext()

	userSelections := &UserSelections{
		Adventure:          true,
		America:            true,
		CarsTrucksTractors: false,
		Goodtimes:          true,
		Grit:               false,
		Home:               false,
		Love:               true,
		HeartBreak:         false,
		Lessons:            true,
		Rebellion:          false,
		Recommendations:    []string{},
	}

	dataCtx.Add("UserSelections", userSelections)

	// Project TODO: Automatically generate this rule from lamdba
	drls := `
    rule Check10000 "Take Me Home, Country Roads" salience 10 {
        when
           UserSelections.IsSongThemeMatch("Adventure", "America", "Home", "Lessons")
        then
            UserSelections.Append("10000");
            Retract("Check10000");
    }

	rule Check10001 "All My Ex’s Live In Texas" salience 10 {
        when
           UserSelections.IsSongThemeMatch("America", "HeartBreak", "Rebellion")
        then
            UserSelections.Append("10001");
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

	fmt.Println("=======================================")
	fmt.Println("Song recommendations for user: ")

	if len(userSelections.Recommendations) > 0 {
		for _, rec := range userSelections.Recommendations {
			fmt.Println(rec)
		}

	} else {
		fmt.Println("No songs matched the user's selections.")
	}
}
