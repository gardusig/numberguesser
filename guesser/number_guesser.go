package guesser

import (
	"fmt"

	"github.com/gardusig/numberguesser/internal"
	pandoraproto "github.com/gardusig/pandoraproto/generated/go"
	"github.com/gardusig/pandoraservice/pandora"
	"github.com/sirupsen/logrus"
)

type Guesser struct {
	pandoraClient *pandora.PandoraServiceClient

	level      uint32
	lowerBound int64
	upperBound int64
}

func NewGuesser(pandoraClient *pandora.PandoraServiceClient) *Guesser {
	return &Guesser{
		pandoraClient: pandoraClient,
	}
}

func (g *Guesser) GetPandoraBox() (*pandoraproto.OpenedPandoraBox, error) {
	var lockedBox *pandoraproto.LockedPandoraBox
	var err error
	g.level = internal.LevelMinThreshold
	for g.level <= internal.LevelMaxThreshold {
		lockedBox, err = g.guessNumberByLevel()
		if err != nil {
			return nil, err
		}
		if lockedBox != nil {
			g.level += 1
			logrus.Debug("Passed to level: ", g.level, ", encryptedMessage: ", lockedBox.EncryptedMessage)
		}
	}
	return g.pandoraClient.SendOpenBoxRequest(lockedBox)
}

func (g *Guesser) guessNumberByLevel() (*pandoraproto.LockedPandoraBox, error) {
	logrus.Debug("attempt to guess number for level: ", g.level)
	g.lowerBound = internal.GuessMinThreshold
	g.upperBound = internal.GuessMaxThreshold
	for g.lowerBound <= g.upperBound {
		guess := g.lowerBound + ((g.upperBound - g.lowerBound) >> 1)
		logrus.Debug("lowerBound:", g.lowerBound, ", upperBound:", g.upperBound, ", guess:", guess)
		resp, err := g.pandoraClient.SendGuessRequest(g.level, guess)
		if err != nil {
			return nil, err
		}
		logrus.Debug("server response:", resp.Result)
		switch resp.Result {
		case internal.Equal:
			return resp.LockedPandoraBox, nil
		case internal.Greater:
			g.upperBound = guess - 1
		case internal.Less:
			g.lowerBound = guess + 1
		default:
			return nil, fmt.Errorf("Unexpected response from server: %v", resp.Result)
		}
	}
	return nil, fmt.Errorf("Failed to guess the right number :/")
}
