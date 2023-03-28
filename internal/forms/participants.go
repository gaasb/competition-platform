package forms

type ParticipantForm struct {
	Team         string        `form:"team" json:"team,omitempty" binding:"lte=20,gte=2"`
	Participants []Participant `form:"participants" json:"participants" binding:"required,dive"`
}

type ParticipantsFromTeam struct {
	Team         string         `form:"team" json:"team" binding:"lte=20,gte=2" boil:"team_alias"`
	Participants []*Participant `json:"participants" boil:"-"`
}

type Participant struct {
	UserAlias string `form:"user_alias" json:"user_alias" binding:"required,lte=20,gte=2" boil:"user_alias"`
	Contact   string `form:"contact" json:"contact" binding:"omitempty,uri" boil:"contact"`
}
