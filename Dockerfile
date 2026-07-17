FROM golang:1.26.4 AS development

WORKDIR /budget

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG GOARCH=arm64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -o ./build/budget ./main.go

FROM gcr.io/distroless/static-debian13:nonroot AS app

EXPOSE 3000

COPY --from=development /budget/build/budget /budget

CMD [ "/budget" ]