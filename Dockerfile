FROM golang:1.26.0 AS development

WORKDIR /budget

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG GOARCH=arm64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -o ./build/budget ./main.go
RUN chmod a+x /budget

FROM alpine:3.23.3 AS app

EXPOSE 3000

COPY --from=development /budget/build/budget /budget

CMD [ "/budget" ]