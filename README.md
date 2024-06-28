### Test repo to show bug in htmx ws extension

When submitting via htmx the form array items are not encoding in json as an array when only a single item exists

To start test:
templ generate
go run main.go 

localhost:8080

Scenarios:
- standard post with single array item - decodes via struct definition correctly
- standard post with 2 array items - decodes via struct definition correctly
- ws message with single arraya item - decodes via struct definition correctly
- ws message with 2 array items in form - fails to decode via struct definition
