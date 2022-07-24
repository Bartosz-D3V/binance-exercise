[![Go](https://github.com/Bartosz-D3V/binance-exercise/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/Bartosz-D3V/binance-exercise/actions/workflows/go.yml)

# binance-exercise

Trade execution service takes an order size and an order price for an asset pair and generates market
orders, when the liquidity in the order book allows it.

## Running

Application expects three system environment variables that provide high configurability:

`SYMBOL` - String - Symbol of the ticker - i.e. `btcusdt` or `ethusdt`

`QUANTITY_TO_SELL` - Decimal - Amount of assets to sell

`MINIMUM_BID` - Decimal - Minimum bid price in a currency denoted in `SYMBOL` variable (i.e. USD for `btcusdt`)

For convenience, all can be set by using `envDefault` in config/config.go

### Standalone

You can run the application simply by building and executing `main.go` file

### Docker

1. Build docker image: `docker build . -t binance-exercise`
2. Run docker container: `docker run -t binance-exercise`

## Approach to the problem

I decided to create a websocket client that connects to Binance API and converts the response to the appropriate go
struct.

Binance client does not process the tick in any way - instead it sends it to a channel.

Trading service listens to channel and per each tick it decides whether application should exit or continue.

To store and retrieve a list of made transactions I created a repository with SQLLite database driver.

Initially, I thought of summing a list of total assets sold on a database level, but I decided to fetch all entries for
two reasons:

1. Summing REAL data types in SQLLite is not safe (floating points)
2. I would still have to fetch all entries when application exits to print the report

## Problems encountered

The first issue was related to data types. Binance returns many numeric values as strings - converting those to floats
in go is unsafe.

I decided to use 3rd party library to handle proper decimal
values: [shopspring/decimal](https://github.com/shopspring/decimal)

The second issue was related to concurrent database access.

In order to make sure that all database calls are safe, I used go channels and the highest transaction isolation level.

The second issue was also the most enjoyable one - finding a solution to this problem was not trivial and required some
work.

The final issue I encountered was error handling - normally I would add a channel where errors could be sent and
handled.
But due to the fact that this is a Minimal Viable Product (MVP) in the financial sector, I decided that every error
should result in an immediate stop to allow investigation and avoid unexpected behavior.

## Next steps

I decided not to spend more than 6 hours on this exercise so there are many areas that could be improved:

1. Add integration tests (i.e., using wiremock)
2. Add database tests for repository
3. Add alerts to the application (i.e. using elastalert)
4. Add database versioning (i.e. using flyway)
5. Remove `envDefault` tags from Config struct - those should be defined by IaC tool (i.e. Terraform)
