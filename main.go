package main

import (
	"fmt"
	"deck"
	"blackjack"
)

type blackjackAI struct {
	score	int
	seen	int
	decks	int
}

func (ai *blackjackAI) Bet(shuffled bool) int {
	if shuffled {
		ai.score  = 0
		ai.seen = 0
	}
	trueScore := ai.score / ((ai.decks * 52 - ai.seen) / 52)
	switch {
	case trueScore > 14:
		return 100000
	case trueScore > 8:
		return 5000
	default:
		return 100
	}
}

func (ai *blackjackAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	pScore := blackjack.Score(hand...)
	dScore := blackjack.Score(dealer)
	if len(hand) == 2 {
		if hand[0] == hand[1] {
			cardScore := blackjack.Score(hand[0])
			if cardScore != 5 && cardScore != 10 {
				if (cardScore >= 8 && dScore != 7 && dScore != 10 && dScore != 11) || (cardScore == 7 && !(dScore < 7)) || (cardScore == 6 && !(dScore == 2) && !(dScore < 7)) || (cardScore <= 3 && dScore > 3 && dScore < 8) {
					return blackjack.MoveSplit
				}
				if cardScore == 9 && (dScore == 7 || dScore >= 10) {
					return blackjack.MoveStand
				}
				return blackjack.MoveHit
			}

		}
		if blackjack.Soft(hand...) {
			if (pScore >= 20) || (pScore == 19 && dScore != 6) || (pScore == 18 && dScore > 6 && dScore < 9) {
				return blackjack.MoveStand
			} else if (pScore == 18 && dScore ==2) || (pScore > 16 && pScore < 19 && dScore ==3) || (pScore > 14 && pScore < 19 && dScore ==4) || (pScore < 19 && dScore == 5) || (pScore < 20 && dScore == 6) {
				return blackjack.MoveDouble
			} else {
				return blackjack.MoveHit
			}
		}
		if (pScore == 11) || (pScore == 10 && dScore < 10) || (pScore == 9 && dScore != 2 && dScore < 7) {
			return blackjack.MoveDouble
		}
	}
	if (pScore < 13 && dScore < 4) || (pScore < 12) || (pScore > 11 && pScore < 17 && dScore > 6) {
		return blackjack.MoveHit
	} else {
		return blackjack.MoveStand
	}
}

func (ai *blackjackAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	for _, card := range dealer {
		ai.count(card)
	}
	for _, hand := range hands {
		for _, card := range hand {
			ai.count(card)
		}
	}
}

func (ai *blackjackAI) count(card deck.Card) {
	score := blackjack.Score(card)
	switch {
	case score >= 10:
		ai.score--
	case score <= 6:
		ai.score++
	}
	ai.seen++
}


func main() {
	fmt.Println("Would you like to play or see how much money an AI can make?\n[play/ai]")
	var input string
	fmt.Scanf("%s\n", &input)
	switch input {
	case "play":
		fmt.Println("How many hands would you like to play?")
		var hands int
		fmt.Scanf("%d\n", &hands)
		opts := blackjack.Options{
			Decks:           4,
			Hands:           hands,
			BlackjackPayout: 1.5,
		}
		game := blackjack.New(opts)
		winnings := game.Play(blackjack.HumanAI())
		var endMessage string
		if winnings < 0 {
			endMessage = fmt.Sprintf("You lost $%d", winnings * -1)
		} else {
			endMessage = fmt.Sprintf("You won $%d", winnings)
		}
		fmt.Println(endMessage)
	case "ai":
		opts := blackjack.Options{
			Decks:           4,
			Hands:           1000000,
			BlackjackPayout: 1.5,
		}
		game := blackjack.New(opts)
		winnings := game.Play(&blackjackAI{
			decks:	4,
		})
		var endMessage string
		if winnings < 0 {
			endMessage = fmt.Sprintf("In 1000000 hands, the AI lost $%d", winnings * -1)
		} else {
			endMessage = fmt.Sprintf("In 1000000 hands, the AI won $%d", winnings)
		}
		fmt.Println(endMessage)
	}
}