# stripe-integration-golang

Stripe integration using Golang.

Simple back-end service to request a payment intent from Stripe.

config.json contains settings:
- RouterUrl: the url to use for the golang service
- AcceptedOrigin: url of the React front-end
- ApiMode: mode in which to run gin
- StripeKey: secret key for the Stripe account