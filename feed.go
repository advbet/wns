package wns

type BetradarBetData struct {
	Type      string    `xml:"DocumentType,attr"`
	Timestamp Timestamp `xml:"Timestamp"`
	Sports    []Sport   `xml:"Sports>Sport"`
}

type Timestamp struct {
	Created string `xml:"CreatedTime,attr"`
	TZ      string `xml:"TimeZone,attr"`
}

type Sport struct {
	ID       int      `xml:"BetradarSportID,attr"`
	Category Category `xml:"Category"`
}

type Category struct {
	ID         int        `xml:"BetradarCategoryID,attr"`
	Country    string     `xml:"IsoName,attr"`
	Tournament Tournament `xml:"Tournament"`
}

type Tournament struct {
	ID    int    `xml:"BetradarTournamentID,attr"`
	Draws []Draw `xml:"Draws>Draw"`
	Bets  []Bet  `xml:"Draws>DrawOdds>Bet"`
}

type Draw struct {
	ID        int    `xml:"BetradarDrawID,attr"`
	DisplayID int    `xml:"DrawDisplayID,attr"`
	Type      string `xml:"DrawType,attr"`
	TimeType  string `xml:"TimeType,attr"`
	GameType  string `xml:"GameType,attr"`
	Name      string `xml:"GameName,attr"`
	DrawDate  string `xml:"Fixture>DateInfo>DrawDate"`
	//Result          []string `xml:"Result>ScoreInfo>Score"`
	Result          Result `xml:"Result"`
	Canceled        bool   `xml:"StatusInfo>Off"`
	BonusBalls      int    `xml:"BonusBalls,attr"`
	BonusBallsDrum  string `xml:"BonusBallsDrum,attr"`
	BonusBallsRange string `xml:"BonusBallsRange,attr"`
	// BetResult parsing
}

type Result struct {
	ScoreInfo []Score `xml:"ScoreInfo>Score"`
}

type Score struct {
	Type  string `xml:"Type,attr"`
	Value string `xml:",innerxml"`
}

type Bet struct {
	OddsType int    `xml:"OddsType,attr"`
	Odds     []Odds `xml:"Odds"`
}

type Odds struct {
	ID              int    `xml:"OutComeId,attr"`
	Outcome         string `xml:"OutCome,attr"`
	SpecialBetValue string `xml:"SpecialBetValue,attr"`
	Odds            string `xml:",innerxml"`
}
