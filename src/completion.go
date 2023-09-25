package src

import (
	"fmt"

	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/predict"
)

func numbers() []string {
	res := make([]string, 0)
	for i := 0; i < 24; i++ {
		res = append(res, fmt.Sprintf("%v", i))
	}
	return res
}

type ServiceSitePredictor struct{}

func (ServiceSitePredictor) Predict(prefix string) []string {
	res, _ := StoreOptionsForCompletionPredictor(prefix)
	return res
}

func GenerateCompletion() {

	serviceSitePredictor := new(ServiceSitePredictor)

	storecmd := &complete.Command{
		Flags: map[string]complete.Predictor{
			"w":  predict.Nothing,
			"sc": predict.Nothing,
			"gp": predict.Set(numbers()),
		},
	}

	getcmd := &complete.Command{
		Flags: map[string]complete.Predictor{
			"clip": serviceSitePredictor,
		},
		Args: serviceSitePredictor,
	}

	deletecmd := &complete.Command{
		Args: serviceSitePredictor,
	}

	listcmd := &complete.Command{
		Flags: map[string]complete.Predictor{
			"u": predict.Nothing,
			"p": predict.Nothing,
		},
	}

	gitauthcmd := &complete.Command{
		Flags: map[string]complete.Predictor{
			"w": predict.Nothing,
		},
	}

	helpcmd := &complete.Command{
		Args: predict.Set{
			"store",
			"get",
			"delete",
			"list",
			"gitauth",
			"pull",
			"push",
		},
	}

	cmd := &complete.Command{
		Sub: map[string]*complete.Command{
			"store":   storecmd,
			"get":     getcmd,
			"delete":  deletecmd,
			"list":    listcmd,
			"gitauth": gitauthcmd,
			"pull":    {},
			"push":    {},
			"help":    helpcmd,
		},
	}
	cmd.Complete("gpman")
}
