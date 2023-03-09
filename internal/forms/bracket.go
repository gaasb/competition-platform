package forms

type BracketForm struct {
}

type BracketsTypes struct {
}
type Bracker interface {
	Do()
}

func (b BracketsTypes) Do() {}

type bracketsType map[string]BracketsTypes

// bracketsType["elimination"].Do(ctx)

func CreateBracket() {}
func UpdateResult()  {}

func Shuffle() {
	// проверить если есть матчи то удалчить все и срандомит заного
}

//tournamet/id/bracket/id?start=true
