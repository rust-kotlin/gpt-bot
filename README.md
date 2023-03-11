# gpt-bot
A robot to connect the go-cqhttp and openai
# The outlines
1. The main.go opens the true main function inside the connect.go, which is responsible for the connection with the go-cqhttp.
2. json.go is responsible for the marshal and unmarshal of the important data about the api key, model and the proxy url.
3. client.go is responsible for creating a client of the openai using the data collected from the config.json manipulated by the json.go.
4. chat.go is responsible for post the question to openai and receive the result, which is then be transfered to the connect.go to sent to the go-cqhttp.
