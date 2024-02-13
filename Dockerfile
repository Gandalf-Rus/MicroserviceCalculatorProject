FROM golang

WORKDIR /Project

COPY . .

# CMD [ "go", "run", "backend/cmd/main.go" ]
