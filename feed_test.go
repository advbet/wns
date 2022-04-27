package wns

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFeed(t *testing.T) {
	tests := []struct {
		msg      string
		xml      string
		expected BetradarBetData
	}{
		{
			msg: "betradar Results type",
			xml: `
<?xml version="1.0" encoding="UTF-8"?>
<BetradarBetData DocumentType="Results">
  <Timestamp CreatedTime="Mon 2017-03-27 09:19:02" TimeZone="UTC"/>
  <Sports>
    <Sport BetradarSportID="108">
      <Category BetradarCategoryID="1103" IsoName="Malta">
        <Tournament BetradarTournamentID="49026">
          <Draws>
            <Draw BetradarDrawID="159777197" DrawDisplayID="217" DrawType="Rng" TimeType="Interval" GameType="20/80" GameName="Quick KENO 20/80">
              <Fixture>
                <DateInfo>
                  <DrawDate>2017-03-27 09:15:00</DrawDate>
                </DateInfo>
                <StatusInfo>
                  <Off>0</Off>
                </StatusInfo>
              </Fixture>
              <Result>
                <ScoreInfo SortType="Numerical">
                  <Score Type="draw_1">2</Score>
                  <Score Type="draw_2">3</Score>
                  <Score Type="draw_3">6</Score>
                  <Score Type="draw_4">10</Score>
                  <Score Type="draw_5">20</Score>
                  <Score Type="draw_6">21</Score>
                  <Score Type="draw_7">31</Score>
                  <Score Type="draw_8">36</Score>
                  <Score Type="draw_9">37</Score>
                  <Score Type="draw_10">38</Score>
                  <Score Type="draw_11">40</Score>
                  <Score Type="draw_12">43</Score>
                  <Score Type="draw_13">47</Score>
                  <Score Type="draw_14">49</Score>
                  <Score Type="draw_15">63</Score>
                  <Score Type="draw_16">69</Score>
                  <Score Type="draw_17">72</Score>
                  <Score Type="draw_18">73</Score>
                  <Score Type="draw_19">75</Score>
                  <Score Type="draw_20">76</Score>
                </ScoreInfo>
              </Result>
              <BetResult>
                <W OddsType="465" OutCome="Odd" OutComeId="9" SpecialBetValue="811"/>
                <W OddsType="466" OutCome="Over" OutComeId="6" SpecialBetValue="811"/>
                <W OddsType="468" OutCome="between 810-829" SpecialBetValue="811"/>
                <W OddsType="513" OutCome="No" OutComeId="5"/>
                <W OddsType="514" OutCome="No" OutComeId="5"/>
              </BetResult>
            </Draw>
          </Draws>
        </Tournament>
      </Category>
    </Sport>
  </Sports>
</BetradarBetData>
			`,
			expected: BetradarBetData{
				XMLName: xml.Name{Local: "BetradarBetData"},
				Type:    "Results",
				Timestamp: Timestamp{
					Created: "Mon 2017-03-27 09:19:02",
					TZ:      "UTC",
				},
				Sports: []Sport{
					{
						ID: 108,
						Category: Category{
							ID:      1103,
							Country: "Malta",
							Tournament: Tournament{
								ID: 49026,
								Draws: []Draw{
									{
										ID:        159777197,
										DisplayID: 217,
										Type:      "Rng",
										TimeType:  "Interval",
										GameType:  "20/80",
										Name:      "Quick KENO 20/80",
										DrawDate:  "2017-03-27 09:15:00",
										Canceled:  false,
										Result: Result{
											ScoreInfo: []Score{
												{
													Type:  "draw_1",
													Value: "2",
												},
												{
													Type:  "draw_2",
													Value: "3",
												},
												{
													Type:  "draw_3",
													Value: "6",
												},
												{
													Type:  "draw_4",
													Value: "10",
												},
												{
													Type:  "draw_5",
													Value: "20",
												},
												{
													Type:  "draw_6",
													Value: "21",
												},
												{
													Type:  "draw_7",
													Value: "31",
												},
												{
													Type:  "draw_8",
													Value: "36",
												},
												{
													Type:  "draw_9",
													Value: "37",
												},
												{
													Type:  "draw_10",
													Value: "38",
												},
												{
													Type:  "draw_11",
													Value: "40",
												},
												{
													Type:  "draw_12",
													Value: "43",
												},
												{
													Type:  "draw_13",
													Value: "47",
												},
												{
													Type:  "draw_14",
													Value: "49",
												},
												{
													Type:  "draw_15",
													Value: "63",
												},
												{
													Type:  "draw_16",
													Value: "69",
												},
												{
													Type:  "draw_17",
													Value: "72",
												},
												{
													Type:  "draw_18",
													Value: "73",
												},
												{
													Type:  "draw_19",
													Value: "75",
												},
												{
													Type:  "draw_20",
													Value: "76",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		{
			msg: "betradar Fixture type",
			xml: `
<?xml version="1.0" encoding="UTF-8"?>
<BetradarBetData DocumentType="Fixtures">
  <Timestamp CreatedTime="Thu 2017-03-23 00:00:04" TimeZone="UTC"/>
  <Sports>
    <Sport BetradarSportID="108">
      <Category BetradarCategoryID="1061" IsoName="Canada">
        <Tournament BetradarTournamentID="47111">
          <Draws>
            <Draw BetradarDrawID="159768499" DrawDisplayID="30" DrawType="Drum" TimeType="Fixed" GameType="20/70" GameName="Keno Atlantic 20/70">
              <Fixture>
                <DateInfo>
                  <DrawDate>2017-03-23 02:30:00</DrawDate>
                </DateInfo>
              </Fixture>
            </Draw>
            <DrawOdds>
              <Bet OddsType="465">
                <Odds OutCome="Odd" OutComeId="9">1.80</Odds>
                <Odds OutCome="Even" OutComeId="10">1.80</Odds>
              </Bet>
              <Bet OddsType="466">
                <Odds OutCome="Over" OutComeId="6" SpecialBetValue="710.5">1.80</Odds>
                <Odds OutCome="Under" OutComeId="7" SpecialBetValue="710.5">1.80</Odds>
              </Bet>
              <Bet OddsType="468">
                <Odds OutCome="between 210-529" OutComeId="669">77</Odds>
                <Odds OutCome="between 530-559" OutComeId="670">50</Odds>
                <Odds OutCome="between 560-589" OutComeId="671">23</Odds>
                <Odds OutCome="between 590-629" OutComeId="672">9</Odds>
                <Odds OutCome="between 630-659" OutComeId="673">8.00</Odds>
                <Odds OutCome="between 660-689" OutComeId="674">6.00</Odds>
                <Odds OutCome="between 690-709" OutComeId="675">8.50</Odds>
                <Odds OutCome="between 710-730" OutComeId="676">8.00</Odds>
                <Odds OutCome="between 731-760" OutComeId="677">6.00</Odds>
                <Odds OutCome="between 761-790" OutComeId="678">8.00</Odds>
                <Odds OutCome="between 791-830" OutComeId="679">9</Odds>
                <Odds OutCome="between 831-860" OutComeId="680">23</Odds>
                <Odds OutCome="between 861-890" OutComeId="681">50</Odds>
                <Odds OutCome="between 891-1210" OutComeId="682">77</Odds>
              </Bet>
              <Bet OddsType="517">
                <Odds OutCome="Hit" SpecialBetValue="1">3.20</Odds>
                <Odds OutCome="Hit" SpecialBetValue="2">11</Odds>
                <Odds OutCome="Hit" SpecialBetValue="3">40</Odds>
                <Odds OutCome="Hit" SpecialBetValue="4">150</Odds>
                <Odds OutCome="Hit" SpecialBetValue="5">500</Odds>
                <Odds OutCome="Hit" SpecialBetValue="6">2000</Odds>
                <Odds OutCome="Hit" SpecialBetValue="7">6500</Odds>
                <Odds OutCome="Hit" SpecialBetValue="8">18000</Odds>
              </Bet>
            </DrawOdds>
          </Draws>
        </Tournament>
      </Category>
    </Sport>
  </Sports>
</BetradarBetData>
			`,
			expected: BetradarBetData{
				XMLName: xml.Name{Local: "BetradarBetData"},
				Type:    "Fixtures",
				Timestamp: Timestamp{
					Created: "Thu 2017-03-23 00:00:04",
					TZ:      "UTC",
				},
				Sports: []Sport{
					{
						ID: 108,
						Category: Category{
							ID:      1061,
							Country: "Canada",
							Tournament: Tournament{
								ID: 47111,
								Draws: []Draw{
									{
										ID:        159768499,
										DisplayID: 30,
										Type:      "Drum",
										TimeType:  "Fixed",
										GameType:  "20/70",
										Name:      "Keno Atlantic 20/70",
										DrawDate:  "2017-03-23 02:30:00",
									},
								},
								Bets: []Bet{
									{
										OddsType: 465,
										Odds: []Odds{
											{
												ID:      9,
												Outcome: "Odd",
												Odds:    "1.80",
											},
											{
												ID:      10,
												Outcome: "Even",
												Odds:    "1.80",
											},
										},
									},
									{
										OddsType: 466,
										Odds: []Odds{
											{
												ID:              6,
												Outcome:         "Over",
												SpecialBetValue: "710.5",
												Odds:            "1.80",
											},
											{
												ID:              7,
												Outcome:         "Under",
												SpecialBetValue: "710.5",
												Odds:            "1.80",
											},
										},
									},
									{
										OddsType: 468,
										Odds: []Odds{
											{
												ID:      669,
												Outcome: "between 210-529",
												Odds:    "77",
											},
											{
												ID:      670,
												Outcome: "between 530-559",
												Odds:    "50",
											},
											{
												ID:      671,
												Outcome: "between 560-589",
												Odds:    "23",
											},
											{
												ID:      672,
												Outcome: "between 590-629",
												Odds:    "9",
											},
											{
												ID:      673,
												Outcome: "between 630-659",
												Odds:    "8.00",
											},
											{
												ID:      674,
												Outcome: "between 660-689",
												Odds:    "6.00",
											},
											{
												ID:      675,
												Outcome: "between 690-709",
												Odds:    "8.50",
											},
											{
												ID:      676,
												Outcome: "between 710-730",
												Odds:    "8.00",
											},
											{
												ID:      677,
												Outcome: "between 731-760",
												Odds:    "6.00",
											},
											{
												ID:      678,
												Outcome: "between 761-790",
												Odds:    "8.00",
											},
											{
												ID:      679,
												Outcome: "between 791-830",
												Odds:    "9",
											},
											{
												ID:      680,
												Outcome: "between 831-860",
												Odds:    "23",
											},
											{
												ID:      681,
												Outcome: "between 861-890",
												Odds:    "50",
											},
											{
												ID:      682,
												Outcome: "between 891-1210",
												Odds:    "77",
											},
										},
									},

									{
										OddsType: 517,
										Odds: []Odds{
											{
												Outcome:         "Hit",
												SpecialBetValue: "1",
												Odds:            "3.20",
											},
											{
												Outcome:         "Hit",
												SpecialBetValue: "2",
												Odds:            "11",
											},
											{
												Outcome:         "Hit",
												SpecialBetValue: "3",
												Odds:            "40",
											},
											{
												Outcome:         "Hit",
												SpecialBetValue: "4",
												Odds:            "150",
											},
											{
												Outcome:         "Hit",
												SpecialBetValue: "5",
												Odds:            "500",
											},
											{
												Outcome:         "Hit",
												SpecialBetValue: "6",
												Odds:            "2000",
											},
											{
												Outcome:         "Hit",
												SpecialBetValue: "7",
												Odds:            "6500",
											},
											{
												Outcome:         "Hit",
												SpecialBetValue: "8",
												Odds:            "18000",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		var actual BetradarBetData
		err := xml.Unmarshal([]byte(test.xml), &actual)
		require.NoError(t, err)
		assert.Equal(t, test.expected, actual)
	}
}
