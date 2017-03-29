package wns

// BetradarBetData is a main structure for parsed document
type BetradarBetData struct {
	Type      string    `xml:"DocumentType,attr"`
	Timestamp Timestamp `xml:"Timestamp"`
	Sports    []Sport   `xml:"Sports>Sport"`
}

// Timestamp represents time when the docoment was created
type Timestamp struct {
	Created string `xml:"CreatedTime,attr"`
	TZ      string `xml:"TimeZone,attr"`
}

// Sport is betradar sport struct(ID for lotteries will (should?) always be 108)
type Sport struct {
	ID       int      `xml:"BetradarSportID,attr"`
	Category Category `xml:"Category"`
}

// Category holds country name and main game struct(tournament)
type Category struct {
	ID         int        `xml:"BetradarCategoryID,attr"`
	Country    string     `xml:"IsoName,attr"`
	Tournament Tournament `xml:"Tournament"`
}

// Tournament is a draws and bets sturcture holder, and also uniquely
// identifies a game
type Tournament struct {
	ID    int    `xml:"BetradarTournamentID,attr"`
	Draws []Draw `xml:"Draws>Draw"`
	Bets  []Bet  `xml:"Draws>DrawOdds>Bet"`
}

// Draw represents single draw instance for given game
type Draw struct {
	ID              int    `xml:"BetradarDrawID,attr"`
	DisplayID       int    `xml:"DrawDisplayID,attr"`
	Type            string `xml:"DrawType,attr"`
	TimeType        string `xml:"TimeType,attr"`
	GameType        string `xml:"GameType,attr"`
	Name            string `xml:"GameName,attr"`
	DrawDate        string `xml:"Fixture>DateInfo>DrawDate"`
	Result          Result `xml:"Result"`
	Canceled        bool   `xml:"StatusInfo>Off"`
	BonusBalls      int    `xml:"BonusBalls,attr"`
	BonusBallsDrum  string `xml:"BonusBallsDrum,attr"`
	BonusBallsRange string `xml:"BonusBallsRange,attr"`
	// BetResult parsing
}

// Result is a holder for Score struct
type Result struct {
	ScoreInfo []Score `xml:"ScoreInfo>Score"`
}

// Score lists drawn balls by type, which is either `draw_x` where x is a
// number for ball drawn, or `draw_bx` for bonus ball drawn.
type Score struct {
	Type  string `xml:"Type,attr"`
	Value string `xml:",innerxml"`
}

// Bet is a holder for odds struct
type Bet struct {
	OddsType int    `xml:"OddsType,attr"`
	Odds     []Odds `xml:"Odds"`
}

// Odds outlines possible odds for betting options
// see 2.7.3 Appendx A of manual
type Odds struct {
	ID              int    `xml:"OutComeId,attr"`
	Outcome         string `xml:"OutCome,attr"`
	SpecialBetValue string `xml:"SpecialBetValue,attr"`
	Odds            string `xml:",innerxml"`
}
