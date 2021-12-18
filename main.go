package main

import (
	"log"
	"net/http"

	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"

	"github.com/sagikazarmark/kube-curiesync-injector/webhook/mutation"
	"github.com/sagikazarmark/kube-curiesync-injector/webhook/mutator"
)

func main() {
	mutator := mutator.Mutator{
		CuriesyncInjector: mutation.CuriesyncInjector{
			CuriesyncImage: "",
			BucketLink:     "",
		},
	}

	wh, err := kwhmutating.NewWebhook(kwhmutating.WebhookConfig{
		ID:      "curiesync-injector",
		Mutator: mutator,
		// Logger:  logger,
	})
	if err != nil {
		panic(err)
	}

	handler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{
		Webhook: wh,
		// Logger: logger
	})
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe("127.0.0.1:8080", handler)
	if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
