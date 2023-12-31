package test

import (
	"testing"

	"github.com/gardusig/numberguesser/guesser"
	"github.com/gardusig/pandoraservice/pandora"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestServerSetup(t *testing.T) {
	server := pandora.NewPandoraServiceServer()
	logrus.Debug("pandora server: ", server)
	err := server.Start()
	if err != nil {
		t.Fatalf("Failed to start Pandora server: %v", err)
	}
	logrus.Debug("started server")
	client, err := pandora.NewPandoraServiceClient()
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	guesser := guesser.NewGuesser(client)
	logrus.Debug("created number guesser")
	openedPandoraBox, err := guesser.GetPandoraBox()
	if err != nil {
		t.Fatalf("failed to guess right number: %v", err)
	}
	if openedPandoraBox == nil {
		t.Fatalf("expected opened pandora box, got nil instead")
	}
	logrus.Debug("message: ", openedPandoraBox.Message)
}
