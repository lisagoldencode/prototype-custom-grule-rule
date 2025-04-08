package main

import (
	"fmt"
	"log"

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
	Recommendations    string
}

// Custom function to check if adventure and good times are true
func (p *UserSelections) IsSongThemeMatch(songThemes ...string) bool {

	const userSelection = "User Selection contained "

	matchCount := 0
	for _, theme := range songThemes {
		fmt.Println(theme)

		if theme == "Adventure" && p.Adventure {
			log.Printf(userSelection + "Adventure")
			matchCount += 1
		}
		if theme == "CarsTrucksTractors" && p.CarsTrucksTractors {
			log.Printf(userSelection + "CarsTrucksTractors")
			matchCount += 1
		}
		if theme == "Goodtimes" && p.Goodtimes {
			log.Printf(userSelection + "Goodtimes")
			matchCount += 1
		}
		if theme == "Grit" && p.Grit {
			log.Printf(userSelection + "Grit")
			matchCount += 1
		}
		if theme == "Home" && p.Home {
			log.Printf(userSelection + "Home")
			matchCount += 1
		}
		if theme == "Love" && p.Love {
			log.Printf(userSelection + "Love")
			matchCount += 1
		}
		if theme == "HeartBreak" && p.HeartBreak {
			log.Printf(userSelection + "HeartBreak")
			matchCount += 1
		}
		if theme == "Lessons" && p.Lessons {
			log.Printf(userSelection + "Lessons")
			matchCount += 1
		}
		if theme == "Rebellion" && p.Rebellion {
			log.Printf(userSelection + "Rebellion")
			matchCount += 1
		}
	}

	if matchCount > 0 {
		log.Printf("%d", matchCount)
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
		America:            false,
		CarsTrucksTractors: false,
		Goodtimes:          true,
		Grit:               false,
		Home:               false,
		Love:               true,
		HeartBreak:         false,
		Lessons:            false,
		Rebellion:          false,
		Recommendations:    "",
	}

	dataCtx.Add("UserSelections", userSelections)

	// Project TODO: Automatically generate this rule from lamdba
	drls := `
    rule Check10000 "Check User Vibes" salience 10 {
        when
           UserSelections.IsSongThemeMatch("Adventure", "America", "Home", "Lessons")
        then
            UserSelections.Recommendations = "10000";
            Retract("Check10000");
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

	fmt.Println("Song recommendations for user: " + userSelections.Recommendations)
}
