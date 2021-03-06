# Implementation

[Game Reference](https://www.youtube.com/watch?v=SreZpq4LCdA&ab_channel=CasinoRamaResort)

-   The Deck on table has 6 \* 52 card (as per [wiki](https://en.wikipedia.org/wiki/Casino_War))
-   Player Enters amount of chips to bet
    -   Min bet amount is 10 chips
    -   Max bet amount is 500 chips

### Rules

-   Two cards are drawn, one for player and one for dealer

    -   If player card is higher than dealers card, dealer pays an amount equal to the bet amount
    -   If player card is lower than dealers card, dealer collects the amount

-   In case of a tie player has two choice
    -   Surrender : player takes half the bet amount and quits
    -   Go to war
        -   Player & dealer put additional chips equavilant to the bet.
            -   Dealer burns 3 cards, draws one card for player
            -   Dealer again burns 3 cards, draws one card for himself
                -   The winner takes the amount
                -   If it is a tie again, player gets 10x his bet

# Server

## Building the image

```sh
cd server

docker build . -t casino-war-server

docker-compose up # Runs mongoDB
```

## Run container with Shell

```sh
docker run -p 8080:8080 --network=host -it -v ${PWD}:/usr/src/casino-war/server casino-war-server sh
```

## Run application

```sh
go build
./main
```

# Client

## Building the image

```sh
$ cd cli-client

$ docker build . -t casino-war-client
```

## Run container with Shell

```sh
docker run --network=host -it -v ${PWD}:/usr/src/casino-war/client casino-war-client sh
```

## Run application

```sh
go build

./main
```

# Test game

```sh
cd server

go test ./game -v
```


# Demo
[Demo Video](https://aeroqubetechnologies-my.sharepoint.com/:v:/g/personal/ashutosh_kumar_aeroqube_com/EddqCfo8aZdPm54z80Zt27cBXievWmY7Ej-aAhhEPaRCkA?e=48NMu5)