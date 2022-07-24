FROM golang:1.18.4-buster

ENV GO111MODULE=on
ENV GOFLAGS="-mod=mod"

# Override those
ENV SYMBOL="ethusdt"
ENV QUANTITY_TO_SELL=12.33
ENV MINIMUM_BID=1605.45

ENV APP_HOME /binance-exercise
ADD . $APP_HOME

WORKDIR "$APP_HOME"
RUN go build -o main $APP_HOME .

CMD ["/binance-exercise/recreate_mocks.sh"]
CMD ["/binance-exercise/main"]
