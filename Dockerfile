FROM golang:1.25.6 AS development

WORKDIR /budget

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/budget ./main.go
RUN chmod a+x /budget

FROM golang:1.25.6-alpine3.22 AS app

EXPOSE 3000

COPY --from=development /budget/build/budget /budget

CMD [ "/budget" ]